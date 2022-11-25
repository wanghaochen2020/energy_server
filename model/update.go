package model

import (
	"context"
	"energy/defs"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func loopTime(t time.Duration, callback func(time.Time)) {
	if t == 0 {
		return
	}
	for {
		now := time.Now().Local()
		nano := time.Duration(now.UnixNano())
		deltaT := nano/t*t + t - nano
		next := now.Add(deltaT)
		time.Sleep(deltaT) //前往下一个整点

		go callback(next)
	}
}

// 更新本月的数据，传入的是本月最后一天的t
func updateMonth(t time.Time) {
	month := int(t.Month())
	yearStr := fmt.Sprintf("%04d", t.Year())
	monthStr := fmt.Sprintf("%s/%02d", yearStr, month)

	data := CalcEnergyCarbonMonth(monthStr) //能源站碳排
	MongoUpdateList(yearStr, month, defs.EnergyCarbonYear, data)
	data = CalcEnergyPayloadMonth(monthStr) //能源站锅炉负载
	MongoUpdateList(yearStr, month, defs.EnergyBoilerPayloadYear, data)

	data = CalcColdCarbonMonth(monthStr) //制冷站碳排
	MongoUpdateList(yearStr, month, defs.ColdCarbonYear, data)

	data = CalcSolarWaterHeatCollectionMonth(monthStr) //太阳能热水集热量
	MongoUpdateList(yearStr, month, defs.SolarWaterHeatCollectionYear, data)
}

// 更新本日的数据
func updateDay(t time.Time) {
	day := t.Day()
	monthStr := fmt.Sprintf("%04d/%02d", t.Year(), t.Month())
	dayStr := fmt.Sprintf("%s/%02d", monthStr, day)

	data := CalcEnergyCarbonDay(dayStr) //能源站碳排
	MongoUpdateList(monthStr, day, defs.EnergyCarbonMonth, data)
	data = CalcEnergyPayloadDay(dayStr) //能源站锅炉负载
	MongoUpdateList(monthStr, day, defs.EnergyBoilerPayloadMonth, data)

	data = CalcColdCarbonDay(dayStr) //制冷站碳排
	MongoUpdateList(monthStr, day, defs.ColdCarbonMonth, data)

	//太阳能热水
	data = CalcSolarWaterHeatCollectionDay(dayStr) //集热量
	MongoUpdateList(monthStr, day, defs.SolarWaterHeatCollectionMonth, data)

	if t.Add(time.Hour*24).Month() != t.Month() {
		updateMonth(t)
	}
}

// 更新本小时的数据
func updateHour(t time.Time) {
	hour := t.Hour()
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	hourStr := fmt.Sprintf("%s %02d", dayStr, hour)
	var data float64

	data = CalcEnergyHeatStorageAndRelease(hourStr) //蓄放热量统计，正值蓄热，负值放热
	MongoUpdateList(dayStr, hour, defs.EnergyHeatStorageAndRelease, data)
	q1 := CalcEnergyBoilerHeatSupply(hourStr) //能源站锅炉供热量
	q2 := CalcEnergyBoilerEnergyCost(hourStr) //锅炉能耗(单位kW·h)
	MongoUpdateList(dayStr, hour, defs.EnergyBoilerEnergyCost, q2)
	data = CalcEnergyBoilerEfficiency(q1, q2) //能源站锅炉效率
	MongoUpdateList(dayStr, hour, defs.EnergyBoilerEfficiencyDay, data)
	data = CalcWatertankEfficiency(hourStr) //能源站蓄热水箱效率
	MongoUpdateList(dayStr, hour, defs.EnergyWatertankEfficiencyDay, data)
	data = CalcEnergyEfficiency(hourStr) //能源站效率
	MongoUpdateList(dayStr, hour, defs.EnergyEfficiencyDay, data)
	data = CalcEnergyCarbonHour(hourStr) //能源站碳排
	MongoUpdateList(dayStr, hour, defs.EnergyCarbonDay, data)
	data = CalcEnergyPayloadHour(q1) //能源站锅炉负载率
	MongoUpdateList(dayStr, hour, defs.EnergyBoilerPayloadDay, data)

	//制冷中心
	q1 = CalcColdEnergyCost(hourStr, defs.ColdMachine1)
	q2 = CalcColdEnergyCost(hourStr, defs.ColdMachine2)
	q3 := CalcColdEnergyCost(hourStr, defs.ColdMachine3)
	q := q1 + q2 + q3
	MongoUpdateList(dayStr, hour, defs.ColdEnergyCostDay, q) //耗能
	//制冷效率（流量没拿到）
	//碳排
	data = CalcColdCarbonHour(q)
	MongoUpdateList(dayStr, hour, defs.ColdCarbonDay, data)
	//负载率（流量没拿到）

	//二次泵站
	data = CalcPumpEnergyCostHour(hourStr) //耗能
	MongoUpdateList(dayStr, hour, defs.PumpEnergyCostDay, data)
	dataList := CalcPumpEHR(hourStr) //输热比
	MongoUpdateList(dayStr, hour, defs.PumpEHR1, dataList[0])
	MongoUpdateList(dayStr, hour, defs.PumpEHR2, dataList[1])

	//太阳能热水
	q1 = CalcSolarWaterHeatCollectionHour(hourStr) //集热量
	MongoUpdateList(dayStr, hour, defs.SolarWaterHeatCollectionDay, q1)
	data = CalcSolarWaterHeatEfficiency(t, q1) //集热效率
	MongoUpdateList(dayStr, hour, defs.SolarWaterHeatEfficiencyDay, data)
	q2 = CalcSolarWaterBoilerPowerConsumptionHour(hourStr) //电加热器耗电
	MongoUpdateList(dayStr, hour, defs.SolarWaterBoilerPowerConsumptionHour, q2)
	data = CalcSolarWaterGuaranteeRate(q1, q2) //保证率
	MongoUpdateList(dayStr, hour, defs.SolarWaterGuaranteeRateDay, data)

	if hour == 23 {
		updateDay(t)
	}
}

// 更新本分钟和上一分钟的数据
func updateMinute(t time.Time) {
	hourStr := fmt.Sprintf("%04d/%02d/%02d %02d", t.Year(), t.Month(), t.Day(), t.Hour())
	lastMinTime := t.Add(-time.Minute)
	lastMin := lastMinTime.Minute()
	lastMinHourStr := fmt.Sprintf("%04d/%02d/%02d %02d", lastMinTime.Year(), lastMinTime.Month(), lastMinTime.Day(), lastMinTime.Hour())
	lastMinStr := fmt.Sprintf("%s:%02d", lastMinHourStr, lastMinTime.Minute())

	//能源站
	data, _ := CalcEnergyOnlineRate(hourStr) //能源站设备在线率
	MongoUpsertOne(defs.EnergyOnlineRate, data)
	data = CalcEnergyBoilerPower(lastMinHourStr, lastMin) //能源站锅炉总功率
	MongoUpsertOne(defs.EnergyBoilerPower, data)
	data = CalcEnergyPowerConsumptionToday(lastMinTime) //能源站今日能耗
	MongoUpsertOne(defs.EnergyPowerConsumptionToday, data)
	data = CalcEnergyBoilerRunningNum(lastMinHourStr, lastMin) //能源站锅炉运行数目
	MongoUpsertOne(defs.EnergyBoilerRunningNum, data)
	data = CalcEnergyTankRunningNum(lastMinHourStr, lastMin) //蓄热水箱运行台数
	MongoUpsertOne(defs.EnergyTankRunningNum, data)

	//设备温度(无)
	//设备供热量（无）

	data = CalcEnergyHeatSupplyToday(t) //总供热量
	MongoUpsertOne(defs.EnergyHeatSupplyToday, data)

	//制冷中心
	q1 := CalcColdMachinePower(lastMinHourStr, lastMin, defs.ColdMachine1)
	q2 := CalcColdMachinePower(lastMinHourStr, lastMin, defs.ColdMachine2)
	q3 := CalcColdMachinePower(lastMinHourStr, lastMin, defs.ColdMachine3)
	q4 := CalcColdCabinetPower(lastMinHourStr, lastMin)
	data = q1 + q2 + q3 + q4 //总功率
	MongoUpsertOne(defs.ColdPowerMin, data)
	data = q1 + q2 + q3 //制冷机功率
	MongoUpsertOne(defs.ColdMachinePowerMin, data)
	data = CalcColdEnergyCostToday(lastMinTime) //今日耗能
	MongoUpsertOne(defs.ColdEnergyCostToday, data)
	data = CalcColdMachineRunningNum(lastMinHourStr, lastMin) //制冷机运行数目
	MongoUpsertOne(defs.ColdMachineRunningNum, data)
	data = CalcColdCoolingWaterInT(lastMinHourStr, lastMin) //冷却进水温度
	MongoUpsertOne(defs.ColdCoolingWaterInT, data)
	data = CalcColdCoolingWaterOutT(lastMinHourStr, lastMin) //冷却出水温度
	MongoUpsertOne(defs.ColdCoolingWaterOut, data)
	data = CalcColdRefrigeratedWaterInT(lastMinHourStr, lastMin) //冷冻进水温度
	MongoUpsertOne(defs.ColdRefrigeratedWaterInT, data)
	data = CalcColdRefrigeratedWaterOutT(lastMinHourStr, lastMin) //冷冻出水温度
	MongoUpsertOne(defs.ColdRefrigeratedWaterOut, data)

	//二次泵站
	data = CalcPumpPowerMin(lastMinHourStr, lastMin) //总功率
	MongoUpsertOne(defs.PumpPowerMin, data)
	data = CalcPumpEnergyCostToday(lastMinTime) //今日耗电量
	MongoUpsertOne(defs.PumpPowerToday, data)

	//太阳能热水
	var GAData defs.LouSolarWater
	err := MongoLoukong.FindOne(context.TODO(), bson.D{{"time", lastMinStr}, {"name", "GA"}}).Decode(&GAData)
	if err == nil {
		data = CalcSolarWaterBoilerPowerConsumptionToday(t) //电加热器今日总耗电量
		MongoUpsertOne(defs.SolarWaterBoilerPowerConsumptionToday, data)
		data = CalcSolarWaterHeatCollecterInT(&GAData) //集热器进口温度
		MongoUpsertOne(defs.SolarWaterHeatCollecterInT, data)
		data = CalcSolarWaterHeatCollecterOutT(&GAData) //集热器出口温度
		MongoUpsertOne(defs.SolarWaterHeatCollecterOutT, data)
		data = CalcSolarWaterJRQT(&GAData) //锅炉温度
		MongoUpsertOne(defs.SolarWaterJRQT, data)
		data = CalcSolarWaterHeatCollectionMin(&GAData) //集热量
		MongoUpdateList(lastMinHourStr, lastMin, defs.SolarWaterHeatCollectionHour, data)
		data = CalcSolarWaterHeatCollectionToday(t) //今日总集热量
		MongoUpsertOne(defs.SolarWaterHeatCollectionToday, data)
		data = CalcSolarWaterPumpRunningNum(&GAData) //水泵运行数目
		MongoUpsertOne(defs.SolarWaterPumpRunningNum, data)
		data = CalcSolarWaterBoilerPowerConsumptionMin(&GAData) //电加热器耗电
		MongoUpdateList(lastMinHourStr, lastMin, defs.SolarWaterBoilerPowerConsumptionHour, data)
	}

	if lastMinTime.Minute() == 59 {
		updateHour(lastMinTime)
	}
}

// 定时更新
func LoopQueryUpdate() {
	go loopTime(time.Minute, updateMinute)
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}

func UpdateData(t1 time.Time, t2 time.Time) {

}
