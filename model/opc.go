package model

import (
	"context"
	"energy/defs"
	"energy/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func GetResultFloat(name string, time string) (float64, bool) {
	var result defs.CalculationResultFloat
	err := MongoResult.FindOne(context.TODO(), bson.D{{"name", name}, {"time", time}}).Decode(&result)
	if err != nil {
		return 0, false
	}
	return result.Value, true
}

func GetOpcBoolList(itemid string, time string) ([]bool, bool) {
	var opcData defs.OpcData
	err := MongoOPC.FindOne(context.TODO(), bson.D{{"itemid", itemid}, {"time", time}}).Decode(&opcData)
	if err != nil {
		return nil, false
	}
	ans := []bool{}
	for _, v := range opcData.Value {
		val, _ := v.(bool) //失败的值视为0
		ans = append(ans, val)
	}
	return ans, true

}

func GetOpcFloatList(itemid string, time string) ([]float64, bool) {
	var opcData defs.OpcData
	err := MongoOPC.FindOne(context.TODO(), bson.D{{"itemid", itemid}, {"time", time}}).Decode(&opcData)
	if err != nil {
		return nil, false
	}
	ans := []float64{}
	for _, v := range opcData.Value {
		val, _ := utils.GetFloat(v) //失败的值视为0
		ans = append(ans, val)
	}
	return ans, true

}

/*
func GetOpcDataList(tableName string, timeType int) []interface{} { //timeType: 0-day, 1-hour，2-近7天, 3-过去一年每月
	var finalData [100]interface{}
	lenFin := 0
	// 根据当前时间查redis有无已计算好的数据
	// now := time.Now().Local()
	now, _ := time.Parse("2006/01/02 15:04:05", "2022/10/13 15:31:00")
	timeStr := fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day()) //如果是近7天，在redis里面要存储近7天的值，但是在mongo里面只存储当天的值，月同理
	tNum := now.Hour()
	if timeType == 1 {
		timeStr = fmt.Sprintf("%s %02d", timeStr, now.Hour())
		tNum = now.Minute()
	}
	if timeType == 2 {
		tNum = 6
	}
	data, err := RedisClient.LRange(timeStr+" "+tableName, 0, int64(tNum)).Result()
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
	// redis没有，去mongo查
	if timeType == 0 || timeType == 1 {
		var result defs.CalculationResultFloatList
		_ = MongoResult.FindOne(context.TODO(), bson.D{{"time", timeStr}, {"name", tableName}}).Decode(&result)
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
			err = MongoResult.FindOne(context.TODO(), bson.D{{"time", fmt.Sprintf("%d/%02d/%02d", startTime.Year(), startTime.Month(), startTime.Day())},
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

	startTime := now.Add(-time.Hour * 24 * 7) //7天前
	for i := 0; i <= tNum; i++ {
		startTime = startTime.Add(time.Hour * 24)
		newStr := fmt.Sprintf("%d/%02d/%02d", startTime.Year(), startTime.Month(), startTime.Day())
		if timeType == 0 {
			if i >= lenFin {
				finalData[i] = Calc(tableName, fmt.Sprintf("%s %02d", timeStr, i))
				lenFin++
			}
		}
		if timeType == 1 {
			if i >= lenFin {
				finalData[i] = Calc(tableName, MinParam{HourStr: timeStr, Min: i})
				lenFin++
			}
		}
		if timeType == 2 {
			if needCalc[i] {
				finalData[i] = Calc(tableName, newStr)
				MongoResult.DeleteOne(context.TODO(), bson.D{{"time", newStr}, {"name", tableName}})
				MongoResult.InsertOne(context.TODO(), bson.D{{"time", newStr}, {"name", tableName}, {"value", finalData[i]}})
			}
		}
	}
	if timeType != 2 {
		// 结果写入mongo
		MongoResult.DeleteOne(context.TODO(), bson.D{{"time", timeStr}, {"name", tableName}})
		MongoResult.InsertOne(context.TODO(), bson.D{{"time", timeStr}, {"name", tableName}, {"value", finalData[:lenFin]}})
	}
	// 并存入redis
	if lredis == tNum {
		// 只用插入最新数据即可
		RedisClient.RPush(timeStr+" "+tableName, finalData[lenFin-1])
	} else {
		// 重新写入数据并设置ttl
		RedisClient.Del(timeStr + " " + tableName)
		for i := 0; i < lenFin; i++ {
			RedisClient.RPush(timeStr+" "+tableName, finalData[i])
		}
		RedisClient.Expire(timeStr+" "+tableName, time.Minute) //每分钟更新一次
	}

	return finalData[:lenFin]
}*/
