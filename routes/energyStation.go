package routes

import (
	"energy/api/energy_analysis"

	"github.com/gin-gonic/gin"
)

func EnergyStationRouter(router *gin.RouterGroup) {
	router.GET("energyStation/boilerEfficiencyDay", energy_analysis.BoilerEfficiencyDay) //电锅炉效率
	router.GET("energyStation/carbonEmission", energy_analysis.CarbonEmission)           //碳排放量统计
	router.GET("energyStation/payLoad", energy_analysis.PayLoad)                         //负载率统计
}
