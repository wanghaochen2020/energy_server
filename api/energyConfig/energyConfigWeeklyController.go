package energyConfig

import (
	"energy/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ()

//var loadDaily = [24]float64{537.41, 586.16, 618.91, 607.23, 608.5, 621.55, 645.52, 890.35, 690.17, 501.28, 1204.25, 915.07, 793.98, 748.76, 714.84, 694.95, 657.61, 681.41, 791.54, 999.22, 1156.91, 1264.27, 828.37, 661.38}
var loadWeekly = [7][24]float64{loadDaily, loadDaily, loadDaily, loadDaily, loadDaily, loadDaily, loadDaily}
var energyWeekly = model.EnergyConfigWeekly{
	Qs: 1000,

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

	c.JSON(http.StatusOK, gin.H{
		"再蓄热量": a,
	})
}

func GetElectricityWeek(c *gin.Context) {
	a := energyWeekly.GetPeakTransferPower(energyWeekly.GetHeatStorageAagin())

	c.JSON(http.StatusOK, gin.H{
		"移峰电量": a,
	})
}

func GetConfigWeek(c *gin.Context) {
	vally, other := energyWeekly.GetBoilerRunningTime()

	c.JSON(http.StatusOK, gin.H{
		"谷电价":  vally,
		"峰平电价": other,
	})
}

func GetEnergySaving(c *gin.Context) {
	a := energyWeekly.GetEnergySaving()

	c.JSON(http.StatusOK, gin.H{
		"data": a,
	})
}

func GetRunningCost(c *gin.Context) {
	a := energyWeekly.GetRunningCost()

	c.JSON(http.StatusOK, gin.H{
		"data": a,
	})
}

func GetCarbonEmission(c *gin.Context) {
	a := energyWeekly.GetCarbonEmission(energyWeekly.GetEnergySaving())

	c.JSON(http.StatusOK, gin.H{
		"data": a,
	})
}
