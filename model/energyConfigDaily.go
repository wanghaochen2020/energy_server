package model

import (
	"fmt"
	"strconv"
)

//import "fmt"

type EnergyConfigDaily struct {
	Qs                      float64 //水箱满蓄蓄热量
	Tank_top_export_temp    float64 //水箱顶部出口侧温度
	Tank_bottom_export_temp float64 //水箱底部出口侧温度

	Vally_cost_time_start  int //谷电价阶段起始时间
	Vally_cost_time_end    int //谷电价阶段结束时间
	Flat_cost_time_1_start int
	Flat_cost_time_1_end   int
	Flat_cost_time_2_start int
	Flat_cost_time_2_end   int
	Flat_cost_time_3_start int
	Flat_cost_time_3_end   int
	Peak_cost_time_1_start int
	Peak_cost_time_1_end   int
	Peak_cost_time_2_start int
	Peak_cost_time_2_end   int

	Startup_1_boiler_lower_limiting_load_value float64 //启动一台电锅炉下限负荷值，单位：kw
	Startup_2_boiler_lower_limiting_load_value float64
	Startup_3_boiler_lower_limiting_load_value float64
	Startup_4_boiler_lower_limiting_load_value float64

	Daily_load_prediction [24]float64 //日负荷值(预测值)，23点给出
}

func (energy EnergyConfigDaily) SortArray2DMAXtoMIN(array [][]float64) (arrayOutput [][]float64) {
	var array_backup [][]float64
	for i := 0; i < len(array); i++ {
		array_backup = append(array_backup, array[i])
	}
	for i := 0; i < len(array); i++ {
		for j := len(array) - 1; j > i; j-- {
			if array_backup[j][1] > array_backup[j-1][1] {
				array_backup[j], array_backup[j-1] = array_backup[j-1], array_backup[j]
			}
		}
	}
	return array_backup
}

/*
 *功能：削峰蓄热计算
 *输入：排序后的二维数组[时间序号][负荷]
 *输出：一维数组（变化后的逐时负荷值，没变化赋0）
 */
func (energy EnergyConfigDaily) ClipPeakStorageCalculation(loadArray [][]float64) []float64 {
	var loadArray_backup [][]float64
	for i := 0; i < len(loadArray); i++ {
		loadArray_backup = append(loadArray_backup, loadArray[i])
	}
	loadArray_backup = energy.SortArray2DMAXtoMIN(loadArray_backup)
	var result []float64
	for i := 0; i < len(loadArray_backup); i++ {
		result = append(result, 0)
	}
	loadAccumulateValue := loadArray_backup[len(loadArray_backup)-1][1]
	for i := len(loadArray_backup) - 1; i > 0; i-- {
		if (loadAccumulateValue+energy.Qs)/(float64(len(loadArray_backup)-i)) <= loadArray_backup[i-1][1] {
			for j := len(loadArray_backup) - 1; j >= i; j-- {
				result[int(loadArray_backup[j][0])] = (loadAccumulateValue + energy.Qs) / (float64(len(loadArray_backup) - i))
			}
			break
		} else {
			loadAccumulateValue = loadAccumulateValue + loadArray_backup[i-1][1]
		}

		if i == 1 {
			for j := len(loadArray_backup) - 1; j >= 0; j-- {
				result[int(loadArray_backup[j][0])] = (loadAccumulateValue + energy.Qs) / (float64(len(loadArray_backup)))
			}
		}
	}

	return result
}

/*
 *功能：削峰放热计算
 *输入：二维数组[时间序号][负荷] 水箱放热量
 *输出：一维数组（变化后的逐时负荷值，没变化赋0）
 */
