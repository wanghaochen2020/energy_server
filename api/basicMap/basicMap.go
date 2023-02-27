package basicMap

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAtmosphere(c *gin.Context) {
	/*
		c.JSON(http.StatusOK, gin.H{
			"风速":   "2.4",
			"湿度":   "30",
			"温度":   "5",
			"总辐射":  "450",
			"大气压力": "920",
		})

	*/
	c.JSON(http.StatusOK, gin.H{
		"data": []string{"5", "30", "450", "2.4", "920"},
	})

}

func GetKekong(c *gin.Context) {

	/*
		c.JSON(http.StatusOK, gin.H{
			"D1": 80,
			"D2": 60,
			"D3": 90,
			"D4": 70,
			"D5": 80,
			"D6": 50,
		})

	*/
	c.JSON(http.StatusOK, gin.H{
		"data": []int{80, 60, 90, 70, 80, 50},
	})
}
