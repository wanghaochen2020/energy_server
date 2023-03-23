package energyConfig

import (
	"energy/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var ()

//var loadDaily = [24]float64{537.41, 586.16, 618.91, 607.23, 608.5, 621.55, 645.52, 890.35, 690.17, 501.28, 1204.25, 915.07, 793.98, 748.76, 714.84, 694.95, 657.61, 681.41, 791.54, 999.22, 1156.91, 1264.27, 828.37, 661.38}
var loadWeekly = [7][24]float64{
	loadDaily,
	{331.61, 327.64, 291.98, 271.11, 338.62, 264.94, 439.08, 319.32, 171.60, 365.83, 368.43, 379.81, 297.32, 425.78, 428.64, 497.79, 318.52, 418.89, 565.24, 552.80, 496.60, 504.43, 519.91, 491.58},
	{454.70, 437.76, 389.20, 419.82, 381.12, 349.09, 530.82, 346.01, 169.92, 213.77, 120.06, 103.37, 131.56, 127.57, 144.52, 204.98, 58.52, 179.02, 317.90, 374.43, 483.58, 419.10, 388.01, 566.03},
	{517.56, 465.86, 506.98, 415.01, 432.81, 423.30, 558.53, 496.81, 279.86, 266.09, 194.44, 193.28, 253.58, 208.42, 204.77, 212.64, 125.83, 205.59, 380.75, 514.41, 444.45, 707.17, 396.03, 576.11},
	{478.42, 442.55, 498.55, 388.87, 456.45, 449.64, 570.63, 357.31, 153.06, 171.74, 219.67, 127.82, 171.08, 136.62, 130.15, 182.17, 58.87, 199.13, 319.11, 280.11, 383.18, 506.15, 305.88, 518.21},
	{484.52, 380.37, 386.27, 409.90, 319.67, 267.47, 417.12, 209.60, 82.29, 38.84, 70.02, 62.94, 80.29, 71.43, 109.92, 76.45, 41.19, 60.57, 62.54, 231.44, 147.66, 267.45, 263.59, 251.67},
	{206.54, 250.18, 214.85, 167.64, 182.05, 191.49, 211.57, 89.44, 27.73, 14.62, 7.68, 32.10, 32.35, 4.84, 33.30, 50.11, 37.97, 5.39, 22.92, 23.98, 87.57, 79.91, 89.96, 203.82}}
var energyWeekly = model.EnergyConfigWeekly{
	Qs: 29768,

	Heat_loss_rectify_coefficiency:       0.05,
	Heat_to_power_transform_coefficiency: 1.11,
	Heat_loss_coefficiency:               0.03,
	Carbon_emission_unit_power:           0.785,

	Vally_cost_time_start:  Vally_cost_time_start,
	Vally_cost_time_end:    Vally_cost_time_end,
	Flat_cost_time_1_start: Flat_cost_time_1_start,
	Flat_cost_time_1_end:   Flat_cost_time_1_end,
	Flat_cost_time_2_start: Flat_cost_time_2_start,
	Flat_cost_time_2_end:   Flat_cost_time_2_end,
	Flat_cost_time_3_start: Flat_cost_time_3_start,
	Flat_cost_time_3_end:   Flat_cost_time_3_end,
	Peak_cost_time_1_start: Peak_cost_time_1_start,
	Peak_cost_time_1_end:   Peak_cost_time_1_end,
	Peak_cost_time_2_start: Peak_cost_time_2_start,
	Peak_cost_time_2_end:   Peak_cost_time_2_end,

	Vally_cost: 0.26,
	Peak_cost:  1.24,
	Flat_cost:  0.73,

	Startup_1_boiler_lower_limiting_load_value: 400,
	Startup_2_boiler_lower_limiting_load_value: 3000,
	Startup_3_boiler_lower_limiting_load_value: 7000,
	Startup_4_boiler_lower_limiting_load_value: 12000,

	Week_load_prediction: loadWeekly,
}

//TODO:未来一周的温度
func GetHeatStorageWeek(c *gin.Context) {
	a := energyWeekly.GetHeatStorageAagin()
	b := []int{4, 4, 2, 5, 6, 7, 6}
	//x := []string{"2023-02-28", "2023-03-01", "2023-03-02", "2023-03-03", "2023-03-04", "2023-03-05", "2023-03-06"}
	x := MakeX()
	c.JSON(http.StatusOK, gin.H{
		"再蓄热量": a,
		"室外温度": b,
		"x轴":   x,
	})
}

func GetElectricityWeek(c *gin.Context) {
	a := energyWeekly.GetPeakTransferPower(energyWeekly.GetHeatStorageAagin())
	b := []int{4, 4, 2, 5, 6, 7, 6}
	//x := []string{"2023-02-28", "2023-03-01", "2023-03-02", "2023-03-03", "2023-03-04", "2023-03-05", "2023-03-06"}
	x := MakeX()
	c.JSON(http.StatusOK, gin.H{
		"移峰电量": a,
		"室外温度": b,
		"x轴":   x,
	})
}

func GetConfigWeek(c *gin.Context) {
	vally, other := energyWeekly.GetBoilerRunningTime()
	//x := []string{"2月28号", "3月1号", "3月2号", "3月3号", "3月4号", "3月5号", "3月6号"}
	x := MakeX()
	c.JSON(http.StatusOK, gin.H{
		"谷电价":  vally,
		"峰平电价": other,
		"x轴":   x,
	})
}

func GetEnergySaving(c *gin.Context) {
	a := energyWeekly.GetEnergySaving()
	//a := []float64{103, 127, 113, 145, 110, 87, 105}
	//x := []string{"2023-02-28", "2023-03-01", "2023-03-02", "2023-03-03", "2023-03-04", "2023-03-05", "2023-03-06"}
	x := MakeX()
	c.JSON(http.StatusOK, gin.H{
		"data": a,
		"x轴":   x,
	})
}

func GetRunningCost(c *gin.Context) {
	a := energyWeekly.GetRunningCost()
	//x := []string{"2023-02-28", "2023-03-01", "2023-03-02", "2023-03-03", "2023-03-04", "2023-03-05", "2023-03-06"}
	x := MakeX()
	c.JSON(http.StatusOK, gin.H{
		"data": a,
		"x轴":   x,
	})
}

func GetCarbonEmission(c *gin.Context) {
	a := energyWeekly.GetCarbonEmission(energyWeekly.GetEnergySaving())
	//a := []float64{54, 43, 47, 51, 61, 41, 52}
	//x := []string{"2023-02-28", "2023-03-01", "2023-03-02", "2023-03-03", "2023-03-04", "2023-03-05", "2023-03-06"}
	x := MakeX()
	c.JSON(http.StatusOK, gin.H{
		"data": a,
		"x轴":   x,
	})
}

func MakeX() []string {
	return []string{model.GetDay(time.Now().Unix()),
		model.GetDay(time.Now().Unix() + 86400),
		model.GetDay(time.Now().Unix() + 86400*2),
		model.GetDay(time.Now().Unix() + 86400*3),
		model.GetDay(time.Now().Unix() + 86400*4),
		model.GetDay(time.Now().Unix() + 86400*5),
		model.GetDay(time.Now().Unix() + 86400*6)}
}
