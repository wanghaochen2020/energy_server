package model

//import "energyConfig/utils"

type EnergyConfigWeekly struct {
	Qs float64 //水箱满蓄蓄热量

	Heat_loss_rectify_coefficiency       float64 //热损失修正系数
	Heat_to_power_transform_coefficiency float64 //热电转化系数
	Heat_loss_coefficiency               float64 //节约能耗损失系数
	Carbon_emission_unit_power           float64 //单位碳排放量

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

	Vally_cost float64 //谷电价  单位：元/千瓦时
	Peak_cost  float64 //峰电价  单位：元/千瓦时
	Flat_cost  float64 //平电价  单位：元/千瓦时

	Startup_1_boiler_lower_limiting_load_value float64 //启动一台电锅炉下限负荷值，单位：kw
	Startup_2_boiler_lower_limiting_load_value float64
	Startup_3_boiler_lower_limiting_load_value float64
	Startup_4_boiler_lower_limiting_load_value float64

	Week_load_prediction [7][24]float64
}

/*
 *功能：再蓄热量
 *输入：无
 *输出：一维数组（逐天再蓄热量）
 */
func (energyWeekly EnergyConfigWeekly) GetHeatStorageAagin() [7]float64 {
	load_flat_and_peak_sum := 0.0
	var heatStorageAgain [7]float64
	for i := 0; i < 7; i++ {
		for j := energyWeekly.Flat_cost_time_1_start + 1; j <= energyWeekly.Flat_cost_time_3_end; j++ {
			load_flat_and_peak_sum = load_flat_and_peak_sum + energyWeekly.Week_load_prediction[i][j]
		}
		if (energyWeekly.Qs - load_flat_and_peak_sum) > 0 {
			heatStorageAgain[i] = energyWeekly.Qs - (1-energyWeekly.Heat_loss_rectify_coefficiency)*(energyWeekly.Qs-load_flat_and_peak_sum)
		} else {
			heatStorageAgain[i] = energyWeekly.Qs
		}

		load_flat_and_peak_sum = 0.0
	}
	return heatStorageAgain
}

/*
 *功能：移峰电量
 *输入：再蓄热量
 *输出：一维数组（逐天移峰电量）
 */
func (energyWeekly EnergyConfigWeekly) GetPeakTransferPower(heatStorageAgain [7]float64) [7]float64 {
	var peakTransferPower [7]float64
	for i := 0; i < 7; i++ {
		peakTransferPower[i] = energyWeekly.Heat_to_power_transform_coefficiency * heatStorageAgain[i]
	}
	return peakTransferPower
}

/*
 *功能：周工况调节
 *输入：无
 *输出：逐天谷电价阶段和非谷电价阶段各台数运行时间
 */
