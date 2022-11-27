package defs

type LouHeat struct {
	Name   string `bson:"name"`
	Status string `bson:"status"`
	CF     string `bson:"CF"`
	InT    string `bson:"inT"`
	OutT   string `bson:"outT"`
}

type LouHeatList struct {
	Time string    `bson:"time"`
	Info []LouHeat `bson:"info"`
}

type LouSolarWaterList struct {
	Time string        `bson:"time"`
	Info LouSolarWater `bson:"info"`
}

// GA的数据
type LouSolarWater struct {
	Heater_1_1  LouSolarWaterStatus      `bson:"Heater_1_1"`
	Heater_1_2  LouSolarWaterStatus      `bson:"Heater_1_2"`
	Heater_1_3  LouSolarWaterStatus      `bson:"Heater_1_3"`
	Heater_1_4  LouSolarWaterStatus      `bson:"Heater_1_4"`
	Heater_1_5  LouSolarWaterStatus      `bson:"Heater_1_5"`
	Heater_2_1  LouSolarWaterStatus      `bson:"Heater_2_1"`
	Heater_2_2  LouSolarWaterStatus      `bson:"Heater_2_2"`
	Heater_2_3  LouSolarWaterStatus      `bson:"Heater_2_3"`
	Heater_2_4  LouSolarWaterStatus      `bson:"Heater_2_4"`
	Heater_2_5  LouSolarWaterStatus      `bson:"Heater_2_5"`
	Heater_3_1  LouSolarWaterStatus      `bson:"Heater_3_1"`
	Heater_3_2  LouSolarWaterStatus      `bson:"Heater_3_2"`
	Heater_3_3  LouSolarWaterStatus      `bson:"Heater_3_3"`
	Heater_3_4  LouSolarWaterStatus      `bson:"Heater_3_4"`
	Heater_3_5  LouSolarWaterStatus      `bson:"Heater_3_5"`
	Heater_4_1  LouSolarWaterStatus      `bson:"Heater_4_1"`
	Heater_4_2  LouSolarWaterStatus      `bson:"Heater_4_2"`
	Heater_4_3  LouSolarWaterStatus      `bson:"Heater_4_3"`
	Heater_4_4  LouSolarWaterStatus      `bson:"Heater_4_4"`
	Heater_4_5  LouSolarWaterStatus      `bson:"Heater_4_5"`
	CollectHeat LouSolarWaterCollectHeat `bson:"CollectHeat"`
	HRPump_1    LouSolarWaterStatus      `bson:"HRPump_1"`
	HRPump_2    LouSolarWaterStatus      `bson:"HRPump_2"`
	JRPump_1    LouSolarWaterStatus      `bson:"JRPump_1"`
	JRPump_2    LouSolarWaterStatus      `bson:"JRPump_2"`
	System      LouSolarWaterSystem      `bson:"System"`
}

// 集热器温度
type LouSolarWaterCollectHeat struct {
	HT string `bson:"HT"`
	LT string `bson:"LT"`
}

// 运行状态
type LouSolarWaterStatus struct {
	Sta string `bson:"Sta"`
}

// GA系统信息
type LouSolarWaterSystem struct {
	JRQ_T string `bson:"JRQ_T"`
}
