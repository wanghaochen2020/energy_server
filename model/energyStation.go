package model

import (
	"context"
	"energy/defs"
	"energy/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type MinParam struct {
	HourStr string
	Min     int
}

//能源站在线率
func CalcEnergyOnlineRate(hourStr string) (float64, bool) {
	var result defs.MongoCountResult
	command := bson.D{{"count", "opc_data"}, {"query", bson.D{{"time", hourStr}}}}
	MongoOPC.Database().RunCommand(context.TODO(), command).Decode(&result)
	return float64(result.N) / 841, result.Ok
}

//能源站锅炉总功率
func CalcEnergyBoilerPower(hourStr string, min int) float64 {
	ans := 0.0
	for i := 1; i <= 4; i++ {
		w, _ := GetOpcFloatList("ZLZ.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(i), hourStr) //功率采集
		if len(w) < min {
			continue
		}
		ans += w[min]
	}
	return ans
}

//能源站今日总耗能，每分钟
func CalcEnergyPowerConsumptionToday(t time.Time) float64 {
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	hourStr := fmt.Sprintf("%s %02d", dayStr, t.Hour())
	ans := 0.0
	cost, ok := GetResultFloatList(defs.EnergyCarbonDay, dayStr)
	if ok {
		minLen := utils.Min(len(cost), t.Hour())
		for i := 0; i <= minLen; i++ {
			ans += cost[i]
		}
	}
	ans += CalcEnergyCarbonHour(hourStr) //如果卡就删掉这一行，然后把这个函数做成一个小时调用一次
	ans *= 1000 / 0.604                  //由tCO2换算到kW·h
	return ans
}

//能源站锅炉运行台数
func CalcEnergyBoilerRunningNum(hourStr string, min int) float64 {
	ans := 0.0
	for i := 1; i <= 4; i++ {
		w, _ := GetOpcFloatList("ZLZ.%E7%B3%BB%E7%BB%9F%E8%BF%90%E8%A1%8C%E4%B8%AD"+fmt.Sprint(i), hourStr) //运行状态
		if len(w) < min {
			continue
		}
		ans += w[min]
	}
	return ans
}

//电极锅炉供热量
func CalcEnergyBoilerHeatSupply(hourStr string) float64 {
	q1 := 0.0
	for i := 1; i <= 4; i++ {
		Tout, ok := GetOpcFloatList("ZLZ.%E9%94%85%E7%82%89%E5%AE%9E%E9%99%85%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6"+fmt.Sprint(i), hourStr) //锅炉实际出水温度
		if !ok {
			continue
		}
		Tin, ok := GetOpcFloatList("ZLZ.%E9%94%85%E7%82%89%E5%AE%9E%E9%99%85%E5%9B%9E%E6%B0%B4%E6%B8%A9%E5%BA%A6"+fmt.Sprint(i), hourStr) //锅炉实际回水温度
		if !ok {
			continue
		}
		Oa, ok := GetOpcFloatList("ZLZ.A%E6%B3%B5%E8%BF%90%E8%A1%8C"+fmt.Sprint(i), hourStr) //A泵运行
		if !ok {
			continue
		}
		Ob, ok := GetOpcFloatList("ZLZ.B%E6%B3%B5%E8%BF%90%E8%A1%8C"+fmt.Sprint(i), hourStr) //B泵运行
		if !ok {
			continue
		}
		minLen := utils.Min(len(Tout), len(Tin), len(Oa), len(Ob))
		for j := 0; j < minLen; j++ {
			if utils.Zero(Tout[j], Tin[j]) {
				continue
			}
			q1 += (Oa[j] + Ob[j]) * (Tout[j] - Tin[j])
		}
	}
	q1 *= 4.2 * 137 * 1e6 / 60
	return q1
}

//锅炉效率
func CalcEnergyBoilerEfficiency(hourStr string, q1 float64) float64 {
	q2 := 0.0
	for i := 1; i <= 4; i++ {
		w, ok := GetOpcFloatList("ZLZ.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(i), hourStr) //功率采集
		if !ok {
			continue
		}
		minLen := len(w)
		for j := 0; j < minLen; j++ {
			q2 += w[j]
		}
	}
	q2 *= 60000 //功率单位kW
	if q2 == 0 {
		return 0
	}
	return q1 / q2
}

//水箱效率
func CalcWatertankEfficiency(hourStr string) float64 {
	q1 := 0.0
	var Tinitial float64
	Taver := 0.0
	h, ok := GetOpcFloatList("ZLZ.OUTPUT_P10", hourStr)
	if !ok {
		return 0
	}
	Tin, ok := GetOpcFloatList("ZLZ.OUTPUT_T3", hourStr)
	if !ok {
		return 0
	}
	Tout, ok := GetOpcFloatList("ZLZ.OUTPUT_T4", hourStr)
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

//能源站效率
func CalcEnergyEfficiency(hourStr string) float64 {
	totalQ1, ok := GetOpcFloatList("ZLZ.RLB_%E7%B4%AF%E8%AE%A1%E7%83%AD%E9%87%8F", hourStr) //RLB_累计热量
	if !ok {
		return 0
	}
	APGL2, ok := GetOpcFloatList("ZLZ.%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6APGL2", hourStr) //有功电度APGL2
	if !ok {
		return 0
	}
	AZ1, ok := GetOpcFloatList("ZLZ.%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6AZ_1", hourStr) //有功电度AZ_1
	if !ok {
		return 0
	}
	w := make([][]float64, 5)
	for i := 1; i <= 4; i++ {
		w[i], _ = GetOpcFloatList("ZLZ.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(i), hourStr) //功率采集
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

//计算这个小时的碳排
func CalcEnergyCarbonHour(hourStr string) float64 {
	APGL2, ok := GetOpcFloatList("ZLZ.%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6APGL2", hourStr) //有功电度APGL2
	if !ok {
		return 0
	}
	AZ1, ok := GetOpcFloatList("ZLZ.%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6AZ_1", hourStr) //有功电度AZ_1
	if !ok {
		return 0
	}

	q23 := 0.0
	for i := 1; i <= 4; i++ {
		w, ok := GetOpcFloatList("ZLZ.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(i), hourStr) //功率采集
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

//计算当天的碳排
func CalcEnergyCarbonDay(dayStr string) float64 {
	sum := 0.0
	for i := 0; i < 24; i++ {
		hourStr := fmt.Sprintf("%s %02d", dayStr, i)
		var result defs.CalculationResultFloat
		err = MongoResult.FindOne(context.TODO(), bson.D{{"time", hourStr}, {"name", defs.EnergyCarbonDay}}).Decode(&result)
		if err == nil {
			sum += result.Value
		}
	}
	return sum
}

//计算本月的碳排
func CalcEnergyCarbonMonth(monthStr string) float64 {
	l, ok := GetResultFloatList(defs.EnergyCarbonMonth, monthStr)
	if !ok {
		return 0
	}
	sum := 0.0
	for _, v := range l {
		sum += v
	}
	return sum
}

//能源站该小时锅炉产热对应的负载
func CalcEnergyPayloadHour(q1 float64) float64 {
	q2 := 4 * 4 * 1e6 * 3600
	return q1 / q2
}

//计算当天的负载
func CalcEnergyPayloadDay(dayStr string) float64 {
	sum := 0.0
	for i := 0; i < 24; i++ {
		hourStr := fmt.Sprintf("%s %02d", dayStr, i)
		var result defs.CalculationResultFloat
		err = MongoResult.FindOne(context.TODO(), bson.D{{"time", hourStr}, {"name", defs.EnergyBoilerPayloadDay}}).Decode(&result)
		if err == nil {
			sum += result.Value
		}
	}
	return sum / 24
}

//计算本月的负载
func CalcEnergyPayloadMonth(monthStr string, maxDay int) float64 {
	if maxDay == 0 {
		return 0
	}
	l, ok := GetResultFloatList(defs.EnergyBoilerPayloadMonth, monthStr)
	if !ok {
		return 0
	}
	sum := 0.0
	for _, v := range l {
		sum += v
	}
	return sum / float64(maxDay)
}
