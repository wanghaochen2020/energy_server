package model

import (
	"context"
	"energy/defs"
	"energy/utils"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

//需要楼控数据

func rightSubLeft(l []float64) float64 {
	left := 0
	right := 0
	for i := 0; i < len(l); i++ {
		if utils.Zero(l[i]) {
			continue
		}
		left = i
		break
	}
	for i := len(l) - 1; i >= 0; i-- {
		if utils.Zero(l[i]) {
			continue
		}
		right = i
		break
	}
	return l[right] - l[left]
}

//当日每小时总耗电量
func CalcPumpPower(hourStr string) float64 {
	ans := 0.0
	l, ok := GetOpcFloatList("ZLZ.T%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6_APRS_1", hourStr)
	if ok {
		ans = rightSubLeft(l)
	}
	return ans
}

//输热比。目前共有2个环路：环路0(D1、D2组团)；环路1(D3~D6组团)
func CalcPumpEHR(hourStr string) []float64 {
	powerStr := [][]string{
		{"ZLZ.T%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6_BPRS1_1", "ZLZ.T%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6_BPRS1_2"},
		{"ZLZ.T%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6_BPRS3_1", "ZLZ.T%E6%9C%89%E5%8A%9F%E7%94%B5%E5%BA%A6_BPRS3_2"}}
	heatMap := map[string]int{
		"D1组团能量表": 0,
		"D2组团能量表": 0,
		"D3组团能量表": 1,
		"D4组团能量表": 1,
		"D5组团能量表": 1,
		"D6组团能量表": 1,
	}
	const L = 2
	var power [L]float64
	var heat [L]float64
	var data defs.LouHeatList
	for i := 0; i < L; i++ {
		for j := 0; j < len(powerStr[i]); j++ {
			l, ok := GetOpcFloatList(powerStr[i][j], hourStr)
			if ok {
				power[i] += rightSubLeft(l)
			}
		}
	}

	err := MongoLoukong.FindOne(context.TODO(), bson.D{{"time", hourStr}, {"name", "heat"}}).Decode(&data)
	if err != nil {
		return []float64{0, 0}
	}
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
		heat[heatMap[v.Name]] += (inT - OutT) * CF * 4200000 //单位为J
	}
	for i := 0; i <= L; i++ {
		if utils.Zero(heat[i]) {
			power[i] = 0
		} else {
			power[i] /= heat[i]
		}
	}
	return power[:]
}
