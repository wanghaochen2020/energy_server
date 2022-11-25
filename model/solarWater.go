package model

import (
	"energy/defs"
	"energy/utils"
	"fmt"
	"math"
	"strconv"
	"time"
)

// 太阳能热水锅炉今日总耗能，每分钟
func CalcSolarWaterBoilerPowerConsumptionToday(t time.Time) float64 {
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	hourStr := fmt.Sprintf("%s %02d", dayStr, t.Hour())
	ans := 0.0
	cost, ok := GetResultFloatList(defs.SolarWaterBoilerPowerConsumptionDay, dayStr)
	if ok {
		minLen := utils.Min(len(cost), t.Hour()+1)
		for i := 0; i < minLen; i++ {
			ans += cost[i]
		}
	}
	ans += CalcSolarWaterBoilerPowerConsumptionHour(hourStr) //如果卡就删掉这一行，然后把这个函数做成一个小时调用一次
	return ans
}

func CalcSolarWaterHeatCollecterInT(data *defs.LouSolarWater) float64 {
	ans, err := strconv.ParseFloat(data.CollectHeat.LT, 64)
	if err != nil {
		return 0
	}
	return ans
}

func CalcSolarWaterHeatCollecterOutT(data *defs.LouSolarWater) float64 {
	ans, err := strconv.ParseFloat(data.CollectHeat.HT, 64)
	if err != nil {
		return 0
	}
	return ans
}

func CalcSolarWaterJRQT(data *defs.LouSolarWater) float64 {
	ans, err := strconv.ParseFloat(data.System.JRQ_T, 64)
	if err != nil {
		return 0
	}
	return ans
}

// 太阳能热水今日总集热量，每分钟
func CalcSolarWaterHeatCollectionToday(t time.Time) float64 {
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	hourStr := fmt.Sprintf("%s %02d", dayStr, t.Hour())
	ans := 0.0
	cost, ok := GetResultFloatList(defs.SolarWaterHeatCollectionDay, dayStr)
	if ok {
		minLen := utils.Min(len(cost), t.Hour()+1)
		for i := 0; i < minLen; i++ {
			ans += cost[i]
		}
	}
	ans += CalcSolarWaterHeatCollectionHour(hourStr) //如果卡就删掉这一行，然后把这个函数做成一个小时调用一次
	return ans
}

func CalcSolarWaterPumpRunningNum(data *defs.LouSolarWater) float64 {
	ans := 0.0
	if data.HRPump_1.Sta == "1" {
		ans++
	}
	if data.HRPump_2.Sta == "1" {
		ans++
	}
	if data.JRPump_1.Sta == "1" {
		ans++
	}
	if data.JRPump_2.Sta == "1" {
		ans++
	}
	return ans
}

func CalcSolarWaterHeatCollectionMin(data *defs.LouSolarWater) float64 {
	HT, err := strconv.ParseFloat(data.CollectHeat.HT, 64)
	if err != nil {
		return 0
	}
	LT, err := strconv.ParseFloat(data.CollectHeat.LT, 64)
	if err != nil {
		return 0
	}
	V := 0.0 //流量
	if data.JRPump_1.Sta != "1" {
		V += 0.5 //集热循环泵流量（注意这里时间差为分钟）
	}
	if data.JRPump_2.Sta != "1" {
		V += 0.5 //集热循环泵流量
	}
	return (HT - LT) * V * 4200000 //单位为J
}

func CalcSolarWaterHeatCollectionHour(hourStr string) float64 {
	l, ok := GetResultFloatList(defs.SolarWaterHeatCollectionHour, hourStr)
	if !ok {
		return 0
	}
	ans := 0.0
	for _, v := range l {
		ans += v
	}
	return ans
}

// 计算当天的集热
func CalcSolarWaterHeatCollectionDay(dayStr string) float64 {
	l, ok := GetResultFloatList(defs.SolarWaterHeatCollectionDay, dayStr)
	if !ok {
		return 0
	}
	sum := 0.0
	for _, v := range l {
		sum += v
	}
	return sum
}

// 计算本月的集热
func CalcSolarWaterHeatCollectionMonth(monthStr string) float64 {
	l, ok := GetResultFloatList(defs.SolarWaterHeatCollectionMonth, monthStr)
	if !ok {
		return 0
	}
	sum := 0.0
	for _, v := range l {
		sum += v
	}
	return sum
}

func CalcSolarWaterBoilerPowerConsumptionMin(data *defs.LouSolarWater) float64 {
	ans := 0.0
	if data.Heater_1_1.Sta == "1" {
		ans++
	}
	if data.Heater_1_2.Sta == "1" {
		ans++
	}
	if data.Heater_1_3.Sta == "1" {
		ans++
	}
	if data.Heater_1_4.Sta == "1" {
		ans++
	}
	if data.Heater_1_5.Sta == "1" {
		ans++
	}
	if data.Heater_2_1.Sta == "1" {
		ans++
	}
	if data.Heater_2_2.Sta == "1" {
		ans++
	}
	if data.Heater_2_3.Sta == "1" {
		ans++
	}
	if data.Heater_2_4.Sta == "1" {
		ans++
	}
	if data.Heater_2_5.Sta == "1" {
		ans++
	}
	if data.Heater_3_1.Sta == "1" {
		ans++
	}
	if data.Heater_3_2.Sta == "1" {
		ans++
	}
	if data.Heater_3_3.Sta == "1" {
		ans++
	}
	if data.Heater_3_4.Sta == "1" {
		ans++
	}
	if data.Heater_3_5.Sta == "1" {
		ans++
	}
	if data.Heater_4_1.Sta == "1" {
		ans++
	}
	if data.Heater_4_2.Sta == "1" {
		ans++
	}
	if data.Heater_4_3.Sta == "1" {
		ans++
	}
	if data.Heater_4_4.Sta == "1" {
		ans++
	}
	if data.Heater_4_5.Sta == "1" {
		ans++
	}
	ans *= 90 / 60 //单位kW·h，单个加热器功率90kW
	return ans
}

func CalcSolarWaterBoilerPowerConsumptionHour(hourStr string) float64 {
	l, ok := GetResultFloatList(defs.SolarWaterBoilerPowerConsumptionHour, hourStr)
	if !ok {
		return 0
	}
	ans := 0.0
	for _, v := range l {
		ans += v
	}
	return ans
}

func CalcSolarWaterHeatEfficiency(t time.Time, heatCol float64) float64 {
	I0 := []float64{0, 1405, 1394, 1378, 1353, 1334, 1316, 1308, 1315, 1330, 1350, 1372, 1392}
	p := 0.5                     //从数据库中拿到p，这个是临时的
	phi := 40.52 / 180 * math.Pi //纬度
	delta := 23.45 * math.Sin(2.0*math.Pi*(284+float64(t.Day()))/365) / 180 * math.Pi
	omiga := float64(t.Hour()-12) * 15 / 180 * math.Pi
	sinh := math.Sin(phi)*math.Sin(delta) + math.Cos(phi)*math.Cos(delta)*math.Cos(omiga)
	h := math.Asin(sinh) / math.Pi * 180
	m := math.Sin(1 / h)
	pM := math.Pow(p, m)
	I := I0[t.Month()] * sinh * (pM + (1-pM)/2.0/(1-1.4*math.Log(p))) / 1000 * 216 //S是光伏板总面积，单位m^2，未知
	return heatCol / I
}

func CalcSolarWaterGuaranteeRate(q1 float64, q2 float64) float64 {
	if q1 == 0 && q2 == 0 {
		return 0
	}
	return q1 / (q1 + q2)
}
