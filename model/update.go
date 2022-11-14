package model

import (
	"energy/defs"
	"fmt"
	"log"
	"time"
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

//更新本月的数据，传入的是本月最后一天的t
func updateMonth(t time.Time) {
	month := int(t.Month())
	yearStr := fmt.Sprintf("%04d", t.Year())
	monthStr := fmt.Sprintf("%s/%02d", yearStr, month)
	dayNum := t.Day()

	data := CalcEnergyCarbonMonth(monthStr, dayNum)
	MongoUpdateList(yearStr, month, defs.EnergyCarbonYear, data)
	data = CalcEnergyPayloadMonth(monthStr, dayNum)
	MongoUpdateList(yearStr, month, defs.EnergyBoilerPayloadYear, data)

	data = CalcColdCarbonMonth(monthStr, dayNum)
	MongoUpdateList(yearStr, month, defs.ColdCarbonYear, data)
}

//更新本日的数据
func updateDay(t time.Time) {
	day := t.Day()
	monthStr := fmt.Sprintf("%04d/%02d", t.Year(), t.Month())
	dayStr := fmt.Sprintf("%s/%02d", monthStr, day)

	data := CalcEnergyCarbonDay(dayStr) //能源站碳排
	MongoUpdateList(monthStr, day, defs.EnergyCarbonMonth, data)
	data = CalcEnergyPayloadDay(dayStr) //能源站锅炉负载
	MongoUpdateList(monthStr, day, defs.EnergyBoilerPayloadMonth, data)

	data = CalcColdCarbonDay(dayStr) //能源站碳排
	MongoUpdateList(monthStr, day, defs.ColdCarbonMonth, data)

	if t.Add(time.Hour*24).Month() != t.Month() {
		updateMonth(t)
	}
}

//更新本小时的数据
func updateHour(t time.Time) {
	hour := t.Hour()
	dayStr := fmt.Sprintf("%04d/%02d/%02d", t.Year(), t.Month(), t.Day())
	hourStr := fmt.Sprintf("%s %02d", dayStr, hour)
	var data float64

	q1 := CalcEnergyBoilerHeatSupply(hourStr)      //能源站锅炉供热量
	data = CalcEnergyBoilerEfficiency(hourStr, q1) //能源站锅炉效率
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
	q1 = CalcColdPower(hourStr, defs.ColdMachine1)
	q2 := CalcColdPower(hourStr, defs.ColdMachine2)
	q3 := CalcColdPower(hourStr, defs.ColdMachine3)
	q := q1 + q2 + q3
	MongoUpdateList(dayStr, hour, defs.ColdPowerDay, q) //耗能
	//制冷效率（流量没拿到）
	//碳排
	data = CalcColdCarbonHour(q)
	MongoUpdateList(dayStr, hour, defs.ColdCarbonDay, data)
	//负载率（流量没拿到）

	//二次泵站
	data = CalcPumpPower(hourStr) //耗能
	MongoUpdateList(dayStr, hour, defs.PumpPowerDay, data)
	dataList := CalcPumpEHR(hourStr) //输热比
	MongoUpdateList(dayStr, hour, defs.PumpEHR1, dataList[0])
	MongoUpdateList(dayStr, hour, defs.PumpEHR2, dataList[1])

	//太阳能热水
	q1 = CalcSolarWaterHeatCollection(hourStr) //集热量
	MongoUpdateList(dayStr, hour, defs.SolarWaterHeatCollectionDay, q1)
	data = CalcSolarWaterHeatEfficiency(t, q1) //集热效率
	MongoUpdateList(dayStr, hour, defs.SolarWaterHeatEfficiencyDay, data)

	if hour == 23 {
		updateDay(t)
	}
}

//更新本分钟和上一分钟的数据
func updateMinute(t time.Time) {
	hourStr := fmt.Sprintf("%04d/%02d/%02d %02d", t.Year(), t.Month(), t.Day(), t.Hour())
	lastMinTime := t.Add(-time.Minute)
	lastMin := lastMinTime.Minute()
	lastMinHourStr := fmt.Sprintf("%04d/%02d/%02d %02d", lastMinTime.Year(), lastMinTime.Month(), lastMinTime.Day(), lastMinTime.Hour())

	data, _ := CalcEnergyOnlineRate(hourStr) //能源站设备在线率
	MongoUpsertOne(defs.EnergyOnlineRate, data)
	data = CalcEnergyBoilerPower(lastMinHourStr, lastMin) //能源站锅炉总功率
	MongoUpsertOne(defs.EnergyBoilerPower, data)
	data = CalcEnergyPowerConsumptionToday(lastMinTime) //能源站今日能耗
	MongoUpsertOne(defs.EnergyPowerConsumptionToday, data)
	data = CalcEnergyBoilerRunningNum(lastMinHourStr, lastMin) //能源站锅炉运行数目
	MongoUpsertOne(defs.EnergyBoilerRunningNum, data)
	//蓄热水箱运行台数
	//设备温度
	//设备供热量
	//总供热量
	//蓄放热量统计，正值蓄热，负值放热
	//锅炉耗电量统计

	if lastMinTime.Minute() == 59 {
		updateHour(lastMinTime)
	}
}

//定时更新
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
