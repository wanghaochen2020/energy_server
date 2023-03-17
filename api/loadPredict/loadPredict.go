package loadPredict

import (
	"energy/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetLoadStatistic(c *gin.Context) {
	index := c.Query("index")
	//a := model.GetLoad(index, "today")
	//b := model.GetData("temperature", int(time.Now().Unix()),"")
	var a, b, x []int

	if index == "D1组团" {
		a = []int{111, 116, 115, 120, 121, 110, 110, 106, 96, 80, 65, 51, 41, 35, 25, 51, 66, 71, 0, 0, 0, 0, 0, 0, 0, 0}
		b = []int{-2, -3, -3, -4, -4, -2, -2, -1, 1, 4, 7, 10, 12, 13, 15, 10, 7, 6, 5, 4, 4, 3, 2, 1, 0, -1}
	}

	x = make([]int, 24)
	for i := 0; i < 24; i++ {
		x[i] = i
	}

	c.JSON(http.StatusOK, gin.H{
		"x轴":  x,
		"负荷量": a,
		"温度":  b,
	})
}

func GetComparison(c *gin.Context) {
	index := c.Query("index")
	var a, b, x []int
	//a := model.GetLoad(index, "today")
	//b := model.LoadPredict(index)

	if index == "D1组团" {
		a = []int{111, 116, 115, 120, 121, 110, 110, 106, 96, 80, 65, 51, 41, 35, 25, 51, 66, 71, 0, 0, 0, 0, 0, 0, 0, 0}
		b = []int{115, 123, 121, 123, 127, 118, 113, 109, 106, 65, 60, 58, 39, 42, 22, 51, 48, 73, 56, 99, 84, 95, 108, 106, 110, 113}
	}
	x = make([]int, 24)
	for i := 0; i < 24; i++ {
		x[i] = i
	}

	c.JSON(http.StatusOK, gin.H{
		"x轴":  x,
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
	a := model.GetData("temperature", int(time.Now().Unix()), "")

	c.JSON(http.StatusOK, gin.H{
		"code": a,
	})
}