func (energyWeekly EnergyConfigWeekly) GetBoilerRunningTime() (C_stage [7][5]int, AB_stage [7][5]int) {
	var energyDaily EnergyConfigDaily
	energyDaily = EnergyConfigDaily{
		Qs:                     energyWeekly.Qs,
		Vally_cost_time_start:  energyWeekly.Vally_cost_time_start,
		Vally_cost_time_end:    energyWeekly.Vally_cost_time_end,
		Flat_cost_time_1_start: energyWeekly.Flat_cost_time_1_start,
		Flat_cost_time_1_end:   energyWeekly.Flat_cost_time_1_end,
		Flat_cost_time_2_start: energyWeekly.Flat_cost_time_2_start,
		Flat_cost_time_2_end:   energyWeekly.Flat_cost_time_2_end,
		Flat_cost_time_3_start: energyWeekly.Flat_cost_time_3_start,
		Flat_cost_time_3_end:   energyWeekly.Flat_cost_time_3_end,
		Peak_cost_time_1_start: energyWeekly.Peak_cost_time_1_start,
		Peak_cost_time_1_end:   energyWeekly.Peak_cost_time_1_end,
		Peak_cost_time_2_start: energyWeekly.Peak_cost_time_2_start,
		Peak_cost_time_2_end:   energyWeekly.Peak_cost_time_2_end,

		Startup_1_boiler_lower_limiting_load_value: energyWeekly.Startup_1_boiler_lower_limiting_load_value,
		Startup_2_boiler_lower_limiting_load_value: energyWeekly.Startup_2_boiler_lower_limiting_load_value,
		Startup_3_boiler_lower_limiting_load_value: energyWeekly.Startup_3_boiler_lower_limiting_load_value,
		Startup_4_boiler_lower_limiting_load_value: energyWeekly.Startup_4_boiler_lower_limiting_load_value,
	}
	var C_stage_1D, AB_stage_1D [5]int
	zeroBoilerNum, oneBoilerNum, twoBoilersNum, threeBoliersNum, fourBoilersNum := 0, 0, 0, 0, 0
	for i := 0; i < 7; i++ {
		energyDaily.Daily_load_prediction = energyWeekly.Week_load_prediction[i]
		var boilerRunningNum = energyDaily.GetBoilerRunningNum()
		for j := 0; j <= energyDaily.Vally_cost_time_end; j++ {
			if boilerRunningNum[j] == 0 {
				zeroBoilerNum++
			} else if boilerRunningNum[j] == 1 {
				oneBoilerNum++
			} else if boilerRunningNum[j] == 2 {
				twoBoilersNum++
			} else if boilerRunningNum[j] == 3 {
				threeBoliersNum++
			} else if boilerRunningNum[j] == 4 {
				fourBoilersNum++
			}
		}
		C_stage_1D[0], C_stage_1D[1], C_stage_1D[2], C_stage_1D[3], C_stage_1D[4] = zeroBoilerNum, oneBoilerNum, twoBoilersNum, threeBoliersNum, fourBoilersNum
		zeroBoilerNum, oneBoilerNum, twoBoilersNum, threeBoliersNum, fourBoilersNum = 0, 0, 0, 0, 0

		for j := energyDaily.Vally_cost_time_end + 1; j < len(boilerRunningNum); j++ {
			if boilerRunningNum[j] == 0 {
				zeroBoilerNum++
			} else if boilerRunningNum[j] == 1 {
				oneBoilerNum++
			} else if boilerRunningNum[j] == 2 {
				twoBoilersNum++
			} else if boilerRunningNum[j] == 3 {
				threeBoliersNum++
			} else if boilerRunningNum[j] == 4 {
				fourBoilersNum++
			}
		}
		AB_stage_1D[0], AB_stage_1D[1], AB_stage_1D[2], AB_stage_1D[3], AB_stage_1D[4] = zeroBoilerNum, oneBoilerNum, twoBoilersNum, threeBoliersNum, fourBoilersNum
		zeroBoilerNum, oneBoilerNum, twoBoilersNum, threeBoliersNum, fourBoilersNum = 0, 0, 0, 0, 0

		for j := 0; j < 5; j++ {
			C_stage[i][j] = C_stage_1D[j]
			AB_stage[i][j] = AB_stage_1D[j]
		}

	}

	return C_stage, AB_stage
}

/*
 *功能：节约能耗
 *输入：无
 *输出：逐天能耗
 */
func (energyWeekly EnergyConfigWeekly) GetEnergySaving() [7]float64 {
	load_flat_and_peak_sum := 0.0
	var energySaving [7]float64
	for i := 0; i < 7; i++ {
		for j := energyWeekly.Flat_cost_time_1_start + 1; j <= energyWeekly.Flat_cost_time_3_end; j++ {
			load_flat_and_peak_sum = load_flat_and_peak_sum + energyWeekly.Week_load_prediction[i][j]
		}
		if (energyWeekly.Qs - load_flat_and_peak_sum) > 0 {
			energySaving[i] = (1 - energyWeekly.Heat_loss_coefficiency) * (energyWeekly.Qs - load_flat_and_peak_sum)
		} else {
			energySaving[i] = 0
		}
		load_flat_and_peak_sum = 0
	}
	return energySaving
}

