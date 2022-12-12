// 导入历史数据
package api

import (
	"energy/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateHistoryData(c *gin.Context) {
	t0str := c.Query("tstart")
	t1str := c.Query("tend")
	t0, err := time.Parse("2006/01/02 15:04", t0str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "时间格式错误: " + err.Error(),
		})
	}
	t1, err := time.Parse("2006/01/02 15:04", t1str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "时间格式错误: " + err.Error(),
		})
	}
	go model.UpdateData(t0, t1)
	c.JSON(http.StatusAccepted, gin.H{
		"code": http.StatusAccepted,
		"msg":  "正在更新中",
	})
}
