package routes

import (
	"energy/api/energy_analysis"
	"energy/api/system_layer"

	"github.com/gin-gonic/gin"
)

func EnergyStationRouter(router *gin.RouterGroup) {
	//系统层
	router.GET("energyStation/onlineRate", system_layer.OnlineRate) //电锅炉效率
	//能效分析
	router.GET("energyStation/boilerEfficiencyDay", energy_analysis.BoilerEfficiencyDay) //电锅炉效率
	router.GET("energyStation/carbonEmission", energy_analysis.CarbonEmission)           //碳排放量统计
	router.GET("energyStation/payLoad", energy_analysis.PayLoad)                         //负载率统计
}
