package routes

import (
	"energy/api/basicMap"
	"github.com/gin-gonic/gin"
)

func BasicMapRouter(router *gin.RouterGroup) {
	router.GET("basicMap/getAtmosphere", basicMap.GetAtmosphere) //错峰用电
	router.GET("basicMap/getKekong", basicMap.GetKekong)         //今日设备运行工况   0
}
