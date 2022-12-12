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

// 当分钟制冷机功率，machineId目前有Z_L,Z_LX1,Z_LX2，分别是螺杆机、离心机1、离心机2
func CalcColdMachinePower(hourStr string, min int, machineId string) float64 {
	const UString = "AB%E7%9B%B8%E7%94%B5%E5%8E%8B"
	const IString = "L1%E7%9B%B8%E7%94%B5%E6%B5%81"
	COP := map[string]float64{defs.ColdMachine1: 6.104, defs.ColdMachine2: 6.431, defs.ColdMachine3: 6.431}
	u := 0.0
	i := 0.0
	U, _ := GetOpcFloatList(fmt.Sprintf("ZLZ.%s_%s", machineId, UString), hourStr)
	I, _ := GetOpcFloatList(fmt.Sprintf("ZLZ.%s_%s", machineId, IString), hourStr)
	if min < len(U) {
		u = U[min]
	}
	if min < len(I) {
		i = I[min]
	}
	return u * i * COP[machineId] * 3 / 1000 //单位kW
}

// 进线柜功率
func CalcColdCabinetPower(hourStr string, min int) float64 {
	U := 0.0
	I := 0.0
	l, _ := GetOpcFloatList("ZLZ.ZA%E7%9B%B8%E7%94%B5%E5%8E%8B_APLJ1_1", hourStr)
	if min < len(l) {
		U = l[min]
	}
	l, _ = GetOpcFloatList("ZLZ.ZA%E7%9B%B8%E7%94%B5%E6%B5%81_APLJ1_1", hourStr)
	if min < len(l) {
		I = l[min]
	}
	return U * I * 3 / 1000
}

// 今日能耗
func CalcColdEnergyCostToday(t time.Time) float64 {
	ans := 0.0
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	hourStr := fmt.Sprintf("%s %02d", dayStr, t.Hour())
	l, ok := GetResultFloatList(defs.ColdEnergyCostDay, dayStr)
	if ok {
		minLen := utils.Min(len(l), t.Hour()+1)
		for i := 0; i < minLen; i++ {
			ans += l[i]
		}
	}
	ans += CalcColdEnergyCost(hourStr, defs.ColdMachine1) //如果卡就删掉这一行，然后把这个函数做成一个小时调用一次
	ans += CalcColdEnergyCost(hourStr, defs.ColdMachine2) //如果卡就删掉这一行，然后把这个函数做成一个小时调用一次
	ans += CalcColdEnergyCost(hourStr, defs.ColdMachine3) //如果卡就删掉这一行，然后把这个函数做成一个小时调用一次
	return ans
}

// 制冷机运行台数
func CalcColdMachineRunningNum(hourStr string, min int) float64 {
	ans := 0.0
	l, ok := GetOpcBoolList("ZLZ.Z_L_%E5%90%AF%E5%81%9C", hourStr)
	if ok && len(l) > min {
		ans += utils.Bool2Float(l[min])
	}
	l, ok = GetOpcBoolList("ZLZ.Z_LX1_%E5%90%AF%E5%81%9C", hourStr)
	if ok && len(l) > min {
		ans += utils.Bool2Float(l[min])
	}
	l, ok = GetOpcBoolList("ZLZ.Z_LX2_%E5%90%AF%E5%81%9C", hourStr)
	if ok && len(l) > min {
		ans += utils.Bool2Float(l[min])
	}
	return ans
}

