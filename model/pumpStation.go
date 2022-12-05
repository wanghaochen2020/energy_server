package model

import (
	"context"
	"energy/defs"
	"energy/utils"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 总功率(kW)
func CalcPumpPowerMin(hourStr string, min int) float64 {
	ans := 0.0
	l, _ := GetOpcFloatList("ZLZ.T_RUN_P1", hourStr)
	if min < len(l) {
		ans += l[min] * 18.5
	}
	l, _ = GetOpcFloatList("ZLZ.T_RUN_P2", hourStr)
	if min < len(l) {
		ans += l[min] * 18.5
	}
	l, _ = GetOpcFloatList("ZLZ.T_RUN_P3", hourStr)
	if min < len(l) {
		ans += l[min] * 15
	}
	l, _ = GetOpcFloatList("ZLZ.T_RUN_P4", hourStr)
	if min < len(l) {
		ans += l[min] * 15
	}
	l, _ = GetOpcFloatList("ZLZ.T_RUN_P5", hourStr)
	if min < len(l) {
		ans += l[min] * 22
	}
	l, _ = GetOpcFloatList("ZLZ.T_RUN_P6", hourStr)
	if min < len(l) {
		ans += l[min] * 22
	}
	return ans
}

// 今日耗电量
func CalcPumpEnergyCostToday(t time.Time) float64 {
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	hourStr := fmt.Sprintf("%s %02d", dayStr, t.Hour())
	ans := 0.0
	l, ok := GetResultFloatList(defs.PumpEnergyCostDay, dayStr)
	if ok {
		minLen := utils.Min(len(l), t.Hour()+1)
		for i := 0; i < minLen; i++ {
			ans += l[i]
		}
	}
	ans += CalcPumpEnergyCostHour(hourStr) //如果卡就删掉这一行，然后把这个函数做成一个小时调用一次
	return ans
}

func CalcPumpRunningState(hourStr string, min int, id int) float64 {
	l, ok := GetOpcFloatList(fmt.Sprintf("ZLZ.T_RUN_P%d", id), hourStr)
	if !ok || len(l) <= min {
		return 0
	}
	return l[min]
}

// 当日每小时总耗电量
func CalcPumpEnergyCostHour(hourStr string) float64 {
	ans := 0.0
	l, ok := GetOpcFloatList("ZLZ.T%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6_APRS_1", hourStr)
	if ok {
		ans = RightSubLeft(l)
	}
	return ans
}

// 通过楼控数据计算当小时每分钟各组团耗热量，用于计算输热比和负荷预测
func CalcPumpHeat(data *defs.LouHeatList) []float64 {
	heatMap := map[string]int{
		"D1组团能量表": 0,
		"D2组团能量表": 1,
		"D3组团能量表": 2,
		"D4组团能量表": 3,
		"D5组团能量表": 4,
		"D6组团能量表": 5,
		"南区能量表":   6,
	}
	const L = 7
	var heat [L]float64
	for _, v := range data.Info {
		if v.Status != "0" {
			continue
		}
		inT, err := strconv.ParseFloat(v.InT, 64)
		if err != nil {
			continue
		}
		OutT, err := strconv.ParseFloat(v.OutT, 64)
		if err != nil {
			continue
		}
		CF, err := strconv.ParseFloat(v.CF, 64)
		if err != nil {
			continue
		}
		heat[heatMap[v.Name]] += (inT - OutT) * CF * 4200000 / 60 //单位为J
	}
	return heat[:]
}

// 输热比。目前共有2个环路：环路0(D1、D2组团)；环路1(D3~D6组团)
func CalcPumpEHR(hourStr string) []float64 {
	powerStr := [][]string{
		{"ZLZ.T%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6_BPRS1_1", "ZLZ.T%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6_BPRS1_2"},
		{"ZLZ.T%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6_BPRS3_1", "ZLZ.T%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6_BPRS3_2"}}
	const L = 2
	var power [L]float64
	var heat [L]float64
	for i := 0; i < L; i++ {
		for j := 0; j < len(powerStr[i]); j++ {
			l, ok := GetOpcFloatList(powerStr[i][j], hourStr)
			if ok {
				power[i] += RightSubLeft(l)
			}
		}
	}
	heat[0] = SumOpcResultList(defs.PumpHeatHour1, hourStr)
	heat[1] = SumOpcResultList(defs.PumpHeatHour2, hourStr)

	for i := 0; i < L; i++ {
		if utils.Zero(heat[i]) {
			power[i] = 0
		} else {
			power[i] /= heat[i]
		}
	}
	return power[:]
}

func CalcHeatConsumptionHour(hourStr string) []float64 {
	tables := []string{defs.GroupHeatConsumptionDay1, defs.GroupHeatConsumptionDay2, defs.GroupHeatConsumptionDay3,
		defs.GroupHeatConsumptionDay4, defs.GroupHeatConsumptionDay5, defs.GroupHeatConsumptionDay6, defs.GroupHeatConsumptionDayPubS}
	ans := make([]float64, len(tables))
	for i, v := range tables {
		ans[i] = SumOpcResultList(v, hourStr)
	}
	return ans
}

var pumpAlarmOpcList = map[string]defs.Alarm{
	"ZLZ.T_ALARM_P1": {"1#空调冷水二次泵", "报警"},
	"ZLZ.T_ALARM_P2": {"2#空调冷水二次泵", "报警"},
	"ZLZ.T_ALARM_P3": {"3#空调冷水二次泵", "报警"},
	"ZLZ.T_ALARM_P4": {"4#空调冷水二次泵", "报警"},
	"ZLZ.T_ALARM_P5": {"5#空调冷水二次泵", "报警"},
	"ZLZ.T_ALARM_P6": {"6#空调冷水二次泵", "报警"},
}

//报警
func UpdatePumpAlarm(hourStr string, min int, t time.Time) {
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
		alarmNum, _ = GetResultFloat(defs.PumpAlarmNumToday, dayStr)
	}

	minStr := fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
	var newList []defs.OpcAlarm

	for k, v := range pumpAlarmOpcList {
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

	MongoUpsertOne(defs.PumpAlarmNumToday, alarmNum)
}
