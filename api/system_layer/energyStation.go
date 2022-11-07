package system_layer

import (
	"energy/model"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func OnlineRate(c *gin.Context) {
	// now := time.Now().Local()
	now, _ := time.Parse("2006/01/02 15:04:05", "2022/10/13 15:31:00")
	timeStr := fmt.Sprintf("%d/%02d/%02d %02d", now.Year(), now.Month(), now.Day(), now.Hour())
	d, _ := model.GetResultFloat("energystation_online_rate_hour", timeStr)
	c.JSON(200, d)
}
