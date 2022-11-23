package defs

//计算结果表名
const (
	EnergyOnlineRate             = "energy_online_rate"              //能源站设备在线率
	EnergyBoilerPower            = "energy_boiler_power"             //能源站锅炉功率
	EnergyPowerConsumptionToday  = "energy_power_consumption_today"  //能源站今日总耗能
	EnergyBoilerRunningNum       = "energy_boiler_running_num"       //能源站锅炉运行数目
	EnergyBoilerEfficiencyDay    = "energy_boiler_efficiency_day"    //能源站每日各小时锅炉效率
	EnergyWatertankEfficiencyDay = "energy_watertank_efficiency_day" //能源站每日各小时蓄热水箱效率
	EnergyEfficiencyDay          = "energy_efficiency_day"           //能源站每日各小时效率
	EnergyCarbonDay              = "energy_carbon_day"               //能源站每日各小时碳排
	EnergyCarbonMonth            = "energy_carbon_month"             //能源站每月各天碳排总和
	EnergyCarbonYear             = "energy_carbon_year"              //能源站每年各月碳排总和
	EnergyBoilerPayloadDay       = "energy_boiler_payload_day"       //能源站每日各小时锅炉负载
	EnergyBoilerPayloadMonth     = "energy_boiler_payload_month"     //能源站每月各天平均锅炉负载
	EnergyBoilerPayloadYear      = "energy_boiler_payload_year"      //能源站每年各月平均锅炉负载

	ColdPowerDay    = "cold_power_day"    //制冷中心每日各小时能耗
	ColdCarbonDay   = "cold_carbon_day"   //制冷中心每日各小时碳排
	ColdCarbonMonth = "cold_carbon_month" //制冷中心每月各天碳排总和
	ColdCarbonYear  = "cold_carbon_year"  //制冷中心每年各月碳排总和

	PumpPowerDay = "pump_power_day" //二次泵站每日各小时能耗
	PumpEHR1     = "pump_EHR1"      //二次泵站环路1每日EHR
	PumpEHR2     = "pump_EHR2"      //二次泵站环路2每日EHR

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
