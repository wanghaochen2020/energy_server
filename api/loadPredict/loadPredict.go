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
		a = []int{130, 150, 60, 230, 224, 100, 218, 135, 80, 147, 260, 200, 150, 60,
			230, 224, 100, 218, 135, 80, 147, 260, 200, 100}
		b = []int{55, 35, 20, 17, 16, 20, 30, 20, 30, 20, 30, 24, 23, 18,
			30, 27, 16, 10, 10, 13, 24, 10, 18, 35}
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
		a = []int{130, 150, 60, 230, 224, 100, 218, 135, 80, 147, 260, 200, 150, 60,
			230, 224, 100, 218, 135, 80, 147, 260, 200, 100}
		b = []int{150, 35, 80, 47, 160, 100, 50, 60, 50, 60, 30, 124, 60, 118,
			80, 47, 160, 100, 100, 130, 124, 100, 118, 35}
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
