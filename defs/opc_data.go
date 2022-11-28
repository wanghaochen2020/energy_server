package defs

//计算结果表名
const (
	EnergyOnlineRate             = "energy_online_rate"              //能源站设备在线率
	EnergyBoilerPower            = "energy_boiler_power"             //能源站锅炉功率
	EnergyPowerConsumptionToday  = "energy_power_consumption_today"  //能源站今日总耗能
	EnergyBoilerRunningNum       = "energy_boiler_running_num"       //能源站锅炉运行数目
	EnergyTankRunningNum         = "energy_tank_running_num"         //能源站蓄热水箱运行数目
	EnergyHeatSupplyToday        = "energy_heat_supply_today"        //能源站今日总供热量
	EnergyHeatStorageAndRelease  = "energy_heat_storage_and_release" //能源站每日各小时水箱蓄放热量
	EnergyBoilerEnergyCost       = "energy_boiler_energy_cost"       //能源站每日各小时锅炉能耗
	EnergyBoilerEfficiencyDay    = "energy_boiler_efficiency_day"    //能源站每日各小时锅炉效率
	EnergyWatertankEfficiencyDay = "energy_watertank_efficiency_day" //能源站每日各小时蓄热水箱效率
	EnergyEfficiencyDay          = "energy_efficiency_day"           //能源站每日各小时效率
	EnergyCarbonDay              = "energy_carbon_day"               //能源站每日各小时碳排
	EnergyCarbonMonth            = "energy_carbon_month"             //能源站每月各天碳排总和
	EnergyCarbonYear             = "energy_carbon_year"              //能源站每年各月碳排总和
	EnergyBoilerPayloadDay       = "energy_boiler_payload_day"       //能源站每日各小时锅炉负载
	EnergyBoilerPayloadMonth     = "energy_boiler_payload_month"     //能源站每月各天平均锅炉负载
	EnergyBoilerPayloadYear      = "energy_boiler_payload_year"      //能源站每年各月平均锅炉负载

	ColdPowerMin              = "cold_power_min"               //制冷中心当分钟功率
	ColdEnergyCostToday       = "cold_energy_cost_today"       //制冷中心今日能耗
	ColdMachineRunningNum     = "cold_machine_running_num"     //制冷中心制冷机运行台数
	ColdCoolingWaterInT       = "cold_cooling_water_InT"       //制冷中心冷却水进水温度
	ColdCoolingWaterOutT      = "cold_cooling_water_OutT"      //制冷中心冷却水出水温度
	ColdRefrigeratedWaterInT  = "cold_refrigerated_water_InT"  //制冷中心冷冻水进水温度
	ColdRefrigeratedWaterOutT = "cold_refrigerated_water_OutT" //制冷中心冷冻水出水温度
	ColdMachinePowerMin       = "cold_machine_power_min"       //制冷中心制冷机当分钟功率
	ColdEnergyCostDay         = "cold_energy_cost_day"         //制冷中心每日各小时能耗
	ColdCarbonDay             = "cold_carbon_day"              //制冷中心每日各小时碳排
	ColdCarbonMonth           = "cold_carbon_month"            //制冷中心每月各天碳排总和
	ColdCarbonYear            = "cold_carbon_year"             //制冷中心每年各月碳排总和

	PumpPowerMin      = "pump_power_min"       //二次泵站功率
	PumpPowerToday    = "pump_power_today"     //二次泵站今日能耗
	PumpEnergyCostDay = "pump_energy_cost_day" //二次泵站每日各小时能耗
	PumpRunningState1 = "pump_running_state1"  //二次泵站泵运行状态
	PumpRunningState2 = "pump_running_state2"  //二次泵站泵运行状态
	PumpRunningState3 = "pump_running_state3"  //二次泵站泵运行状态
	PumpRunningState4 = "pump_running_state4"  //二次泵站泵运行状态
	PumpRunningState5 = "pump_running_state5"  //二次泵站泵运行状态
	PumpRunningState6 = "pump_running_state6"  //二次泵站泵运行状态
	PumpHeatHour1     = "pump_heat_hour1"      //二次泵站当小三每分钟环路1输热量
	PumpHeatHour2     = "pump_heat_hour2"      //二次泵站当小三每分钟环路2输热量
	PumpEHR1          = "pump_EHR1"            //二次泵站环路1每日EHR
	PumpEHR2          = "pump_EHR2"            //二次泵站环路2每日EHR

	SolarWaterBoilerPowerConsumptionToday = "solar_water_boiler_power_comsumption_today" //太阳能热水电加热器今日总耗电量
	SolarWaterHeatCollecterInT            = "solar_water_heat_collecter_in_temp"         //太阳能热水集热器进口温度
	SolarWaterHeatCollecterOutT           = "solar_water_heat_collecter_out_temp"        //太阳能热水集热器出口温度
	SolarWaterJRQT                        = "solar_water_JRQ_temp"                       //太阳能热水加热器温度
	SolarWaterHeatCollectionToday         = "solar_water_heat_collection_today"          //太阳能热水今日总集热量
	SolarWaterPumpRunningNum              = "solar_water_pump_running_num"               //太阳能热水水泵运行数目
	SolarWaterHeatCollectionHour          = "solar_water_heat_collection_hour"           //太阳能热水集热量当小时每分钟
	SolarWaterHeatCollectionDay           = "solar_water_heat_collection_day"            //太阳能热水集热量当日每小时
	SolarWaterHeatCollectionMonth         = "solar_water_heat_collection_month"          //太阳能热水集热量每月各天总和
	SolarWaterHeatCollectionYear          = "solar_water_heat_collection_year"           //太阳能热水集热量每年各月总和
	SolarWaterBoilerPowerConsumptionHour  = "solar_water_boiler_power_comsumption_hour"  //太阳能热水电加热器耗电量当小时每分钟
	SolarWaterBoilerPowerConsumptionDay   = "solar_water_boiler_power_comsumption_day"   //太阳能热水电加热器耗电量当日每小时
	SolarWaterHeatEfficiencyDay           = "solar_water_heat_efficiency_day"            //太阳能热水集热效率当日每小时
	SolarWaterGuaranteeRateDay            = "solar_water_guarantee_rate"                 //太阳能热水保证率当日每小时
)

//其它常数
const (
	ColdMachine1 = "Z_L"
	ColdMachine2 = "Z_LX1"
	ColdMachine3 = "Z_LX2"
)

type OpcData struct {
	Time   string        `bson:"time"`
	ItemId string        `bson:"itemid"`
	Value  []interface{} `bson:"value"`
}

type CalculationResultFloatList struct {
	Time  string    `bson:"time"`
	Name  string    `bson:"name"`
	Value []float64 `bson:"value"`
}

type CalculationResultFloat struct {
	Time  string  `bson:"time"`
	Name  string  `bson:"name"`
	Value float64 `bson:"value"`
}
type OpcUpdateTime struct {
	UpdateTime string `bson:"update_time"`
	Group      string `bson:"group"`
}
