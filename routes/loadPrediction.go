package routes

import (
	"energy/api/loadPredict"
	"github.com/gin-gonic/gin"
)

func LoadPredictRouter(router *gin.RouterGroup) {
	//router.GET("loadPredict/loadStatistic", loadPredict.GetLoadStatistic) //负荷实时统计
	//router.GET("loadPredict/comparison", loadPredict.GetComparison)       //对比

	router.GET("loadPredict/test", loadPredict.GetRealLoad) //对比
}
