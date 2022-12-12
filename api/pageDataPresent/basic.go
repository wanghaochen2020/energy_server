package pageDataPresent

import (
	"encoding/json"
	"energy/defs"
	"energy/model"
	"io/ioutil"
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

func OpcDataList(c *gin.Context) {
	name := c.Query("name")
	time := c.Query("time")
	d, _ := model.GetOpcFloatList(name, time)
	c.JSON(http.StatusOK, d)
}

func BasicDataSet(c *gin.Context) {
	r := c.Request
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "bad_request",
		})
		return
	}
	var ubody defs.BasicDataSetRequest
	if err = json.Unmarshal(res, &ubody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "bad_request",
		})
		return
	}
	ans := make(map[string]interface{})
	for _, v := range ubody.Data.BasicData {
		d, _ := model.GetResultFloatNoTime(v)
		ans[v] = d
	}
	for _, v := range ubody.Data.BasicDataListDay {
		d, _ := model.GetResultFloatList(v, ubody.DayStr)
		ans[v] = d
	}
	for _, v := range ubody.Data.BasicDataListHour {
		d, _ := model.GetResultFloatList(v, ubody.HourStr)
		ans[v] = d
	}
	for _, v := range ubody.Data.BasicOpcList {
		d, _ := model.GetOpcFloatList(v, ubody.HourStr)
		ans[v] = d
	}
	c.JSON(http.StatusOK, ans)
}
