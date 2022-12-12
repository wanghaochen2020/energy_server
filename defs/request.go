package defs

type BasicDataSet struct {
	BasicData         []string `json:"basic_data"`
	BasicDataListDay  []string `json:"basic_data_list_day"`
	BasicDataListHour []string `json:"basic_data_list_hour"`
	BasicOpcList      []string `json:"basic_opc_list"`
}

type BasicDataSetRequest struct {
	Data    BasicDataSet `json:"data"`
	DayStr  string       `json:"day_str"`
	HourStr string       `json:"hour_str"`
}
