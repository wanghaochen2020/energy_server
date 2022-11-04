package energy_analysis

import (
	"energy/model"

	"github.com/gin-gonic/gin"
)

func BoilerEfficiencyDay(c *gin.Context) {
	d := make([]map[string]interface{}, 3)
	d[0] = make(map[string]interface{})
	d[0]["name"] = "电锅炉"
	d[0]["data"] = model.GetOpcDataList("boiler_efficiency_day", 0)

	d[1] = make(map[string]interface{})
	d[1]["name"] = "蓄热水箱"
	d[1]["data"] = model.GetOpcDataList("watertank_efficiency_day", 0)

	d[2] = make(map[string]interface{})
	d[2]["name"] = "能源站系统"
	d[2]["data"] = model.GetOpcDataList("energystation_efficiency_day", 0)

	c.JSON(200, d)
}

func CarbonEmission(c *gin.Context) {
	d := make([]map[string]interface{}, 3)

	d[0] = make(map[string]interface{})
	d[0]["data"] = model.GetOpcDataList("energystation_carbon_day", 0)

	d[1] = make(map[string]interface{})
	d[1]["data"] = model.GetOpcDataList("energystation_carbon_week", 2)

	d[2] = make(map[string]interface{})
	d[2]["data"] = []float64{0}

	c.JSON(200, d)
}

func PayLoad(c *gin.Context) {
	d := make([]map[string]interface{}, 4)

	d[0] = make(map[string]interface{})
	d[0]["data"] = model.GetOpcDataList("energy_pay_load", 0)

	c.JSON(200, d)
}
