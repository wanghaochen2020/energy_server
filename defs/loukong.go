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

//GA的数据
type LouSolarWater struct {
	CollectHeat LouSolarWaterCollectHeat `bson:"CollectHeat"`
	JRPump_1    LouSolarWaterJRPump      `bson:"JRPump_1"`
	JRPump_2    LouSolarWaterJRPump      `bson:"JRPump_2"`
}

//集热器温度
type LouSolarWaterCollectHeat struct {
	HT string `bson:"HT"`
	LT string `bson:"LT"`
}

//集热泵运行状态
type LouSolarWaterJRPump struct {
	Sta string `bson:"Sta"`
}
