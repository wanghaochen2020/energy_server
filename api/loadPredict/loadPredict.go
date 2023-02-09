package loadPredict

import (
	"energy/defs"
	"energy/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
func GetLoadStatistic(c *gin.Context) {
	index := c.Query("index")

	c.JSON(http.StatusOK, gin.H{
		"code": a,
	})
}

func GetComparison(c *gin.Context) {
	index := c.Query("index")

	c.JSON(http.StatusOK, gin.H{
		"code": a,
	})
}


*/
func GetRealLoad(c *gin.Context) {
	index := c.Query("index")
	var a []float64

	switch index {
	case "1":
		a, _ = model.GetResultFloatList(defs.GroupHeatConsumptionHour1, "2022/10/13 13")
	}

	fmt.Println(a)
	//fmt.Println(b)
	c.JSON(http.StatusOK, gin.H{
		"code": a,
	})
}
