package routes

import (
	"energy/api/energyConfig"
	"github.com/gin-gonic/gin"
)

func EnergyConfigRouter(router *gin.RouterGroup) {
	router.GET("energyConfig/getPeriod", energyConfig.GetPeriod)                       //错峰用电
	router.GET("energyConfig/getDeviceWorkTime", energyConfig.GetDeviceWorkTime)       //今日设备运行工况   0
	router.GET("energyConfig/getLoadDetail", energyConfig.GetLoadDetail)               //负荷明细统计
	router.GET("energyConfig/getBoilerConfigDaily", energyConfig.GetBoilerConfigDaily) //电锅炉逐时建议工况
	router.GET("energyConfig/getTankConfigDaily", energyConfig.GetTankConfigDaily)     //蓄热水箱逐时建议工况
	router.GET("energyConfig/getDeviceWorkState", energyConfig.GetDeviceWorkState)     //设备运行状态

	router.GET("energyConfig/getHeatStorageWeek", energyConfig.GetHeatStorageWeek) //未来七天再蓄热量
	router.GET("energyConfig/getElectricityWeek", energyConfig.GetElectricityWeek) //未来七天移峰电量
	router.GET("energyConfig/getConfigWeek", energyConfig.GetConfigWeek)           //周工况调节

	router.GET("energyConfig/getEnergySaving", energyConfig.GetEnergySaving)     //节约能耗
	router.GET("energyConfig/getRunningCost", energyConfig.GetRunningCost)       //运行费用
	router.GET("energyConfig/getCarbonEmission", energyConfig.GetCarbonEmission) //减少碳排放
}
