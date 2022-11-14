package model

import (
	"context"
	"energy/defs"
	"math"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func CalcSolarWaterHeatCollection(hourStr string) float64 {
	var data defs.LouSolarWater
	err := MongoLoukong.FindOne(context.TODO(), bson.D{{"time", hourStr}, {"name", "GA"}}).Decode(&data)
	if err != nil {
		return 0
	}
	HT, err := strconv.ParseFloat(data.CollectHeat.HT, 64)
	if err != nil {
		return 0
	}
	LT, err := strconv.ParseFloat(data.CollectHeat.LT, 64)
	if err != nil {
		return 0
	}
	V := 0.0 //流量
	if data.JRPump_1.Sta != "0" {
		V += 0 //集热循环泵流量，未知
	}
	if data.JRPump_2.Sta != "0" {
		V += 0 //集热循环泵流量，未知
	}
	return (HT - LT) * V * 4200000 //单位为J
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
	I := I0[t.Month()] * sinh * (pM + (1-pM)/2.0/(1-1.4*math.Log(p)))
	return heatCol / I
}