// 冷却水进水温度
func CalcColdCoolingWaterInT(hourStr string, min int) float64 {
	ans := 0.0
	num := 0
	l, _ := GetOpcFloatList("ZLZ.Z_L_%E5%86%B7%E5%8D%B4%E8%BF%9B%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_L_冷却进水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	l, _ = GetOpcFloatList("ZLZ.Z_LX1_%E5%86%B7%E5%8D%B4%E8%BF%9B%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_LX1_冷却进水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	l, _ = GetOpcFloatList("ZLZ.Z_LX2_%E5%86%B7%E5%8D%B4%E8%BF%9B%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_LX2_冷却进水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	if num == 0 {
		return 0
	}
	return ans / float64(num)
}

// 冷却水出水温度
func CalcColdCoolingWaterOutT(hourStr string, min int) float64 {
	ans := 0.0
	num := 0
	l, _ := GetOpcFloatList("ZLZ.Z_L_%E5%86%B7%E5%8D%B4%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_L_冷却出水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	l, _ = GetOpcFloatList("ZLZ.Z_LX1_%E5%86%B7%E5%8D%B4%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_LX1_冷却出水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	l, _ = GetOpcFloatList("ZLZ.Z_LX2_%E5%86%B7%E5%8D%B4%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_LX2_冷却出水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	if num == 0 {
		return 0
	}
	return ans / float64(num)
}

// 冷冻水进水温度
func CalcColdRefrigeratedWaterInT(hourStr string, min int) float64 {
	ans := 0.0
	num := 0
	l, _ := GetOpcFloatList("ZLZ.Z_L_%E5%86%B7%E5%86%BB%E8%BF%9B%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_L_冷冻进水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	l, _ = GetOpcFloatList("ZLZ.Z_LX1_%E5%86%B7%E5%86%BB%E8%BF%9B%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_LX1_冷冻进水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	l, _ = GetOpcFloatList("ZLZ.Z_LX2_%E5%86%B7%E5%86%BB%E8%BF%9B%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_LX2_冷冻进水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	if num == 0 {
		return 0
	}
	return ans / float64(num)
}

// 冷冻水出水温度
func CalcColdRefrigeratedWaterOutT(hourStr string, min int) float64 {
	ans := 0.0
	num := 0
	l, _ := GetOpcFloatList("ZLZ.Z_L_%E5%86%B7%E5%86%BB%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_L_冷冻出水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	l, _ = GetOpcFloatList("ZLZ.Z_LX1_%E5%86%B7%E5%86%BB%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_LX1_冷冻出水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	l, _ = GetOpcFloatList("ZLZ.Z_LX2_%E5%86%B7%E5%86%BB%E5%87%BA%E6%B0%B4%E6%B8%A9%E5%BA%A6", hourStr) //ZLZ.Z_LX2_冷冻出水温度
	if min < len(l) {
		ans += l[min]
		num++
	}
	if num == 0 {
		return 0
	}
	return ans / float64(num)
}

// 计算给定小时和给定机器的能耗，machineId目前有Z_L,Z_LX1,Z_LX2，分别是螺杆机、离心机1、离心机2
func CalcColdEnergyCost(hourStr string, machineId string) float64 { //制冷站能效比:离心机COP≥6.431，螺杆机COP≥6.104；新闻中心螺杆机COP≥5.777
	ans := 0.0
	const UString = "AB%E7%9B%B8%E7%94%B5%E5%8E%8B"
	const IString = "L1%E7%9B%B8%E7%94%B5%E6%B5%81"
	COP := map[string]float64{defs.ColdMachine1: 6.104, defs.ColdMachine2: 6.431, defs.ColdMachine3: 6.431}
	U, _ := GetOpcFloatList(fmt.Sprintf("ZLZ.%s_%s", machineId, UString), hourStr)
	I, _ := GetOpcFloatList(fmt.Sprintf("ZLZ.%s_%s", machineId, IString), hourStr)
	minLen := utils.Min(len(U), len(I))
	for i := 0; i < minLen; i++ {
		ans += U[i] * I[i] * COP[machineId] * 3 / 1000 //kW·h
	}
	return ans
}

func ColdEfficiency(hourStr string, machineId string) float64 {
	//差制冷量所需的流量数据
	return 0.0
}

// q单位为kW·h
func CalcColdCarbonHour(q float64) float64 {
	return q / 1000 * 0.604
}

// 计算本日的碳排
func CalcColdCarbonDay(dayStr string) float64 {
	return SumOpcResultList(defs.ColdCarbonDay, dayStr)
}

// 计算本月的碳排
func CalcColdCarbonMonth(monthStr string) float64 {
	return SumOpcResultList(defs.ColdCarbonMonth, monthStr)
}

func ColdPayloadDay(hourStr string) float64 {
	return 0.0
}

func ColdPayloadWeek(hourStr string) float64 {
	return 0.0
}

func ColdPayloadMonth(hourStr string) float64 {
	return 0.0
}

var coldAlarmOpcList = map[string]defs.Alarm{
	"ZLZ.Z_L_%E5%BD%93%E5%89%8D%E6%95%85%E9%9A%9C":                   {"螺杆机", "故障"},
	"ZLZ.Z_L_%E5%BD%93%E5%89%8D%E5%81%9C%E6%9C%BA%E6%95%85%E9%9A%9C": {"螺杆机", "停机故障"},
}

//报警
func UpdateColdAlarm(hourStr string, min int, t time.Time) {
	var oldList defs.MongoAlarmList
	MongoResult.FindOne(context.TODO(), bson.D{{"name", defs.ColdAlarmToday}, {"time", hourStr}}).Decode(&oldList)
	alarmMap := make(map[string]int)
	for _, v := range oldList.Info {
		if v.State == 0 {
			alarmMap[v.Name] = 1
		}
	}

	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())

	alarmNum := 0.0
	if t.Hour() != 0 || t.Minute() != 0 {
		alarmNum, _ = GetResultFloat(defs.ColdAlarmNumToday, dayStr)
	}

	minStr := fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
	var newList []defs.OpcAlarm

	for k, v := range coldAlarmOpcList {
		l, ok := GetOpcFloatList(k, hourStr)
		if !ok || len(l) <= min {
			continue
		}
		if !utils.Zero(l[min]) { //新报警
			var newAlarm defs.OpcAlarm
			newAlarm.Name = v.Name
			newAlarm.State = 0
			newAlarm.Time = minStr
			newAlarm.Type = v.Type
			newList = append(newList, newAlarm)
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
	_, err = MongoResult.UpdateOne(context.TODO(), bson.D{{"name", defs.ColdAlarmToday}, {"time", hourStr}}, bson.D{{"$set", bson.D{{"info", oldList.Info}}}}, opts)
	if err != nil {
		log.Print(err)
	}

	MongoUpsertOne(defs.ColdAlarmNumToday, alarmNum)
}
