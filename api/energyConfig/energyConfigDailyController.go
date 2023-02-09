package energyConfig

import (
	"energy/defs"
	_ "energy/defs"
	"energy/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EnergyConfigDailyController struct {
}

var (
	Vally_cost_time_start  = 23
	Vally_cost_time_end    = 7
	Flat_cost_time_1_start = 7
	Flat_cost_time_1_end   = 10
	Flat_cost_time_2_start = 15
	Flat_cost_time_2_end   = 18
	Flat_cost_time_3_start = 21
	Flat_cost_time_3_end   = 23
	Peak_cost_time_1_start = 10
	Peak_cost_time_1_end   = 15
	Peak_cost_time_2_start = 18
	Peak_cost_time_2_end   = 21
)

var loadDaily = [24]float64{537.41, 586.16, 618.91, 607.23, 608.5, 621.55, 645.52, 890.35, 690.17, 501.28, 1204.25, 915.07, 793.98, 748.76, 714.84, 694.95, 657.61, 681.41, 791.54, 999.22, 1156.91, 1264.27, 828.37, 661.38}
var energy = model.EnergyConfigDaily{
	Qs:                      1000,
	Tank_top_export_temp:    80,
	Tank_bottom_export_temp: 80,
	Vally_cost_time_start:   Vally_cost_time_start,
	Vally_cost_time_end:     Vally_cost_time_end,
	Flat_cost_time_1_start:  Flat_cost_time_1_start,
	Flat_cost_time_1_end:    Flat_cost_time_1_end,
	Flat_cost_time_2_start:  Flat_cost_time_2_start,
	Flat_cost_time_2_end:    Flat_cost_time_2_end,
	Flat_cost_time_3_start:  Flat_cost_time_3_start,
	Flat_cost_time_3_end:    Flat_cost_time_3_end,
	Peak_cost_time_1_start:  Peak_cost_time_1_start,
	Peak_cost_time_1_end:    Peak_cost_time_1_end,
	Peak_cost_time_2_start:  Peak_cost_time_2_start,
	Peak_cost_time_2_end:    Peak_cost_time_2_end,

	Startup_1_boiler_lower_limiting_load_value: 400,
	Startup_2_boiler_lower_limiting_load_value: 3000,
	Startup_3_boiler_lower_limiting_load_value: 7000,
	Startup_4_boiler_lower_limiting_load_value: 12000,

	Daily_load_prediction: loadDaily,
}

func GetPeriod(c *gin.Context) {
	flat := [6]int{Flat_cost_time_1_start, Flat_cost_time_1_end, Flat_cost_time_2_start, Flat_cost_time_2_end, Flat_cost_time_3_start, Flat_cost_time_3_end}
	peak := [4]int{Peak_cost_time_1_start, Peak_cost_time_1_end, Peak_cost_time_2_start, Peak_cost_time_2_end}
	vally := [2]int{Vally_cost_time_start, Vally_cost_time_end}
	c.JSON(http.StatusOK, gin.H{
		"平电价": flat,
		"峰电价": peak,
		"谷电价": vally,
	})
}