func (energy EnergyConfigDaily) ClipPeakHeatingCalculation(loadArray [][]float64, Qs float64) []float64 {
	var loadArray_backup [][]float64
	for i := 0; i < len(loadArray); i++ {
		loadArray_backup = append(loadArray_backup, loadArray[i])
	}
	loadArray_backup = energy.SortArray2DMAXtoMIN(loadArray_backup)
	var result []float64
	for i := 0; i < len(loadArray_backup); i++ {
		result = append(result, 0)
	}
	loadAccumulateValue := loadArray_backup[0][1]

	for i := 0; i < len(loadArray_backup); i++ {
		if i == len(loadArray_backup)-1 {
			if Qs >= loadAccumulateValue {
				break
			} else {
				for j := 0; j < len(loadArray_backup); j++ {
					result[int(loadArray_backup[j][0])] = (loadAccumulateValue - Qs) / float64(len(loadArray_backup))
				}
				break
			}
		}
		if (loadAccumulateValue-Qs)/float64(i+1) >= loadArray_backup[i+1][1] {
			for j := 0; j <= i; j++ {
				result[int(loadArray_backup[j][0])] = (loadAccumulateValue - Qs) / float64(i+1)
			}
			for k := i + 1; k < len(loadArray_backup); k++ {
				result[int(loadArray_backup[k][0])] = loadArray_backup[k][1]
			}
			break
		} else {
			loadAccumulateValue = loadAccumulateValue + loadArray_backup[i+1][1]
		}
	}
	return result
}

/*
 *功能：获取水箱逐时蓄热量
 *输入：无
 *输出：水箱逐时蓄热量（一维数组）
 */
func (energy EnergyConfigDaily) GetTankStorageLoad() []float64 {
	vally_time_hour_num := energy.Vally_cost_time_end + (24 - energy.Vally_cost_time_start)
	var tankStorage []float64
	var tankStorageLoad [][]float64

	for i := 0; i < vally_time_hour_num; i++ {
		var tankStorageLoad1D = []float64{float64(i), energy.Daily_load_prediction[i]}
		tankStorageLoad = append(tankStorageLoad, tankStorageLoad1D)
	}

	tankStorage = energy.ClipPeakStorageCalculation(tankStorageLoad)

	for i := 0; i < len(tankStorage); i++ {
		if tankStorage[i] != 0 {
			tankStorage[i] = tankStorage[i] - energy.Daily_load_prediction[i]
		}
	}
	return tankStorage
}

/*
 *功能：获取水箱逐时放热量
 *输入：无
 *输出：水箱逐时放热量（一维数组）
 */
