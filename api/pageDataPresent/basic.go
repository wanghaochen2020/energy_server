package pageDataPresent

import (
	"energy/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicData(c *gin.Context) {
	name := c.Query("name")
	d, _ := model.GetResultFloatNoTime(name)
	c.JSON(http.StatusOK, d)
}

func BasicDataList(c *gin.Context) {
	name := c.Query("name")
	time := c.Query("time")
	d, _ := model.GetResultFloatList(name, time)
	c.JSON(http.StatusOK, d)
}
