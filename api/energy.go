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

func boiler_efficiency_day(c *gin.Context) []float64 {
	var finalData []float64
	// 根据当前时间查redis有无已计算好的数据
	// now := time.Now().Local()
	now, _ := time.Parse("2006/01/02 15:04:05", "2022/09/13 03:37:02")
	dayStr := fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day())
	hourStr := fmt.Sprintf("%d/%02d/%02d %02d", now.Year(), now.Month(), now.Day(), now.Hour())
	data, err := model.RedisClient.LRange(dayStr+" boiler_efficiency_day", 0, int64(now.Hour())).Result()
	lredis := len(data)
	if err == nil && lredis == now.Hour()+1 {
		available := true
		for i := 0; i < lredis; i++ {
			floatData, err := strconv.ParseFloat(data[i], 64)
			if err != nil {
				available = false
				break
			}
			finalData = append(finalData, floatData)
		}
		if available {
			return finalData
		}
	}
	// redis没有，去mongo查当天前几个小时的数据
	var result defs.CalculationResultFloat
	err = model.MongoResult.FindOne(context.TODO(), bson.D{{"time", dayStr}, {"name", "boiler_efficiency_day"}}).Decode(&result)
	finalData = result.Value
	// mongo没有的数据和这个小时的数据重新计算
	if len(finalData) == now.Hour()+1 {
		finalData = finalData[:len(finalData)-1]
	}
	for i := len(finalData); i <= now.Hour(); i++ {
		ans := calc.BoilerEfficiency(hourStr)
		finalData = append(finalData, ans)
	}
	// 结果写入mongo
	model.MongoResult.DeleteOne(context.TODO(), bson.D{{"time", dayStr}, {"name", "boiler_efficiency_day"}})
	model.MongoResult.InsertOne(context.TODO(), bson.D{{"time", dayStr}, {"name", "boiler_efficiency_day"}, {"value", finalData}})
	// 并存入redis
	if lredis == now.Hour() {
		// 只用插入最新数据即可
		model.RedisClient.RPush(dayStr+" boiler_efficiency_day", finalData[len(finalData)-1])
	} else {
		// 重新写入数据并设置ttl
		model.RedisClient.Del(dayStr + " boiler_efficiency_day")
		for i := 0; i < len(finalData); i++ {
			model.RedisClient.RPush(dayStr+" boiler_efficiency_day", finalData[i])
		}
		model.RedisClient.Expire(dayStr+" boiler_efficiency_day", time.Minute) //每分钟更新一次
	}
	return finalData
}

func GetPageData(c *gin.Context) {
	page := c.Query("page")
	switch page {
	case "analyse-energy-station": //能效分析-能源站
		// 电锅炉热效率
		d1 := boiler_efficiency_day(c)
		print(d1)
	}
}