func (energy EnergyConfigDaily) GetTankHeatingLoad() []float64 {
	Qs1 := energy.Qs
	peak_time_1_hour_num := energy.Peak_cost_time_1_end - energy.Peak_cost_time_1_start
	peak_time_2_hour_num := energy.Peak_cost_time_2_end - energy.Peak_cost_time_2_start
	flat_time_1_hour_num := energy.Flat_cost_time_1_end - energy.Flat_cost_time_1_start
	flat_time_2_hour_num := energy.Flat_cost_time_2_end - energy.Flat_cost_time_2_start
	flat_time_3_hour_num := energy.Flat_cost_time_3_end - energy.Flat_cost_time_3_start
	tank_heating_hour_num := energy.Flat_cost_time_3_end - energy.Flat_cost_time_1_start

	peak_time_1_load_sum, peak_time_2_load_sum, flat_time_1_load_sum, flat_time_2_load_sum, flat_time_3_load_sum := 0.0, 0.0, 0.0, 0.0, 0.0

	for i := energy.Peak_cost_time_1_start + 1; i <= energy.Peak_cost_time_1_start+peak_time_1_hour_num; i++ {
		peak_time_1_load_sum = peak_time_1_load_sum + energy.Daily_load_prediction[i]
	}
	for i := energy.Peak_cost_time_2_start + 1; i <= energy.Peak_cost_time_2_start+peak_time_2_hour_num; i++ {
		peak_time_2_load_sum = peak_time_2_load_sum + energy.Daily_load_prediction[i]
	}
	for i := energy.Flat_cost_time_1_start + 1; i <= energy.Flat_cost_time_1_start+flat_time_1_hour_num; i++ {
		flat_time_1_load_sum = flat_time_1_load_sum + energy.Daily_load_prediction[i]
	}
	for i := energy.Flat_cost_time_2_start + 1; i <= energy.Flat_cost_time_2_start+flat_time_2_hour_num; i++ {
		flat_time_2_load_sum = flat_time_2_load_sum + energy.Daily_load_prediction[i]
	}
	for i := energy.Flat_cost_time_3_start + 1; i <= energy.Flat_cost_time_3_start+flat_time_3_hour_num; i++ {
		flat_time_3_load_sum = flat_time_3_load_sum + energy.Daily_load_prediction[i]
	}

	var tankHeating []float64
	for i := 0; i < tank_heating_hour_num; i++ {
		tankHeating = append(tankHeating, 0)
	}

	if Qs1 >= peak_time_1_load_sum+peak_time_2_load_sum {
		for i := energy.Peak_cost_time_1_start + 1; i <= energy.Peak_cost_time_1_start+peak_time_1_hour_num; i++ {
			tankHeating[i-(24-tank_heating_hour_num)] = energy.Daily_load_prediction[i]
		}
		for i := energy.Peak_cost_time_2_start + 1; i <= energy.Peak_cost_time_2_start+peak_time_2_hour_num; i++ {
			tankHeating[i-(24-tank_heating_hour_num)] = energy.Daily_load_prediction[i]
		}
		Qs1 = Qs1 - peak_time_1_load_sum - peak_time_2_load_sum
	} else {
		var tankHeatingPeak []float64
		var tankHeatingLoadPeak [][]float64
		for i := 0; i < peak_time_1_hour_num; i++ {
			tankHeatingLoadPeak1D := []float64{float64(i), energy.Daily_load_prediction[energy.Peak_cost_time_1_start+1+i]}
			tankHeatingLoadPeak = append(tankHeatingLoadPeak, tankHeatingLoadPeak1D)
		}
		for i := peak_time_1_hour_num; i < peak_time_1_hour_num+peak_time_2_hour_num; i++ {
			tankHeatingLoadPeak1D := []float64{float64(i), energy.Daily_load_prediction[energy.Peak_cost_time_2_start+1+i-peak_time_1_hour_num]}
			tankHeatingLoadPeak = append(tankHeatingLoadPeak, tankHeatingLoadPeak1D)
		}

		tankHeatingPeak = energy.ClipPeakHeatingCalculation(tankHeatingLoadPeak, Qs1)
		//获取水箱放热量
		for i := energy.Peak_cost_time_1_start + 1; i <= energy.Peak_cost_time_1_start+peak_time_1_hour_num; i++ {
			tankHeating[i-(24-tank_heating_hour_num)] = energy.Daily_load_prediction[i] - tankHeatingPeak[i-energy.Peak_cost_time_1_start-1]
		}
		for i := energy.Peak_cost_time_2_start + 1; i <= energy.Peak_cost_time_2_start+peak_time_2_hour_num; i++ {
			tankHeating[i-(24-tank_heating_hour_num)] = energy.Daily_load_prediction[i] - tankHeatingPeak[i+peak_time_1_hour_num-energy.Peak_cost_time_2_start-1]
		}

		return tankHeating

	}
	//平三阶段水箱放热量计算
	if Qs1 >= flat_time_3_load_sum {
		for i := energy.Flat_cost_time_3_start + 1; i <= energy.Flat_cost_time_3_start+flat_time_3_hour_num; i++ {
			tankHeating[i-(24-tank_heating_hour_num)] = energy.Daily_load_prediction[i]
		}
		Qs1 = Qs1 - flat_time_3_load_sum
	} else {
		var tankHeatingF3 []float64
		var tankHeatingLoadF3 [][]float64
		//传入平三阶段数据
		for i := 0; i < flat_time_3_hour_num; i++ {
			tankHeatingLoadF31D := []float64{float64(i), energy.Daily_load_prediction[energy.Flat_cost_time_3_start+1+i]}
			tankHeatingLoadF3 = append(tankHeatingLoadF3, tankHeatingLoadF31D)
		}
		tankHeatingF3 = energy.ClipPeakHeatingCalculation(tankHeatingLoadF3, Qs1)
		//获取水箱放热量
		for i := energy.Flat_cost_time_3_start + 1; i <= energy.Flat_cost_time_3_start+flat_time_3_hour_num; i++ {
			tankHeating[i-(24-tank_heating_hour_num)] = energy.Daily_load_prediction[i] - tankHeatingF3[i-energy.Flat_cost_time_3_start-1]
		}
		return tankHeating

	}
	//平二阶段水箱放热量计算
	if Qs1 >= flat_time_2_load_sum {
		for i := energy.Flat_cost_time_2_start + 1; i <= energy.Flat_cost_time_2_start+flat_time_2_hour_num; i++ {
			tankHeating[i-(24-tank_heating_hour_num)] = energy.Daily_load_prediction[i]
		}

		Qs1 = Qs1 - flat_time_2_load_sum
	} else {
		var tankHeatingF2 []float64
		var tankHeatingLoadF2 [][]float64
		//传入平二阶段数据
		for i := 0; i < flat_time_2_hour_num; i++ {
			tankHeatingLoadF21D := []float64{float64(i), energy.Daily_load_prediction[energy.Flat_cost_time_2_start+1+i]}
			tankHeatingLoadF2 = append(tankHeatingLoadF2, tankHeatingLoadF21D)
		}
		tankHeatingF2 = energy.ClipPeakHeatingCalculation(tankHeatingLoadF2, Qs1)
		//获取水箱放热量
		for i := energy.Flat_cost_time_2_start + 1; i <= energy.Flat_cost_time_2_start+flat_time_2_hour_num; i++ {
			tankHeating[i-(24-tank_heating_hour_num)] = energy.Daily_load_prediction[i] - tankHeatingF2[i-energy.Flat_cost_time_2_start-1]
		}
		return tankHeating
	}
	//平一阶段水箱放热量计算
	if Qs1 >= flat_time_1_load_sum {
		for i := energy.Flat_cost_time_1_start + 1; i <= energy.Flat_cost_time_1_start+flat_time_1_hour_num; i++ {
			tankHeating[i-(24-tank_heating_hour_num)] = energy.Daily_load_prediction[i]
		}

		Qs1 = Qs1 - flat_time_1_load_sum
	} else {
		var tankHeatingF1 []float64
		var tankHeatingLoadF1 [][]float64
		//传入平二阶段数据
		for i := 0; i < flat_time_1_hour_num; i++ {
			tankHeatingLoadF11D := []float64{float64(i), energy.Daily_load_prediction[energy.Flat_cost_time_1_start+1+i]}
			tankHeatingLoadF1 = append(tankHeatingLoadF1, tankHeatingLoadF11D)
		}
		tankHeatingF1 = energy.ClipPeakHeatingCalculation(tankHeatingLoadF1, Qs1)
		//获取水箱放热量
		for i := energy.Flat_cost_time_1_start + 1; i <= energy.Flat_cost_time_1_start+flat_time_1_hour_num; i++ {
			tankHeating[i-(24-tank_heating_hour_num)] = energy.Daily_load_prediction[i] - tankHeatingF1[i-energy.Flat_cost_time_1_start-1]
		}
		return tankHeating
	}
	return tankHeating
}

