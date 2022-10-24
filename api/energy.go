package api

import (
	"context"
	"energy/calc"
	"energy/defs"
	"energy/model"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	err error
	ok  bool
)

func getOpcDataList(c *gin.Context, tableName string, timeType int) []interface{} { //timeType: 0-day, 1-hour，2-近7天
	// boiler_efficiency_day
	var finalData [100]interface{}
	lenFin := 0
	// 根据当前时间查redis有无已计算好的数据
	// now := time.Now().Local()
	now, _ := time.Parse("2006/01/02 15:04:05", "2022/09/13 03:37:02")
	timeStr := fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day()) //如果是近7天，在redis里面要存储近7天的值，但是在mongo里面只存储当天的值
	tNum := now.Hour()
	if timeType == 1 {
		timeStr = fmt.Sprintf("%s %02d", timeStr, now.Hour())
		tNum = now.Minute()
	}
	if timeType == 2 {
		tNum = 6
	}
	data, err := model.RedisClient.LRange(timeStr+" "+tableName, 0, int64(tNum)).Result()
	lredis := len(data)
	if err == nil && lredis == tNum+1 {
		available := true
		for i := 0; i < lredis; i++ {
			floatData, err := strconv.ParseFloat(data[i], 64)
			if err != nil {
				available = false
				break
			}
			finalData[i] = floatData
			lenFin++
		}
		if available {
			return finalData[:lenFin]
		}
	}
	// redis没有，去mongo查当天前几个小时的数据
	if timeType == 0 || timeType == 1 {
		var result defs.CalculationResultFloatList
		_ = model.MongoResult.FindOne(context.TODO(), bson.D{{"time", timeStr}, {"name", tableName}}).Decode(&result)
		for i, v := range result.Value {
			finalData[i] = v
		}
		lenFin = len(result.Value)
	}
	var needCalc [7]bool
	if timeType == 2 {
		startTime := now.Add(-time.Hour * 24 * 6) //6天前
		var result defs.CalculationResultFloat
		for i := 0; i < 7; i++ {
			startTime = startTime.Add(time.Hour * 24)
			err = model.MongoResult.FindOne(context.TODO(), bson.D{{"time", fmt.Sprintf("%d/%02d/%02d", startTime.Year(), startTime.Month(), startTime.Day())},
				{"name", tableName}}).Decode(&result)
			if err != nil {
				needCalc[i] = true
			}
			finalData[i] = result.Value
		}
		lenFin = 7
	}
	// mongo没有的数据和这个小时的数据重新计算
	if timeType == 2 {
		needCalc[lenFin-1] = true
	} else {
		if lenFin == tNum+1 {
			lenFin--
		}
	}

	startTime := now.Add(-time.Hour * 24 * 6) //6天前
	for i := 0; i <= tNum; i++ {
		startTime = startTime.Add(time.Hour * 24)
		newStr := fmt.Sprintf("%d/%02d/%02d", startTime.Year(), startTime.Month(), startTime.Day())
		if timeType == 0 {
			if i >= lenFin {
				finalData[i] = calc.Calc(tableName, fmt.Sprintf("%s %02d", timeStr, i))
				lenFin++
			}
		}
		if timeType == 1 {
			if i >= lenFin {
				finalData[i] = calc.Calc(tableName, calc.MinParam{HourStr: timeStr, Min: i})
				lenFin++
			}
		}
		if timeType == 2 {
			if needCalc[i] {
				finalData[i] = calc.Calc(tableName, newStr)
				model.MongoResult.DeleteOne(context.TODO(), bson.D{{"time", newStr}, {"name", tableName}})
				model.MongoResult.InsertOne(context.TODO(), bson.D{{"time", newStr}, {"name", tableName}, {"value", finalData}})
			}
		}
	}
	if timeType != 2 {
		// 结果写入mongo
		model.MongoResult.DeleteOne(context.TODO(), bson.D{{"time", timeStr}, {"name", tableName}})
		model.MongoResult.InsertOne(context.TODO(), bson.D{{"time", timeStr}, {"name", tableName}, {"value", finalData}})
	}
	// 并存入redis
	if lredis == tNum {
		// 只用插入最新数据即可
		model.RedisClient.RPush(timeStr+" "+tableName, finalData[lenFin-1])
	} else {
		// 重新写入数据并设置ttl
		model.RedisClient.Del(timeStr + " " + tableName)
		for i := 0; i < lenFin; i++ {
			model.RedisClient.RPush(timeStr+" "+tableName, finalData[i])
		}
		model.RedisClient.Expire(timeStr+" "+tableName, time.Minute) //每分钟更新一次
	}

	return finalData[:lenFin]
}

func GetPageData(c *gin.Context) {
	page := c.Query("page")
	switch page {
	case "system-energy-station": //系统层-能源站
		//设备在线率
		dor := getOpcDataList(c, "device_online_rate_hour", 1)
		print(dor)
		//锅炉总功率
		bp := getOpcDataList(c, "boiler_power_hour", 1)
		print(bp)
	case "analyse-energy-station": //能效分析-能源站
		// 电锅炉热效率
		d11 := getOpcDataList(c, "boiler_efficiency_day", 0)
		print(d11)
		// 蓄热水箱热效率
		d12 := getOpcDataList(c, "watertank_efficiency_day", 0)
		print(d12)
		// 总效率
		d13 := getOpcDataList(c, "energystation_efficiency_day", 0)
		print(d13)
		//日碳排
		d21 := getOpcDataList(c, "energystation_carbon_day", 0)
		print(d21)
		//周碳排
		d22 := getOpcDataList(c, "energystation_carbon_week", 2)
		print(d22)
	}
}
