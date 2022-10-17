package calc

import (
	"context"
	"energy/defs"
	"energy/model"
	"energy/utils"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	err error
)

func bool2Float(b bool) float64 {
	if b {
		return 1.0
	} else {
		return 0.0
	}
}

func getFloat(val defs.OpcData) (float64, bool) {
	if val1, ok := val.Value[0].(float64); ok {
		return val1, true
	}
	if val2, ok := val.Value[0].(int32); ok {
		return float64(val2), true
	}
	return 0, false
}

func getBool(val defs.OpcData) (bool, bool) {
	if val, ok := val.Value[0].(bool); ok {
		return val, true
	}
	return false, false
}

func getOpcFloat(itemid string, time string) (float64, bool) {
	var opcData defs.OpcData
	err := model.MongoOPC.FindOne(context.TODO(), bson.D{{"itemid", itemid}, {"time", time}}).Decode(&opcData)
	if err != nil {
		return 0, false
	}
	val, ok := getFloat(opcData)
	if !ok {
		return 0, false
	}
	return val, true
}

func getOpcBool(itemid string, time string) (bool, bool) {
	var opcData defs.OpcData
	err := model.MongoOPC.FindOne(context.TODO(), bson.D{{"itemid", itemid}, {"time", time}}).Decode(&opcData)
	if err != nil {
		return false, false
	}
	val, ok := getBool(opcData)
	if !ok {
		return false, false
	}
	return val, true
}

func BoilerEfficiency(hourStr string) float64 {
	q1 := 0.0
	q2 := 0.0
	for i := 0; i < 60; i += utils.DeltaT {
		minStr := fmt.Sprintf("%s:%02d", hourStr, i)
		for j := 1; j <= 4; j++ {
			Tout, ok := getOpcFloat("server.%E9%94%85%E7%82%89%E5%AE%9E%E9%99%85%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6"+fmt.Sprint(j), minStr) //锅炉实际出水温度
			if !ok {
				continue
			}
			Tin, ok := getOpcFloat("server.%E9%94%85%E7%82%89%E5%AE%9E%E9%99%85%E5%9B%9E%E6%B0%B4%E6%B8%A9%E5%BA%A6"+fmt.Sprint(j), minStr) //锅炉实际回水温度
			if !ok {
				continue
			}
			Oa, ok := getOpcBool("server.A%E6%B3%B5%E8%BF%90%E8%A1%8C"+fmt.Sprint(j), minStr) //A泵运行
			if !ok {
				continue
			}
			Ob, ok := getOpcBool("server.B%E6%B3%B5%E8%BF%90%E8%A1%8C"+fmt.Sprint(j), minStr) //B泵运行
			if !ok {
				continue
			}
			w, ok := getOpcFloat("server.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(j), minStr) //功率采集
			if !ok {
				continue
			}

			q1 += (bool2Float(Oa) + bool2Float(Ob)) * (Tout - Tin)
			q2 += w
		}
	}
	q1 *= 4.2 * 137 * 1e6 / 60
	q2 *= 60000 //目前功率单位未知，公式暂用kW
	if q2 == 0 {
		return 0
	}
	return q1 / q2
}