/*
 *功能：获取电锅炉承担负荷
 *输入：无
 *输出：电锅炉逐时负荷
 */
func (energy EnergyConfigDaily) GetBoilerLoad() [24]float64 {
	var boilerLoad [24]float64
	for i := 0; i < len(energy.GetTankStorageLoad()); i++ {
		boilerLoad[i] = energy.GetTankStorageLoad()[i] + energy.Daily_load_prediction[i]
	}
	for i := len(energy.GetTankStorageLoad()); i < len(energy.GetTankStorageLoad())+len(energy.GetTankHeatingLoad()); i++ {
		boilerLoad[i] = energy.Daily_load_prediction[i] - energy.GetTankHeatingLoad()[i-len(energy.GetTankStorageLoad())]
	}
	return boilerLoad
}

/*
 *功能：获取电锅炉逐时运行台数
 *输入：无
 *输出：实时电锅炉运行台数
 */
func (energy EnergyConfigDaily) GetBoilerRunningNum() [24]int {
	var boilerRunningNum [24]int
	var boilerLoad [24]float64
	for i := 0; i < 24; i++ {
		boilerLoad[i] = energy.GetBoilerLoad()[i]
	}
	for i := 0; i < 24; i++ {
		if boilerLoad[i] >= energy.Startup_1_boiler_lower_limiting_load_value && boilerLoad[i] < energy.Startup_2_boiler_lower_limiting_load_value {
			boilerRunningNum[i] = 1
		} else if boilerLoad[i] >= energy.Startup_2_boiler_lower_limiting_load_value && boilerLoad[i] < energy.Startup_3_boiler_lower_limiting_load_value {
			boilerRunningNum[i] = 2
		} else if boilerLoad[i] >= energy.Startup_3_boiler_lower_limiting_load_value && boilerLoad[i] < energy.Startup_4_boiler_lower_limiting_load_value {
			boilerRunningNum[i] = 3
		} else if boilerLoad[i] >= energy.Startup_4_boiler_lower_limiting_load_value {
			boilerRunningNum[i] = 4
		} else if boilerLoad[i] < energy.Startup_1_boiler_lower_limiting_load_value {
			boilerRunningNum[i] = 0
		}
	}
	return boilerRunningNum
}

