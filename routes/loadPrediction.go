package routes

import (
	"energy/api/loadPredict"
	"github.com/gin-gonic/gin"
)

func LoadPredictRouter(router *gin.RouterGroup) {
	router.GET("loadPredict/loadStatistic", loadPredict.GetLoadStatistic) //负荷实时统计
	router.GET("loadPredict/comparison", loadPredict.GetComparison)       //对比

	router.GET("loadPredict/realLoad", loadPredict.GetRealLoad)    //真实值
	router.GET("loadPredict/forecast", loadPredict.GetLoadPredict) //预测值
	router.GET("loadPredict/temp", loadPredict.GetTemp)            //温度
}
