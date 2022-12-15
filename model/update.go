package model

import (
	"context"
	"energy/defs"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func loopTime(t time.Duration, callback func(time.Time, bool)) {
	if t == 0 {
		return
	}
	for {
		now := time.Now().Local()
		nano := time.Duration(now.UnixNano())
		deltaT := nano/t*t + t - nano
		next := now.Add(deltaT)
		time.Sleep(deltaT) //前往下一个整点

		go callback(next, true)
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

	//二次泵站
	data = CalcPumpCarbonMonth(monthStr)
	MongoUpdateList(yearStr, month, defs.PumpCarbonYear, data)

	data = CalcSolarWaterHeatCollectionMonth(monthStr) //太阳能热水集热量
	MongoUpdateList(yearStr, month, defs.SolarWaterHeatCollectionYear, data)
	data = CalcSolarWaterHeatEfficiencyMonth(monthStr) //集热效率
	MongoUpdateList(yearStr, month, defs.SolarWaterHeatEfficiencyYear, data)
	data = CalcSolarWaterGuaranteeRateMonth(monthStr) //保证率
	MongoUpdateList(yearStr, month, defs.SolarWaterGuaranteeRateYear, data)
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

	//二次泵站
	data = CalcPumpCarbonDay(dayStr)
	MongoUpdateList(monthStr, day, defs.PumpCarbonMonth, data)

	//太阳能热水
	data = CalcSolarWaterHeatCollectionDay(dayStr) //集热量
	MongoUpdateList(monthStr, day, defs.SolarWaterHeatCollectionMonth, data)
	data = CalcSolarWaterHeatEfficiencyDay(dayStr) //集热效率
	MongoUpdateList(monthStr, day, defs.SolarWaterHeatEfficiencyMonth, data)
	data = CalcSolarWaterGuaranteeRateDay(dayStr) //保证率
	MongoUpdateList(monthStr, day, defs.SolarWaterGuaranteeRateMonth, data)

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

	q3 := CalcEnergyHeatStorageAndRelease(hourStr) //蓄放热量统计，正值蓄热，负值放热
	MongoUpdateList(dayStr, hour, defs.EnergyHeatStorageAndRelease, data)
	q1 := CalcEnergyBoilerHeatSupply(hourStr)     //能源站锅炉供热量
	q2List := CalcEnergyBoilerEnergyCost(hourStr) //各锅炉能耗(单位kW·h)
	q2 := q2List[0] + q2List[1] + q2List[2] + q2List[3]
	MongoUpdateList(dayStr, hour, defs.EnergyBoilerPowerConsumptionDay1, q2List[0])
	MongoUpdateList(dayStr, hour, defs.EnergyBoilerPowerConsumptionDay2, q2List[1])
	MongoUpdateList(dayStr, hour, defs.EnergyBoilerPowerConsumptionDay3, q2List[2])
	MongoUpdateList(dayStr, hour, defs.EnergyBoilerPowerConsumptionDay4, q2List[3])
	MongoUpdateList(dayStr, hour, defs.EnergyBoilerEnergyCost, q2)
	data = CalcEnergyBoilerEfficiency(q1, q2) //能源站锅炉效率
	MongoUpdateList(dayStr, hour, defs.EnergyBoilerEfficiencyDay, data)
	data = CalcWatertankEfficiency(q3, hourStr) //能源站蓄热水箱效率
	MongoUpdateList(dayStr, hour, defs.EnergyWatertankEfficiencyDay, data)
	data = CalcEnergyEfficiency(hourStr) //能源站效率
	MongoUpdateList(dayStr, hour, defs.EnergyEfficiencyDay, data)
	data = CalcEnergyCarbonHour(hourStr, q2) //能源站碳排
	MongoUpdateList(dayStr, hour, defs.EnergyCarbonDay, data)
	data = CalcEnergyPayloadHour(q1) //能源站锅炉负载率
	MongoUpdateList(dayStr, hour, defs.EnergyBoilerPayloadDay, data)
	dataList := CalcEnergyRunningTimeHour(hourStr) //设备运行时间（分钟）
	for i := 0; i < 9; i++ {
		MongoUpdateList(dayStr, hour, defs.EnergyRunningTimeDay[i], dataList[i])
	}

	//制冷中心
	q1 = CalcColdEnergyCost(hourStr, defs.ColdMachine1)
	q2 = CalcColdEnergyCost(hourStr, defs.ColdMachine2)
	q3 = CalcColdEnergyCost(hourStr, defs.ColdMachine3)
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
	dataList = CalcPumpEHR(hourStr) //输热比
	MongoUpdateList(dayStr, hour, defs.PumpEHR1, dataList[0])
	MongoUpdateList(dayStr, hour, defs.PumpEHR2, dataList[1])

	dataList = CalcHeatConsumptionHour(hourStr) //耗热统计
	MongoUpdateList(dayStr, hour, defs.GroupHeatConsumptionDay1, dataList[0])
	MongoUpdateList(dayStr, hour, defs.GroupHeatConsumptionDay2, dataList[1])
	MongoUpdateList(dayStr, hour, defs.GroupHeatConsumptionDay3, dataList[2])
	MongoUpdateList(dayStr, hour, defs.GroupHeatConsumptionDay4, dataList[3])
	MongoUpdateList(dayStr, hour, defs.GroupHeatConsumptionDay5, dataList[4])
	MongoUpdateList(dayStr, hour, defs.GroupHeatConsumptionDay6, dataList[5])
	MongoUpdateList(dayStr, hour, defs.GroupHeatConsumptionDayPubS, dataList[6])

	//太阳能热水
	q1 = CalcSolarWaterHeatCollectionHour(hourStr) //集热量
	MongoUpdateList(dayStr, hour, defs.SolarWaterHeatCollectionDay, q1)
	data = CalcSolarWaterHeatEfficiency(t, q1) //集热效率
	MongoUpdateList(dayStr, hour, defs.SolarWaterHeatEfficiencyDay, data)
	q2 = CalcSolarWaterBoilerPowerConsumptionHour(hourStr) //电加热器耗电
	MongoUpdateList(dayStr, hour, defs.SolarWaterBoilerPowerConsumptionDay, q2)
	data = CalcSolarWaterGuaranteeRate(q1, q2) //保证率
	MongoUpdateList(dayStr, hour, defs.SolarWaterGuaranteeRateDay, data)

	if hour == 23 {
		updateDay(t)
	}
}

// 更新上一分钟的数据。如果仅仅是导入过去数据upsert设为false，需要更新页面设为true
func updateMinute(t time.Time, upsert bool) {
	lastMinTime := t.Add(-time.Minute)
	lastMin := lastMinTime.Minute()
	lastMinDayStr := fmt.Sprintf("%04d/%02d/%02d", lastMinTime.Year(), lastMinTime.Month(), lastMinTime.Day())
	lastMinHourStr := fmt.Sprintf("%s %02d", lastMinDayStr, lastMinTime.Hour())
	// lastMinStr := fmt.Sprintf("%s:%02d", lastMinHourStr, lastMinTime.Minute())
	var data float64

	//报警数据
	UpdateEnergyAlarm(lastMinHourStr, lastMin, lastMinTime) //能源站
	UpdateColdAlarm(lastMinHourStr, lastMin, lastMinTime)   //制冷中心
	UpdatePumpAlarm(lastMinHourStr, lastMin, lastMinTime)   //二次泵站
	//之后计算要用的数据

	exampleTime := "2022/05/01 08:05"
	//二次泵站
	var HeatData defs.LouHeatList
	// err := MongoLoukong.FindOne(context.TODO(), bson.D{{"time", lastMinStr}, {"name", "heat"}}).Decode(&HeatData)
	err := MongoLoukong.FindOne(context.TODO(), bson.D{{"time", exampleTime}, {"name", "heat"}}).Decode(&HeatData)
	if err == nil {
		dataList := CalcPumpHeat(&HeatData) //统计输热量
		MongoUpdateList(lastMinHourStr, lastMin, defs.GroupHeatConsumptionHour1, dataList[0])
		MongoUpdateList(lastMinHourStr, lastMin, defs.GroupHeatConsumptionHour2, dataList[1])
		MongoUpdateList(lastMinHourStr, lastMin, defs.GroupHeatConsumptionHour3, dataList[2])
		MongoUpdateList(lastMinHourStr, lastMin, defs.GroupHeatConsumptionHour4, dataList[3])
		MongoUpdateList(lastMinHourStr, lastMin, defs.GroupHeatConsumptionHour5, dataList[4])
		MongoUpdateList(lastMinHourStr, lastMin, defs.GroupHeatConsumptionHour6, dataList[5])
		MongoUpdateList(lastMinHourStr, lastMin, defs.GroupHeatConsumptionHourPubS, dataList[6])
		MongoUpdateList(lastMinHourStr, lastMin, defs.PumpHeatHour1, dataList[0]+dataList[1])
		MongoUpdateList(lastMinHourStr, lastMin, defs.PumpHeatHour2, dataList[2]+dataList[3]+dataList[4]+dataList[5])
	}

	//太阳能热水
	var GAData defs.LouSolarWaterList
	// GAerr := MongoLoukong.FindOne(context.TODO(), bson.D{{"time", lastMinStr}, {"name", "GA"}}).Decode(&GAData)
	GAerr := MongoLoukong.FindOne(context.TODO(), bson.D{{"time", exampleTime}, {"name", "GA"}}).Decode(&GAData)
	if GAerr == nil {
		data = CalcSolarWaterHeatCollectionMin(&GAData.Info) //集热量
		MongoUpdateList(lastMinHourStr, lastMin, defs.SolarWaterHeatCollectionHour, data)
		data = CalcSolarWaterBoilerPowerConsumptionMin(&GAData.Info) //电加热器耗电
		MongoUpdateList(lastMinHourStr, lastMin, defs.SolarWaterBoilerPowerConsumptionHour, data)
	}

	//实时展示数据
	if upsert {
		//能源站
		data, _ = CalcEnergyOnlineRate(lastMinHourStr) //能源站设备在线率
		MongoUpsertOne(defs.EnergyOnlineRate, data)
		data = CalcEnergyBoilerPower(lastMinHourStr, lastMin) //能源站锅炉总功率
		MongoUpsertOne(defs.EnergyBoilerPower, data)
		dataList := CalcEnergyBoilerEnergyCost(lastMinHourStr) //本小时各锅炉能耗(单位kW·h)
		q23 := dataList[0] + dataList[1] + dataList[2] + dataList[3]
		dataList = CalcEnergyBoilerEnergyCostToday(lastMinDayStr, dataList) //今日各锅炉能耗
		MongoUpsertOne(defs.EnergyBoilerPowerConsumptionToday1, dataList[0])
		MongoUpsertOne(defs.EnergyBoilerPowerConsumptionToday2, dataList[1])
		MongoUpsertOne(defs.EnergyBoilerPowerConsumptionToday3, dataList[2])
		MongoUpsertOne(defs.EnergyBoilerPowerConsumptionToday4, dataList[3])
		data = CalcEnergyPowerConsumptionToday(lastMinTime, q23) //能源站今日能耗
		MongoUpsertOne(defs.EnergyPowerConsumptionToday, data)
		data = CalcEnergyBoilerRunningNum(lastMinHourStr, lastMin) //能源站锅炉运行数目
		MongoUpsertOne(defs.EnergyBoilerRunningNum, data)
		data = CalcEnergyTankRunningNum(lastMinHourStr, lastMin) //蓄热水箱运行台数
		MongoUpsertOne(defs.EnergyTankRunningNum, data)
		dataList = CalcEnergyRunningTimeToday(lastMinTime) //设备今日运行时长
		MongoUpsertOne(defs.EnergyRunningTimeToday, dataList)

		data = CalcEnergyHeatSupplyToday(t) //总供热量
		MongoUpsertOne(defs.EnergyHeatSupplyToday, data)

		//制冷中心
		q1 := CalcColdMachinePower(lastMinHourStr, lastMin, defs.ColdMachine1)
		MongoUpsertOne(defs.ColdMachinePowerMin1, q1)
		q2 := CalcColdMachinePower(lastMinHourStr, lastMin, defs.ColdMachine2)
		MongoUpsertOne(defs.ColdMachinePowerMin2, q2)
		q3 := CalcColdMachinePower(lastMinHourStr, lastMin, defs.ColdMachine3)
		MongoUpsertOne(defs.ColdMachinePowerMin3, q3)
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
		MongoUpsertOne(defs.ColdCoolingWaterOutT, data)
		data = CalcColdRefrigeratedWaterInT(lastMinHourStr, lastMin) //冷冻进水温度
		MongoUpsertOne(defs.ColdRefrigeratedWaterInT, data)
		data = CalcColdRefrigeratedWaterOutT(lastMinHourStr, lastMin) //冷冻出水温度
		MongoUpsertOne(defs.ColdRefrigeratedWaterOutT, data)

		//二次泵站
		data = CalcPumpPowerMin(lastMinHourStr, lastMin) //总功率
		MongoUpsertOne(defs.PumpPowerMin, data)
		data = CalcPumpEnergyCostToday(lastMinTime) //今日耗电量
		MongoUpsertOne(defs.PumpPowerToday, data)
		data = CalcPumpRunningState(lastMinHourStr, lastMin, 1) //泵运行状态
		MongoUpsertOne(defs.PumpRunningState1, data)
		data = CalcPumpRunningState(lastMinHourStr, lastMin, 2) //泵运行状态
		MongoUpsertOne(defs.PumpRunningState2, data)
		data = CalcPumpRunningState(lastMinHourStr, lastMin, 3) //泵运行状态
		MongoUpsertOne(defs.PumpRunningState3, data)
		data = CalcPumpRunningState(lastMinHourStr, lastMin, 4) //泵运行状态
		MongoUpsertOne(defs.PumpRunningState4, data)
		data = CalcPumpRunningState(lastMinHourStr, lastMin, 5) //泵运行状态
		MongoUpsertOne(defs.PumpRunningState5, data)
		data = CalcPumpRunningState(lastMinHourStr, lastMin, 6) //泵运行状态
		MongoUpsertOne(defs.PumpRunningState6, data)

		//太阳能热水
		if GAerr == nil {
			data = CalcSolarWaterBoilerPowerConsumptionToday(t) //电加热器今日总耗电量
			MongoUpsertOne(defs.SolarWaterBoilerPowerConsumptionToday, data)
			data = CalcSolarWaterHeatCollecterInT(&GAData.Info) //集热器进口温度
			MongoUpsertOne(defs.SolarWaterHeatCollecterInT, data)
			data = CalcSolarWaterHeatCollecterOutT(&GAData.Info) //集热器出口温度
			MongoUpsertOne(defs.SolarWaterHeatCollecterOutT, data)
			data = CalcSolarWaterJRQT(&GAData.Info) //锅炉温度
			MongoUpsertOne(defs.SolarWaterJRQT, data)
			data = CalcSolarWaterHeatCollectionToday(t) //今日总集热量
			MongoUpsertOne(defs.SolarWaterHeatCollectionToday, data)
			data = CalcSolarWaterPumpRunningNum(&GAData.Info) //水泵运行数目
			MongoUpsertOne(defs.SolarWaterPumpRunningNum, data)
		}
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

// 从t1到t2计算数据，每分钟更新一次
func UpdateData(t1 time.Time, t2 time.Time) {
	tend := t2.Add(time.Minute)
	for t := t1.Add(time.Minute); t.Before(tend); t = t.Add(time.Minute) {
		updateMinute(t, false)
	}
	log.Print("Update Complete")
}