func GetDeviceWorkTime(c *gin.Context) {
	var result [9]int
	for i := 0; i < 9; i++ {
		//data, _ := model.GetResultFloatList(defs.EnergyRunningTimeDay[i], model.UnixToString(int(time.Now().Unix())))
		data, _ := model.GetResultFloatList(defs.EnergyRunningTimeDay[i], "2022/10/12")
		for j := 0; j < len(data); j++ {
			result[i] += int(data[j])
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func GetLoadDetail(c *gin.Context) {
	fullBoilerLoad := energy.GetBoilerLoad()
	tankHeating := energy.GetTankHeatingLoad()
	var tankHeating2 [24]int
	var tankStorage [24]int
	var boilerLoad [24]int

	len := 24 - len(tankHeating)

	for i := len; i < 24; i++ {
		tankHeating2[i] = int(tankHeating[i-len])
	}
	for i := 0; i < len; i++ {
		tankStorage[i] = int(fullBoilerLoad[i] - loadDaily[i])
		boilerLoad[i] = int(loadDaily[i])
	}
	for i := len; i < 24; i++ {
		boilerLoad[i] = int(loadDaily[i] - tankHeating[i-(len)])
	}

	c.JSON(http.StatusOK, gin.H{
		"电锅炉负荷":  boilerLoad,
		"水箱蓄热负荷": tankStorage,
		"水箱放热负荷": tankHeating2,
	})
}

func GetBoilerConfigDaily(c *gin.Context) {
	a, _ := model.GetResultFloatList(defs.EnergyBoilerRunningNum, "2022/10/12")
	c.JSON(http.StatusOK, gin.H{
		"实际": a,
		"建议": energy.GetBoilerRunningNum(),
	})
}

func GetData(c *gin.Context) {
	// fmt.Println("水箱放热：", energy.GetTankHeatingLoad())
	c.JSON(http.StatusOK, gin.H{
		"水箱放热":      energy.GetTankHeatingLoad(),
		"电锅炉承担逐时负荷": energy.GetBoilerLoad(),
	})
}

func GetConsumption(c *gin.Context) {
	time := c.Query("time")
	a, b := model.GetResultFloatList(defs.GroupHeatConsumptionHour4, time)
	fmt.Println(a)
	fmt.Println(b)
	c.JSON(http.StatusOK, gin.H{
		"code": a,
	})
}

func GetTankConfigDaily(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": energy.GetTankRecommendedHourlyWorkingCondition(),
	})
}

func GetWorkTime(c *gin.Context) {
	a, _ := model.GetResultFloatList(defs.EnergyRunningTimeDay[0], "2022/10/12")
	fmt.Println(a)
	c.JSON(http.StatusOK, gin.H{
		"code": a,
	})
}

func GetDeviceWorkState(c *gin.Context) {
	/*0就是关
	1-4 锅炉
	5-10 泵
	11-18 DV 1~8
	19-20 DVT
	*/
	//var array = [...]string{"ZLZ.系统运行中1", "ZLZ.系统运行中2", "ZLZ.系统运行中3", "ZLZ.系统运行中4", "ZLZ.RUN_P1", "ZLZ.RUN_P2", "ZLZ.RUN_P3", "ZLZ.RUN_P7", "ZLZ.RUN_P8", "ZLZ.RUN_P9", "ZLZ.OPEN_V1", "ZLZ.OPEN_V2", "ZLZ.OPEN_V3", "ZLZ.OPEN_V4", "ZLZ.OPEN_V5", "ZLZ.OPEN_V6", "ZLZ.OPEN_V8", "ZLZ.OPEN_V11", "ZLZ.OUTPUT_T29", "ZLZ.OUTPUT_T30"}
	var array = [...]string{"ZLZ.%E7%B3%BB%E7%BB%9F%E8%BF%90%E8%A1%8C%E4%B8%AD1", "ZLZ.%E7%B3%BB%E7%BB%9F%E8%BF%90%E8%A1%8C%E4%B8%AD2", "ZLZ.%E7%B3%BB%E7%BB%9F%E8%BF%90%E8%A1%8C%E4%B8%AD3", "ZLZ.%E7%B3%BB%E7%BB%9F%E8%BF%90%E8%A1%8C%E4%B8%AD4", "ZLZ.RUN_P1", "ZLZ.RUN_P2", "ZLZ.RUN_P3", "ZLZ.RUN_P7", "ZLZ.RUN_P8", "ZLZ.RUN_P9", "ZLZ.OPEN_V1", "ZLZ.OPEN_V2", "ZLZ.OPEN_V3", "ZLZ.OPEN_V4", "ZLZ.OPEN_V5", "ZLZ.OPEN_V6", "ZLZ.OPEN_V8", "ZLZ.OPEN_V11", "ZLZ.OUTPUT_T29", "ZLZ.OUTPUT_T30"}
	var array2 [20]int
	var result [22]int

	for i := 0; i < len(array); i++ {
		a, _ := model.GetOpcFloatList(array[i], "2022/10/12 13")
		if a[0] == 0 {
			array2[i] = 0
		} else {
			array2[i] = 1
		}
	}

	for i := 0; i < 4; i++ {
		result[i] = array2[i]
	}
	if array2[10] == 0 && array2[15] == 1 {
		result[4] = 1
		result[5] = 1
	}
	for i := 6; i < 21; i++ {
		result[i] = array2[i-2]
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
