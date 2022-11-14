package model

import (
	"context"
	"energy/defs"
	"energy/utils"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func CalcColdMachineRunningNum(hourStr string, min int) int { //制冷机运行台数
	ans := 0
	l, ok := GetOpcBoolList("ZLZ.Z_L_%E5%90%AF%E5%81%9C", hourStr)
	if ok && len(l) > min {
		ans += utils.Bool2Int(l[min])
	}
	l, ok = GetOpcBoolList("ZLZ.Z_LX1_%E5%90%AF%E5%81%9C", hourStr)
	if ok && len(l) > min {
		ans += utils.Bool2Int(l[min])
	}
	l, ok = GetOpcBoolList("ZLZ.Z_LX2_%E5%90%AF%E5%81%9C", hourStr)
	if ok && len(l) > min {
		ans += utils.Bool2Int(l[min])
	}
	return ans
}

//计算给定小时和给定机器的能耗，machineId目前有Z_L,Z_LX1,Z_LX2，分别是螺杆机、离心机1、离心机2
func CalcColdPower(hourStr string, machineId string) float64 { //制冷站能效比:离心机COP≥6.431，螺杆机COP≥6.104；新闻中心螺杆机COP≥5.777
	ans := 0.0
	const UString = "AB%E7%9B%B8%E7%94%B5%E5%8E%8B"
	const IString = "L1%E7%9B%B8%E7%94%B5%E6%B5%81"
	//螺杆机
	U, _ := GetOpcFloatList(fmt.Sprintf("ZLZ.%s_%s", machineId, UString), hourStr)
	I, _ := GetOpcFloatList(fmt.Sprintf("ZLZ.%s_%s", machineId, IString), hourStr)
	minLen := utils.Min(len(U), len(I))
	for i := 0; i < minLen; i++ {
		ans += U[i] * I[i] * 6.104 * 3 / 60 //kW
	}
	return ans
}

func ColdEfficiency(hourStr string, machineId string) float64 {
	//差制冷量所需的流量数据
	return 0.0
}

//q单位为kW·h
func CalcColdCarbonHour(q float64) float64 {
	return q / 1000 * 0.604
}

func CalcColdCarbonDay(dayStr string) float64 {
	sum := 0.0
	for i := 0; i < 24; i++ {
		hourStr := fmt.Sprintf("%s %02d", dayStr, i)
		var result defs.CalculationResultFloat
		err = MongoResult.FindOne(context.TODO(), bson.D{{"time", hourStr}, {"name", defs.ColdCarbonDay}}).Decode(&result)
		if err == nil {
			sum += result.Value
		}
	}
	return sum
}

//计算本月的碳排
func CalcColdCarbonMonth(monthStr string, maxDay int) float64 {
	sum := 0.0
	for i := 1; i <= maxDay; i++ {
		var result defs.CalculationResultFloat
		err = MongoResult.FindOne(context.TODO(), bson.D{{"time", monthStr}, {"name", defs.ColdCarbonMonth}}).Decode(&result)
		if err == nil {
			sum += result.Value
		}
	}
	return sum
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
