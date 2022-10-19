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

func getFloat(val interface{}) (float64, bool) {
	if val1, ok := val.(float64); ok {
		return val1, true
	}
	if val2, ok := val.(int32); ok {
		return float64(val2), true
	}
	if val3, ok := val.(int64); ok {
		return float64(val3), true
	}
	return 0, false
}

func getOpcBoolList(itemid string, time string) ([]bool, bool) {
	var opcData defs.OpcData
	err := model.MongoOPC.FindOne(context.TODO(), bson.D{{"itemid", itemid}, {"time", time}}).Decode(&opcData)
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

func getOpcFloatList(itemid string, time string) ([]float64, bool) {
	var opcData defs.OpcData
	err := model.MongoOPC.FindOne(context.TODO(), bson.D{{"itemid", itemid}, {"time", time}}).Decode(&opcData)
	if err != nil {
		return nil, false
	}
	ans := []float64{}
	for _, v := range opcData.Value {
		val, _ := getFloat(v) //失败的值视为0
		ans = append(ans, val)
	}
	return ans, true

}

func BoilerEfficiency(hourStr string) float64 {
	q1 := 0.0
	q2 := 0.0
	for i := 1; i <= 4; i++ {
		Tout, ok := getOpcFloatList("server.%E9%94%85%E7%82%89%E5%AE%9E%E9%99%85%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6"+fmt.Sprint(i), hourStr) //锅炉实际出水温度
		if !ok {
			continue
		}
		Tin, ok := getOpcFloatList("server.%E9%94%85%E7%82%89%E5%AE%9E%E9%99%85%E5%9B%9E%E6%B0%B4%E6%B8%A9%E5%BA%A6"+fmt.Sprint(i), hourStr) //锅炉实际回水温度
		if !ok {
			continue
		}
		Oa, ok := getOpcBoolList("server.A%E6%B3%B5%E8%BF%90%E8%A1%8C"+fmt.Sprint(i), hourStr) //A泵运行
		if !ok {
			continue
		}
		Ob, ok := getOpcBoolList("server.B%E6%B3%B5%E8%BF%90%E8%A1%8C"+fmt.Sprint(i), hourStr) //B泵运行
		if !ok {
			continue
		}
		w, ok := getOpcFloatList("server.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(i), hourStr) //功率采集
		if !ok {
			continue
		}
		minLen := utils.Min(len(Tout), len(Tin), len(Oa), len(Ob), len(w))
		for j := 0; j < minLen; j++ {
			q1 += (utils.Bool2Float(Oa[j]) + utils.Bool2Float(Ob[j])) * (Tout[j] - Tin[j])
			q2 += w[j]
		}
	}
	q1 *= 4.2 * 137 * 1e6 / 60
	q2 *= 60000 //目前功率单位未知，公式暂用kW
	if q2 == 0 {
		return 0
	}
	return q1 / q2
}