/*
 *功能：获取碳排放量
 *输入：无
 *输出：逐天碳排放量
 */
func (energyWeekly EnergyConfigWeekly) GetCarbonEmission(energySaving [7]float64) [7]float64 {
	var carbonEmission [7]float64
	for i := 0; i < 7; i++ {
		carbonEmission[i] = energyWeekly.Carbon_emission_unit_power * energySaving[i]
	}
	return carbonEmission
}

/*
 *功能：获取运行费用
 *输入：无
 *输出：逐天运行费用
 */
func (energyWeekly EnergyConfigWeekly) GetRunningCost() [7]float64 {
	var runningCost, heatStorage [7]float64
	load_flat_sum, load_peak_sum, load_vally_sum := 0.0, 0.0, 0.0
	heatStorage = energyWeekly.GetHeatStorageAagin()

	for i := 0; i < 7; i++ {
		for j := energyWeekly.Peak_cost_time_1_start + 1; j <= energyWeekly.Peak_cost_time_1_end; j++ {
			load_peak_sum = load_peak_sum + energyWeekly.Week_load_prediction[i][j]
		}
		for j := energyWeekly.Peak_cost_time_2_start + 1; j <= energyWeekly.Peak_cost_time_2_end; j++ {
			load_peak_sum = load_peak_sum + energyWeekly.Week_load_prediction[i][j]
		}
		for j := energyWeekly.Flat_cost_time_1_start + 1; j <= energyWeekly.Flat_cost_time_1_end; j++ {
			load_flat_sum = load_flat_sum + energyWeekly.Week_load_prediction[i][j]
		}
		for j := energyWeekly.Flat_cost_time_2_start + 1; j <= energyWeekly.Flat_cost_time_2_end; j++ {
			load_flat_sum = load_flat_sum + energyWeekly.Week_load_prediction[i][j]
		}
		for j := energyWeekly.Flat_cost_time_3_start + 1; j <= energyWeekly.Flat_cost_time_3_end; j++ {
			load_flat_sum = load_flat_sum + energyWeekly.Week_load_prediction[i][j]
		}
		if energyWeekly.Vally_cost_time_start == 23 {
			for j := 0; j <= energyWeekly.Vally_cost_time_end; j++ {
				load_vally_sum = load_vally_sum + energyWeekly.Week_load_prediction[i][j]
			}
		} else {
			for j := 0; j <= energyWeekly.Vally_cost_time_end; j++ {
				load_vally_sum = load_vally_sum + energyWeekly.Week_load_prediction[i][j]
			}
			for j := energyWeekly.Flat_cost_time_3_start + 1; j <= 23; j++ {
				load_vally_sum = load_vally_sum + energyWeekly.Week_load_prediction[i][j]
			}
		}
		if (energyWeekly.Qs - load_peak_sum) > 0 {
			if (energyWeekly.Qs - load_peak_sum - load_flat_sum) > 0 {
				runningCost[i] = energyWeekly.Vally_cost * (load_vally_sum + heatStorage[i])
			} else {
				runningCost[i] = energyWeekly.Flat_cost*(load_flat_sum+load_peak_sum-energyWeekly.Qs) + energyWeekly.Vally_cost*(load_vally_sum+heatStorage[i])
			}
		} else {
			runningCost[i] = energyWeekly.Peak_cost*(load_peak_sum-energyWeekly.Qs) + energyWeekly.Flat_cost*load_flat_sum + energyWeekly.Vally_cost*(load_vally_sum+heatStorage[i])
		}
		load_flat_sum, load_peak_sum, load_vally_sum = 0.0, 0.0, 0.0
	}
	return runningCost
}
