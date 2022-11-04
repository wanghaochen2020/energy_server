package model

import (
	"context"
	"energy/defs"
	"energy/utils"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type MinParam struct {
	HourStr string
	Min     int
}

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

func getOpcFloatList(itemid string, time string) ([]float64, bool) {
	var opcData defs.OpcData
	err := MongoOPC.FindOne(context.TODO(), bson.D{{"itemid", itemid}, {"time", time}}).Decode(&opcData)
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

func boilerEfficiency(hourStr string) float64 {
	q1 := 0.0
	q2 := 0.0
	for i := 1; i <= 4; i++ {
		Tout, ok := getOpcFloatList("ZLZ.%E9%94%85%E7%82%89%E5%AE%9E%E9%99%85%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6"+fmt.Sprint(i), hourStr) //锅炉实际出水温度
		if !ok {
			continue
		}
		Tin, ok := getOpcFloatList("ZLZ.%E9%94%85%E7%82%89%E5%AE%9E%E9%99%85%E5%9B%9E%E6%B0%B4%E6%B8%A9%E5%BA%A6"+fmt.Sprint(i), hourStr) //锅炉实际回水温度
		if !ok {
			continue
		}
		Oa, ok := getOpcBoolList("ZLZ.A%E6%B3%B5%E8%BF%90%E8%A1%8C"+fmt.Sprint(i), hourStr) //A泵运行
		if !ok {
			continue
		}
		Ob, ok := getOpcBoolList("ZLZ.B%E6%B3%B5%E8%BF%90%E8%A1%8C"+fmt.Sprint(i), hourStr) //B泵运行
		if !ok {
			continue
		}
		w, ok := getOpcFloatList("ZLZ.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(i), hourStr) //功率采集
		if !ok {
			continue
		}
		minLen := utils.Min(len(Tout), len(Tin), len(Oa), len(Ob), len(w))
		for j := 0; j < minLen; j++ {
			if utils.Zero(Tout[j], Tin[j]) {
				continue
			}
			q1 += (utils.Bool2Float(Oa[j]) + utils.Bool2Float(Ob[j])) * (Tout[j] - Tin[j])
			q2 += w[j]
		}
	}
	q1 *= 4.2 * 137 * 1e6 / 60
	q2 *= 60000 //目前功率单位未知，公式暂用kW，最后的值的单位是J
	if q2 == 0 {
		return 0
	}
	return q1 / q2
}

func watertankEfficiency(hourStr string) float64 {
	q1 := 0.0
	var Tinitial float64
	Taver := 0.0
	h, ok := getOpcFloatList("ZLZ.OUTPUT_P10", hourStr)
	if !ok {
		return 0
	}
	Tin, ok := getOpcFloatList("ZLZ.OUTPUT_T3", hourStr)
	if !ok {
		return 0
	}
	Tout, ok := getOpcFloatList("ZLZ.OUTPUT_T4", hourStr)
	if !ok {
		return 0
	}
	minLen := utils.Min(len(h), len(Tin), len(Tout))
	if minLen == 0 {
		return 0
	}
	if utils.Zero(Tin[0]) {
		return 0
	}
	Tinitial = Tin[0]
	Taver = Tinitial
	averCount := 1
	for i := 1; i < minLen; i++ {
		if utils.Zero(h[i], h[i-1], Tin[i], Tout[i]) {
			continue
		}
		q1 += ((h[i]-h[i-1])*100 + 0*0 /*预留给小水箱的位置，目前未得到相关数据*/) * (Tin[i] - Tout[i])
		Taver += Tin[i]
		averCount++
	}
	Taver /= float64(averCount)
	q2 := (Taver - Tinitial) * 1055
	if q2 == 0 {
		return 0
	}
	return q1 / q2
}

func energystationEfficiency(hourStr string) float64 {
	totalQ1, ok := getOpcFloatList("ZLZ.RLB_%E7%B4%AF%E8%AE%A1%E7%83%AD%E9%87%8F", hourStr) //RLB_累计热量
	if !ok {
		return 0
	}
	APGL2, ok := getOpcFloatList("ZLZ.%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6APGL2", hourStr) //有功电度APGL2
	if !ok {
		return 0
	}
	AZ1, ok := getOpcFloatList("ZLZ.%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6AZ_1", hourStr) //有功电度AZ_1
	if !ok {
		return 0
	}
	w := make([][]float64, 5)
	for i := 1; i <= 4; i++ {
		w[i], _ = getOpcFloatList("ZLZ.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(i), hourStr) //功率采集
	}
	minLen := utils.Min(len(totalQ1), len(APGL2), len(AZ1), utils.Max(len(w[1]), len(w[2]), len(w[3]), len(w[4])))
	left := 0
	right := 0
	for i := 0; i < minLen; i++ {
		if utils.Zero(totalQ1[i], APGL2[i], AZ1[i]) {
			continue
		}
		left = i
		break
	}
	for i := minLen - 1; i >= 0; i-- {
		if utils.Zero(totalQ1[i], APGL2[i], AZ1[i]) {
			continue
		}
		right = i
		break
	}
	if left >= right {
		return 0
	}
	q23 := 0.0
	for i := 1; i <= 4; i++ {
		for j := left; j < utils.Min(right+1, len(w[i])); j++ {
			q23 += w[i][j]
		}
	}
	q23 /= 60 //目前功率单位未知，公式暂用kW，最后的值的单位是kW·h
	q1 := totalQ1[right] - totalQ1[left]
	q2 := APGL2[right] - APGL2[left] + AZ1[right] - AZ1[left] + q23
	if q2 == 0 {
		return 0
	}
	return q1 / q2
}

func deviceOnlineRate(minStr string) float64 {
	// MongoOPC.Find()
	return 0
}

func boilerPower(hourStr string, min int) float64 {
	ans := 0.0
	for i := 1; i <= 4; i++ {
		w, _ := getOpcFloatList("ZLZ.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(i), hourStr) //功率采集
		if len(w) < min {
			continue
		}
		ans += w[min]
	}
	return ans
}

func energystationCarbonHour(hourStr string) float64 { //计算这个小时的碳排
	APGL2, ok := getOpcFloatList("ZLZ.%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6APGL2", hourStr) //有功电度APGL2
	if !ok {
		return 0
	}
	AZ1, ok := getOpcFloatList("ZLZ.%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6AZ_1", hourStr) //有功电度AZ_1
	if !ok {
		return 0
	}
	q23 := 0.0
	for i := 1; i <= 4; i++ {
		w, ok := getOpcFloatList("ZLZ.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(i), hourStr) //功率采集
		if !ok {
			continue
		}
		for j := 0; j < len(w); j++ {
			q23 += w[j]
		}
	}
	l1 := -1
	for i := 0; i < len(APGL2); i++ {
		if utils.Zero(APGL2[i]) {
			continue
		}
		l1 = i
		break
	}
	r1 := -1
	for i := len(APGL2) - 1; i >= 0; i-- {
		if utils.Zero(APGL2[i]) {
			continue
		}
		r1 = i
		break
	}
	l2 := -1
	for i := 0; i < len(AZ1); i++ {
		if utils.Zero(AZ1[i]) {
			continue
		}
		l2 = i
		break
	}
	r2 := -1
	for i := len(AZ1) - 1; i >= 0; i-- {
		if utils.Zero(AZ1[i]) {
			continue
		}
		r2 = i
		break
	}
	q21 := 0.0
	q22 := 0.0
	if l1 != -1 {
		q21 = APGL2[r1] - APGL2[l1]
	}
	if l2 != -1 {
		q21 = AZ1[r2] - AZ1[l2]
	}
	t := (q21 + q22 + q23/60) / 1000 * 0.604 //吨CO2
	return t
}

func energystationCarbonDay(dayStr string) float64 { //计算当天的碳排
	sum := 0.0
	for i := 0; i < 24; i++ {
		hourStr := fmt.Sprintf("%s %02d", dayStr, i)
		var result defs.CalculationResultFloat
		err = MongoResult.FindOne(context.TODO(), bson.D{{"time", hourStr}, {"name", "energystation_carbon_week"}}).Decode(&result)
		if err != nil {
			ans := energystationCarbonHour(hourStr)
			MongoResult.InsertOne(context.TODO(), bson.D{{"time", hourStr}, {"name", "energystation_carbon_week"}, {"value", ans}})
			sum += ans
		} else {
			sum += result.Value
		}
	}
	return sum
}

func energyPayload(hourStr string) float64 {
	q1 := 0.0
	for i := 1; i <= 4; i++ {
		Tout, ok := getOpcFloatList("ZLZ.%E9%94%85%E7%82%89%E5%AE%9E%E9%99%85%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6"+fmt.Sprint(i), hourStr) //锅炉实际出水温度
		if !ok {
			continue
		}
		Tin, ok := getOpcFloatList("ZLZ.%E9%94%85%E7%82%89%E5%AE%9E%E9%99%85%E5%9B%9E%E6%B0%B4%E6%B8%A9%E5%BA%A6"+fmt.Sprint(i), hourStr) //锅炉实际回水温度
		if !ok {
			continue
		}
		Oa, ok := getOpcBoolList("ZLZ.A%E6%B3%B5%E8%BF%90%E8%A1%8C"+fmt.Sprint(i), hourStr) //A泵运行
		if !ok {
			continue
		}
		Ob, ok := getOpcBoolList("ZLZ.B%E6%B3%B5%E8%BF%90%E8%A1%8C"+fmt.Sprint(i), hourStr) //B泵运行
		if !ok {
			continue
		}
		minLen := utils.Min(len(Tout), len(Tin), len(Oa), len(Ob))
		for j := 0; j < minLen; j++ {
			if utils.Zero(Tout[j], Tin[j]) {
				continue
			}
			q1 += (utils.Bool2Float(Oa[j]) + utils.Bool2Float(Ob[j])) * (Tout[j] - Tin[j])
		}
	}
	// q1 *= 4.2 * 137 * 1e6 / 60
	// q2 := 4 * 4 * 1e6 * 3600
	q1 *= 2.3975
	q2 := 14400.0
	return q1 / q2
}

func Calc(tableName string, params interface{}) interface{} {
	switch tableName {
	case "boiler_efficiency_day":
		p, ok := params.(string)
		if ok {
			return boilerEfficiency(p)
		}
	case "watertank_efficiency_day":
		p, ok := params.(string)
		if ok {
			return watertankEfficiency(p)
		}
	case "energystation_efficiency_day":
		p, ok := params.(string)
		if ok {
			return energystationEfficiency(p)
		}
	case "device_online_rate_hour":
		p, ok := params.(string)
		if ok {
			return deviceOnlineRate(p)
		}
	case "boiler_power_hour":
		p, ok := params.(MinParam)
		if ok {
			return boilerPower(p.HourStr, p.Min)
		}
	case "energystation_carbon_day":
		p, ok := params.(string)
		if ok {
			return energystationCarbonHour(p) //表的名字是day，但是求的是小时的，下面同理
		}
	case "energystation_carbon_week":
		p, ok := params.(string)
		if ok {
			return energystationCarbonDay(p)
		}
	case "energy_pay_load":
		p, ok := params.(string)
		if ok {
			return energyPayload(p)
		}
	}
	return nil
}
