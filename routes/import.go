// 导入历史数据
package routes

import (
	"energy/api"

	"github.com/gin-gonic/gin"
)

func ImportDBRouter(router *gin.RouterGroup) {
	router.GET("update_history_data", api.UpdateHistoryData) //更新给定时间的历史数据
}