/*
 *功能：获取水箱逐时建议工况
 *输入：无
 *输出：[C B1 A B2 A B3]
 */
func (energy EnergyConfigDaily) GetTankRecommendedHourlyWorkingCondition() [6]float64 {
	peak_time_1_hour_num := energy.Peak_cost_time_1_end - energy.Peak_cost_time_1_start
	peak_time_2_hour_num := energy.Peak_cost_time_2_end - energy.Peak_cost_time_2_start
	flat_time_1_hour_num := energy.Flat_cost_time_1_end - energy.Flat_cost_time_1_start
	flat_time_2_hour_num := energy.Flat_cost_time_2_end - energy.Flat_cost_time_2_start
	flat_time_3_hour_num := energy.Flat_cost_time_3_end - energy.Flat_cost_time_3_start

	peak_time_1_load_sum, peak_time_2_load_sum, flat_time_1_load_sum, flat_time_2_load_sum, flat_time_3_load_sum := 0.0, 0.0, 0.0, 0.0, 0.0

	for i := energy.Peak_cost_time_1_start + 1; i <= energy.Peak_cost_time_1_start+peak_time_1_hour_num; i++ {
		peak_time_1_load_sum = peak_time_1_load_sum + energy.Daily_load_prediction[i]
	}
	for i := energy.Peak_cost_time_2_start + 1; i <= energy.Peak_cost_time_2_start+peak_time_2_hour_num; i++ {
		peak_time_2_load_sum = peak_time_2_load_sum + energy.Daily_load_prediction[i]
	}
	for i := energy.Flat_cost_time_1_start + 1; i <= energy.Flat_cost_time_1_start+flat_time_1_hour_num; i++ {
		flat_time_1_load_sum = flat_time_1_load_sum + energy.Daily_load_prediction[i]
	}
	for i := energy.Flat_cost_time_2_start + 1; i <= energy.Flat_cost_time_2_start+flat_time_2_hour_num; i++ {
		flat_time_2_load_sum = flat_time_2_load_sum + energy.Daily_load_prediction[i]
	}
	for i := energy.Flat_cost_time_3_start + 1; i <= energy.Flat_cost_time_3_start+flat_time_3_hour_num; i++ {
		flat_time_3_load_sum = flat_time_3_load_sum + energy.Daily_load_prediction[i]
	}

	var Q_z1, Q_z2 float64
	if (energy.Qs - peak_time_1_load_sum - peak_time_2_load_sum - flat_time_3_load_sum) > 0 {
		Q_z2 = energy.Qs - peak_time_1_load_sum - peak_time_2_load_sum - flat_time_3_load_sum
	} else {
		Q_z2 = 0.0
		Q_z1 = 0.0
	}
	if energy.Qs-peak_time_1_load_sum-peak_time_2_load_sum-flat_time_3_load_sum-flat_time_2_load_sum > 0 {
		Q_z1 = energy.Qs - peak_time_1_load_sum - peak_time_2_load_sum - flat_time_3_load_sum - flat_time_2_load_sum
	} else {
		Q_z1 = 0.0
	}

	var tankRecommendedHourlyWorkingCondition [6]float64
	tankRecommendedHourlyWorkingCondition[0] = energy.Tank_bottom_export_temp
	tankRecommendedHourlyWorkingCondition[2] = energy.Tank_top_export_temp
	tankRecommendedHourlyWorkingCondition[4] = energy.Tank_top_export_temp
	tankRecommendedHourlyWorkingCondition[5] = energy.Tank_top_export_temp
	tankRecommendedHourlyWorkingCondition[1], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", Q_z1), 64)
	//tankRecommendedHourlyWorkingCondition[3] = Q_z2
	tankRecommendedHourlyWorkingCondition[3], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", Q_z2), 64)

	return tankRecommendedHourlyWorkingCondition
}
