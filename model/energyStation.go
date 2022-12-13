package model

import (
	"context"
	"energy/defs"
	"energy/utils"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MinParam struct {
	HourStr string
	Min     int
}

// 能源站在线率
func CalcEnergyOnlineRate(hourStr string) (float64, bool) {
	var result defs.MongoCountResult
	command := bson.D{{"count", "opc_data"}, {"query", bson.D{{"time", hourStr}}}}
	MongoOPC.Database().RunCommand(context.TODO(), command).Decode(&result)
	return float64(result.N) / 841, result.Ok
}

// 能源站锅炉总功率
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

// 能源站今日总耗能，每分钟
func CalcEnergyPowerConsumptionToday(t time.Time, q23 float64) float64 {
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	hourStr := fmt.Sprintf("%s %02d", dayStr, t.Hour())
	ans := 0.0
	cost, ok := GetResultFloatList(defs.EnergyCarbonDay, dayStr)
	if ok {
		minLen := utils.Min(len(cost), t.Hour()+1)
		for i := 0; i < minLen; i++ {
			ans += cost[i]
		}
	}
	ans += CalcEnergyCarbonHour(hourStr, q23) //如果卡就删掉这一行，然后把这个函数做成一个小时调用一次
	ans *= 1000 / 0.604                       //由tCO2换算到kW·h
	return ans
}

// 能源站锅炉运行台数
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

//设备该小时运行时间(分钟)
func CalcEnergyRunningTimeHour(hourStr string) []float64 {
	ans := make([]float64, 9)
	for i := 1; i <= 4; i++ {
		w, _ := GetOpcFloatList("ZLZ.%E7%B3%BB%E7%BB%9F%E8%BF%90%E8%A1%8C%E4%B8%AD"+fmt.Sprint(i), hourStr) //运行状态
		l := len(w)
		for j := 0; j < l; j++ {
			ans[i-1] += w[j]
		}
	}
	for i := 1; i <= 3; i++ {
		w, _ := GetOpcFloatList("ZLZ.RUN_P"+fmt.Sprint(i+3), hourStr) //运行状态
		l := len(w)
		for j := 0; j < l; j++ {
			ans[i+3] += w[j]
		}
	}

	w1, _ := GetOpcFloatList("ZLZ.OPEN_V6", hourStr)  //运行状态
	w2, _ := GetOpcFloatList("ZLZ.CLOSE_V1", hourStr) //运行状态
	l := utils.Min(len(w1), len(w2))
	for j := 0; j < l; j++ {
		if w1[j] != 0 && w2[j] != 0 {
			ans[7]++
		}
	}
	ans[8] = ans[7]
	return ans
}

//设备今日运行时间（小时）
func CalcEnergyRunningTimeToday(t time.Time) []float64 {
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	hourStr := fmt.Sprintf("%s %02d", dayStr, t.Hour())
	ans := CalcEnergyRunningTimeHour(hourStr)
	for i := 0; i < 9; i++ {
		ans[i] += SumOpcResultList(defs.EnergyRunningTimeDay[i], dayStr) / 60
	}
	return ans
}

// 能源站蓄热水箱运行台数
func CalcEnergyTankRunningNum(hourStr string, min int) float64 {
	items := []string{"ZLZ.OPEN_V6", "ZLZ.OPEN_V7"}
	ans := 0.0
	for _, v := range items {
		w, _ := GetOpcFloatList(v, hourStr) //运行状态
		if len(w) < min {
			continue
		}
		ans += w[min]
	}
	return ans
}

// 今日总供热量,如果卡就把这个函数做成一个小时调用一次
func CalcEnergyHeatSupplyToday(t time.Time) float64 {
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	//从前往后查
	q1 := 0.0
	i0 := 0
	for i0 = 0; i0 <= t.Hour(); i0++ {
		l, ok := GetOpcFloatList("ZLZ.RLB_%E7%B4%AF%E8%AE%A1%E7%83%AD%E9%87%8F", fmt.Sprintf("%s %02d", dayStr, i0)) //累计热量
		if !ok {
			continue
		}
		for _, v := range l {
			if !utils.Zero(v) {
				q1 = v
				break
			}
		}
		if !utils.Zero(q1) {
			break
		}
	}
	//从后往前查
	q2 := 0.0
	for i := t.Hour(); i >= i0; i-- {
		l, ok := GetOpcFloatList("ZLZ.RLB_%E7%B4%AF%E8%AE%A1%E7%83%AD%E9%87%8F", fmt.Sprintf("%s %02d", dayStr, i)) //累计热量
		if !ok {
			continue
		}
		ll := len(l)
		for j := ll - 1; j >= 0; j-- {
			if !utils.Zero(l[j]) {
				q2 = l[j]
				break
			}
		}
		if !utils.Zero(q2) {
			break
		}
	}
	return q2 - q1
}

// 能源站每日各小时水箱蓄放热量(正值蓄热，负值放热)
func CalcEnergyHeatStorageAndRelease(hourStr string) float64 {
	//通过小时初末的水箱（温度差*体积*密度*比热容）得出，体积变化不大视为定值（北建大意见）
	t1 := make([][]float64, 11)
	ok1 := make([]bool, 11)
	maxLen := 0
	for i := 0; i < 11; i++ {
		t1[i], ok1[i] = GetOpcFloatList(fmt.Sprintf("ZLZ.T1_%d", i+1), hourStr)
		maxLen = utils.Max(maxLen, len(t1[i]))
	}
	left := 0
	leftTemp := 0.0
	for left = 0; left < maxLen; left++ {
		leftTempNum := 0
		for i := 0; i < 11; i++ {
			if !ok1[i] || len(t1[i]) <= left || t1[i][left] == 0 {
				continue
			}
			leftTemp += t1[i][left]
			leftTempNum++
		}
		if leftTempNum != 0 {
			leftTemp /= float64(leftTempNum)
			break
		}
	}
	right := 0
	rightTemp := 0.0
	for right = maxLen; right > left; right-- {
		rightTempNum := 0
		for i := 0; i < 11; i++ {
			if !ok1[i] || len(t1[i]) <= right || t1[i][right] == 0 {
				continue
			}
			rightTemp += t1[i][left]
			rightTempNum++
		}
		if rightTempNum != 0 {
			rightTemp /= float64(rightTempNum)
			break
		}
	}
	if right == left {
		return 0
	}
	return (rightTemp - leftTemp) * 5.6 * 100 * 1e3 * 4.2 * 1e3 //单位J
}

// 能源站各锅炉能耗
func CalcEnergyBoilerEnergyCost(hourStr string) []float64 {
	q2 := make([]float64, 4)
	for i := 1; i <= 4; i++ {
		w, ok := GetOpcFloatList("ZLZ.%E5%8A%9F%E7%8E%87%E9%87%87%E9%9B%86"+fmt.Sprint(i), hourStr) //功率采集
		if !ok {
			continue
		}
		minLen := len(w)
		for j := 0; j < minLen; j++ {
			q2[i-1] += w[j] / 60
		}
	}
	return q2
}

var qlist = []string{
	defs.EnergyBoilerPowerConsumptionDay1,
	defs.EnergyBoilerPowerConsumptionDay2,
	defs.EnergyBoilerPowerConsumptionDay3,
	defs.EnergyBoilerPowerConsumptionDay4,
}

// 能源站各锅炉今日能耗
func CalcEnergyBoilerEnergyCostToday(dayStr string, q []float64) []float64 {
	q2 := make([]float64, 4)
	for i, v := range qlist {
		q2[i] = SumOpcResultList(v, dayStr)
		if len(q) > i {
			q2[i] += q[i]
		}
	}
	return q2
}

// 电极锅炉供热量
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

// 锅炉效率
func CalcEnergyBoilerEfficiency(q1 float64, q2 float64) float64 {
	q2 *= 3600000 //单位kW->J
	if q2 == 0 {
		return 0
	}
	return q1 / q2
}

// 水箱效率
func CalcWatertankEfficiency(q1 float64, hourStr string) float64 {
	Tinitial := 0.0
	Taver := 0.0
	Tout, ok := GetOpcFloatList("ZLZ.OUTPUT_T3", hourStr) //供水温度
	if !ok {
		return 0
	}
	averCount := 0
	for _, t := range Tout {
		if !utils.Zero(t) {
			if utils.Zero(Tinitial) {
				Tinitial = t
			}
			Taver += t
			averCount++
		}
	}
	if averCount == 0 {
		return 0
	}
	Taver /= float64(averCount)
	q2 := (Taver - Tinitial) * 5.6 * 100 * 1e3 * 4.2 * 1e3
	if q2 == 0 {
		return 0
	}
	return q1 / q2
}

// 能源站效率
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

// 计算这个小时的碳排
func CalcEnergyCarbonHour(hourStr string, q23 float64) float64 {
	APGL2, ok := GetOpcFloatList("ZLZ.%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6APGL2", hourStr) //有功电度APGL2
	if !ok {
		return 0
	}
	AZ1, ok := GetOpcFloatList("ZLZ.%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6AZ_1", hourStr) //有功电度AZ_1
	if !ok {
		return 0
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
		q22 = AZ1[r2] - AZ1[l2]
	}
	t := (q21 + q22 + q23) / 1000 * 0.604 //吨CO2
	return t
}

// 计算当天的碳排
func CalcEnergyCarbonDay(dayStr string) float64 {
	return SumOpcResultList(defs.EnergyCarbonDay, dayStr)
}

// 计算本月的碳排
func CalcEnergyCarbonMonth(monthStr string) float64 {
	return SumOpcResultList(defs.EnergyCarbonMonth, monthStr)
}

// 能源站该小时锅炉产热对应的负载
func CalcEnergyPayloadHour(q1 float64) float64 {
	q2 := 4 * 4 * 1e6 * 3600
	return q1 / q2
}

// 计算当天的负载
func CalcEnergyPayloadDay(dayStr string) float64 {
	return AvgOpcResultList(defs.EnergyBoilerPayloadDay, dayStr)
}

// 计算本月的负载
func CalcEnergyPayloadMonth(monthStr string) float64 {
	return AvgOpcResultList(defs.EnergyBoilerPayloadMonth, monthStr)
}

var energyAlarmOpcList = map[string]defs.Alarm{
	"ZLZ.ALARM_P1": {"蓄热循环泵1", "报警"}, "ZLZ.ALARM_P2": {"蓄热循环泵2", "报警"},
	"ZLZ.ALARM_P3": {"蓄热循环泵3", "报警"}, "ZLZ.ALARM_P4": {"放热循环泵1", "报警"},
	"ZLZ.ALARM_P5": {"放热循环泵2", "报警"}, "ZLZ.ALARM_P6": {"放热循环泵3", "报警"},
	"ZLZ.ALARM_P7": {"供热水泵1", "报警"}, "ZLZ.ALARM_P8": {"供热水泵2", "报警"},
	"ZLZ.ALARM_P9": {"供热水泵3", "报警"}, "ZLZ.ALARM_P10": {"蓄热系统补水泵1", "报警"},
	"ZLZ.ALARM_P11": {"蓄热系统补水泵2", "报警"},

	"ZLZ.%E5%8A%A0%E8%8D%AF%E6%B3%B5%E8%B7%B3%E9%97%B81": {"1#锅炉", "加药泵跳闸"},
	"ZLZ.%E5%8A%A0%E8%8D%AF%E6%B3%B5%E8%B7%B3%E9%97%B82": {"2#锅炉", "加药泵跳闸"},
	"ZLZ.%E5%8A%A0%E8%8D%AF%E6%B3%B5%E8%B7%B3%E9%97%B83": {"3#锅炉", "加药泵跳闸"},
	"ZLZ.%E5%8A%A0%E8%8D%AF%E6%B3%B5%E8%B7%B3%E9%97%B84": {"4#锅炉", "加药泵跳闸"},

	"ZLZ.%E6%8D%A2%E7%83%AD%E7%B3%BB%E7%BB%9F%E6%9C%AA%E5%90%AF%E5%8A%A8%E4%BF%9D%E6%8A%A41": {"1#锅炉", "换热系统未启动保护"},
	"ZLZ.%E6%8D%A2%E7%83%AD%E7%B3%BB%E7%BB%9F%E6%9C%AA%E5%90%AF%E5%8A%A8%E4%BF%9D%E6%8A%A42": {"2#锅炉", "换热系统未启动保护"},
	"ZLZ.%E6%8D%A2%E7%83%AD%E7%B3%BB%E7%BB%9F%E6%9C%AA%E5%90%AF%E5%8A%A8%E4%BF%9D%E6%8A%A43": {"3#锅炉", "换热系统未启动保护"},
	"ZLZ.%E6%8D%A2%E7%83%AD%E7%B3%BB%E7%BB%9F%E6%9C%AA%E5%90%AF%E5%8A%A8%E4%BF%9D%E6%8A%A44": {"4#锅炉", "换热系统未启动保护"},

	"ZLZ.%E8%BF%90%E8%A1%8C%E4%B8%AD%E6%B0%B4%E6%B3%B5%E6%95%85%E9%9A%9C%E6%88%96%E5%81%9C%E6%AD%A21": {"1#锅炉", "运行中水泵故障或停止"},
	"ZLZ.%E8%BF%90%E8%A1%8C%E4%B8%AD%E6%B0%B4%E6%B3%B5%E6%95%85%E9%9A%9C%E6%88%96%E5%81%9C%E6%AD%A22": {"2#锅炉", "运行中水泵故障或停止"},
	"ZLZ.%E8%BF%90%E8%A1%8C%E4%B8%AD%E6%B0%B4%E6%B3%B5%E6%95%85%E9%9A%9C%E6%88%96%E5%81%9C%E6%AD%A23": {"3#锅炉", "运行中水泵故障或停止"},
	"ZLZ.%E8%BF%90%E8%A1%8C%E4%B8%AD%E6%B0%B4%E6%B3%B5%E6%95%85%E9%9A%9C%E6%88%96%E5%81%9C%E6%AD%A24": {"4#锅炉", "运行中水泵故障或停止"},

	"ZLZ.%E9%94%85%E7%82%89%E5%86%85%E9%83%A8%E6%80%A5%E5%81%9C1": {"1#锅炉", "内部急停"},
	"ZLZ.%E9%94%85%E7%82%89%E5%86%85%E9%83%A8%E6%80%A5%E5%81%9C2": {"2#锅炉", "内部急停"},
	"ZLZ.%E9%94%85%E7%82%89%E5%86%85%E9%83%A8%E6%80%A5%E5%81%9C3": {"3#锅炉", "内部急停"},
	"ZLZ.%E9%94%85%E7%82%89%E5%86%85%E9%83%A8%E6%80%A5%E5%81%9C4": {"4#锅炉", "内部急停"},

	"ZLZ.%E9%94%85%E7%82%89%E5%87%BA%E6%B0%B4%E5%8E%8B%E9%AB%98%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%891": {"1#锅炉", "出水压高保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%87%BA%E6%B0%B4%E5%8E%8B%E9%AB%98%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%892": {"2#锅炉", "出水压高保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%87%BA%E6%B0%B4%E5%8E%8B%E9%AB%98%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%893": {"3#锅炉", "出水压高保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%87%BA%E6%B0%B4%E5%8E%8B%E9%AB%98%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%894": {"4#锅炉", "出水压高保护停炉"},

	"ZLZ.%E9%94%85%E7%82%89%E5%87%BA%E6%B0%B4%E6%B8%A9%E9%AB%98%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%891": {"1#锅炉", "出水压温保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%87%BA%E6%B0%B4%E6%B8%A9%E9%AB%98%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%892": {"2#锅炉", "出水压温保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%87%BA%E6%B0%B4%E6%B8%A9%E9%AB%98%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%893": {"3#锅炉", "出水压温保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%87%BA%E6%B0%B4%E6%B8%A9%E9%AB%98%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%894": {"4#锅炉", "出水压温保护停炉"},

	"ZLZ.%E9%94%85%E7%82%89%E5%8A%9F%E7%8E%87%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%891": {"1#锅炉", "功率保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%8A%9F%E7%8E%87%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%892": {"2#锅炉", "功率保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%8A%9F%E7%8E%87%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%893": {"3#锅炉", "功率保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%8A%9F%E7%8E%87%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%894": {"4#锅炉", "功率保护停炉"},

	"ZLZ.%E9%94%85%E7%82%89%E5%9B%9E%E6%B0%B4%E5%8E%8B%E4%BD%8E%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%891": {"1#锅炉", "回水压低保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%9B%9E%E6%B0%B4%E5%8E%8B%E4%BD%8E%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%892": {"2#锅炉", "回水压低保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%9B%9E%E6%B0%B4%E5%8E%8B%E4%BD%8E%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%893": {"3#锅炉", "回水压低保护停炉"},
	"ZLZ.%E9%94%85%E7%82%89%E5%9B%9E%E6%B0%B4%E5%8E%8B%E4%BD%8E%E4%BF%9D%E6%8A%A4%E5%81%9C%E7%82%894": {"4#锅炉", "回水压低保护停炉"},

	"ZLZ.%E9%94%85%E7%82%89%E5%A4%96%E9%83%A8%E6%80%A5%E5%81%9C1": {"1#锅炉", "外部急停"},
	"ZLZ.%E9%94%85%E7%82%89%E5%A4%96%E9%83%A8%E6%80%A5%E5%81%9C2": {"2#锅炉", "外部急停"},
	"ZLZ.%E9%94%85%E7%82%89%E5%A4%96%E9%83%A8%E6%80%A5%E5%81%9C3": {"3#锅炉", "外部急停"},
	"ZLZ.%E9%94%85%E7%82%89%E5%A4%96%E9%83%A8%E6%80%A5%E5%81%9C4": {"4#锅炉", "外部急停"},

	"ZLZ.%E9%94%85%E7%82%89%E5%BE%AA%E7%8E%AF%E6%B3%B5A%E8%B7%B3%E9%97%B81": {"1#锅炉", "循环泵A跳闸"},
	"ZLZ.%E9%94%85%E7%82%89%E5%BE%AA%E7%8E%AF%E6%B3%B5A%E8%B7%B3%E9%97%B82": {"2#锅炉", "循环泵A跳闸"},
	"ZLZ.%E9%94%85%E7%82%89%E5%BE%AA%E7%8E%AF%E6%B3%B5A%E8%B7%B3%E9%97%B83": {"3#锅炉", "循环泵A跳闸"},
	"ZLZ.%E9%94%85%E7%82%89%E5%BE%AA%E7%8E%AF%E6%B3%B5A%E8%B7%B3%E9%97%B84": {"4#锅炉", "循环泵A跳闸"},

	"ZLZ.%E9%94%85%E7%82%89%E5%BE%AA%E7%8E%AF%E6%B3%B5B%E8%B7%B3%E9%97%B81": {"1#锅炉", "循环泵B跳闸"},
	"ZLZ.%E9%94%85%E7%82%89%E5%BE%AA%E7%8E%AF%E6%B3%B5B%E8%B7%B3%E9%97%B82": {"2#锅炉", "循环泵B跳闸"},
	"ZLZ.%E9%94%85%E7%82%89%E5%BE%AA%E7%8E%AF%E6%B3%B5B%E8%B7%B3%E9%97%B83": {"3#锅炉", "循环泵B跳闸"},
	"ZLZ.%E9%94%85%E7%82%89%E5%BE%AA%E7%8E%AF%E6%B3%B5B%E8%B7%B3%E9%97%B84": {"4#锅炉", "循环泵B跳闸"},

	"ZLZ.%E9%94%85%E7%82%89%E6%89%A7%E8%A1%8C%E5%99%A8%E6%95%85%E9%9A%9C1": {"1#锅炉", "执行器故障"},
	"ZLZ.%E9%94%85%E7%82%89%E6%89%A7%E8%A1%8C%E5%99%A8%E6%95%85%E9%9A%9C2": {"2#锅炉", "执行器故障"},
	"ZLZ.%E9%94%85%E7%82%89%E6%89%A7%E8%A1%8C%E5%99%A8%E6%95%85%E9%9A%9C3": {"3#锅炉", "执行器故障"},
	"ZLZ.%E9%94%85%E7%82%89%E6%89%A7%E8%A1%8C%E5%99%A8%E6%95%85%E9%9A%9C4": {"4#锅炉", "执行器故障"},

	"ZLZ.%E9%94%85%E7%82%89%E6%89%A7%E8%A1%8C%E5%99%A8%E8%B7%B3%E9%97%B81": {"1#锅炉", "执行器跳闸"},
	"ZLZ.%E9%94%85%E7%82%89%E6%89%A7%E8%A1%8C%E5%99%A8%E8%B7%B3%E9%97%B82": {"2#锅炉", "执行器跳闸"},
	"ZLZ.%E9%94%85%E7%82%89%E6%89%A7%E8%A1%8C%E5%99%A8%E8%B7%B3%E9%97%B83": {"3#锅炉", "执行器跳闸"},
	"ZLZ.%E9%94%85%E7%82%89%E6%89%A7%E8%A1%8C%E5%99%A8%E8%B7%B3%E9%97%B84": {"4#锅炉", "执行器跳闸"},
}

//报警
func UpdateEnergyAlarm(hourStr string, min int, t time.Time) {
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	minStr := fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
	var oldList defs.MongoAlarmList
	MongoResult.FindOne(context.TODO(), bson.D{{"name", defs.EnergyAlarmToday}, {"time", dayStr}}).Decode(&oldList)
	alarmMap := make(map[string]int)
	for _, v := range oldList.Info {
		if v.State == 0 {
			alarmMap[v.Name] = 1
		}
	}

	alarmNum := 0.0
	if t.Hour() != 0 || t.Minute() != 0 {
		alarmNum, _ = GetResultFloat(defs.EnergyAlarmNumToday, dayStr)
	}

	var newList []defs.OpcAlarm

	for k, v := range energyAlarmOpcList {
		l, ok := GetOpcFloatList(k, hourStr)
		if !ok || len(l) <= min {
			continue
		}
		if !utils.Zero(l[min]) {
			if alarmMap[k] == 1 {
				continue
			} //新报警
			var newAlarm defs.OpcAlarm
			newAlarm.Name = v.Name
			newAlarm.State = 0
			newAlarm.Time = minStr
			newAlarm.Type = v.Type
			newList = append(newList, newAlarm)
			alarmNum++
		} else if alarmMap[k] == 1 { //已处理的旧报警
			alarmMap[k] = 2
		}
	}

	for _, v := range oldList.Info {
		if alarmMap[v.Name] == 2 {
			v.State = 1
		}
	}

	oldList.Info = append(oldList.Info, newList...)

	opts := options.Update().SetUpsert(true)
	_, err = MongoResult.UpdateOne(context.TODO(), bson.D{{"name", defs.EnergyAlarmToday}, {"time", hourStr}}, bson.D{{"$set", bson.D{{"info", oldList.Info}}}}, opts)
	if err != nil {
		log.Print(err)
	}

	MongoUpsertOne(defs.EnergyAlarmNumToday, alarmNum)
}
