package loadPredict

import (
	"energy/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetLoadStatistic(c *gin.Context) {
	index := c.Query("index")
	a := model.GetLoad(index, "today")
	b := model.GetData("temperature", int(time.Now().Unix()))

	c.JSON(http.StatusOK, gin.H{
		"负荷量": a,
		"温度":  b,
	})
}

func GetComparison(c *gin.Context) {
	index := c.Query("index")
	a := model.GetLoad(index, "today")
	b := model.LoadPredict(index)

	c.JSON(http.StatusOK, gin.H{
		"真实值": a,
		"预测值": b,
	})
}

func GetRealLoad(c *gin.Context) {
	index := c.Query("index")
	a := model.GetLoad(index, "today")

	c.JSON(http.StatusOK, gin.H{
		"真实值": a,
	})
}

func GetLoadPredict(c *gin.Context) {
	index := c.Query("index")
	a := model.LoadPredict(index)

	c.JSON(http.StatusOK, gin.H{
		"预测值": a,
	})
}

func GetTemp(c *gin.Context) {
	a := model.GetData("temperature", int(time.Now().Unix()))

	c.JSON(http.StatusOK, gin.H{
		"code": a,
	})
}
