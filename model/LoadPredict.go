package model

import (
	"context"
	"encoding/json"
	"energy/defs"
	"energy/utils"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	getWeatherForecastUrl = "https://api.qweather.com/v7/weather/168h?location=101010800&key=2dee7efdb9a54d06830b1c3af13857db"
)

type Input struct {
	Date          [168]string  `json:"日期"`
	Temperature   [168]float64 `json:"温度"`
	Humidity      [168]int     `json:"湿度"`
	Radiation     [168]int     `json:"辐射"`
	Wind          [168]float64 `json:"风速"`
	RoomRate      [168]float64 `json:"在室率"`
	OccupancyRate [168]float64 `json:"入住率"`
	Load          [168]float64 `json:"负荷"`
}

type Output struct {
	Result [168]float64 `json:"result"`
}

type Forecast struct {
	Hourly []Forecast2 `json:"hourly"`
}

type Forecast2 struct {
	Temp      string `json:"temp"`
	Humidity  string `json:"humidity"`
	WindSpeed string `json:"windSpeed"`
}

type Atmosphere struct {
	Time   int         `json:"time"`
	Result Atmosphere2 `json:"data"`
}

type Atmosphere2 struct {
	Wind        Atmosphere3
	Temperature Atmosphere3
	Humidity    Atmosphere3
	Radiation   Atmosphere3
}

type Atmosphere3 struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Kekong struct {
	Time int     `json:"time"`
	D4   float64 `json:"d4"`
	D5   float64 `json:"d5"`
	D6   float64 `json:"d6"`
}

func LoadPredict(index string) Output {
	input := MakeInputBody(index)

	data := Output{}
	resp, err := http.Post(utils.LoadPredictRouter, "application/json", strings.NewReader(string(input)))
	if err != nil {
		log.Println(err)
		return Output{}
	}
	defer resp.Body.Close()
	n, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(n, &data)
	return data
}

func MakeInputBody(index string) []byte {
	var input Input
	start := FindStart(int(time.Now().Unix()))
	for i := 0; i < 168; i++ {
		input.Date[i] = UnixToString(start + i*3600)
	}

	//前一天数据
	load := GetLoad(index, "yesterday")
	temperature := GetData("temperature", int(time.Now().Unix()-86400), index)
	humidity := GetData("humidity", int(time.Now().Unix()-86400), index)
	radiation := GetData("radiation", int(time.Now().Unix()-86400), index)
	wind := GetData("wind", int(time.Now().Unix()-86400), index)
	roomRate := GetData("roomRate", int(time.Now().Unix()-86400), index)
	occupancyRate := GetData("roomRate", int(time.Now().Unix()-86400), index)

	for i := 0; i < 24; i++ {
		input.Load[i] = load[i]
		input.Temperature[i] = temperature[i]
		input.Humidity[i] = int(humidity[i])
		input.Radiation[i] = int(radiation[i])
		input.Wind[i] = wind[i]
		input.RoomRate[i] = roomRate[i]
		input.OccupancyRate[i] = occupancyRate[i]
	}

	//后六天数据
	forecast := GetForecast()

	for i := 24; i < 168; i++ {
		input.Temperature[i], _ = strconv.ParseFloat(forecast.Hourly[i-24].Temp, 64)
		input.Humidity[i], _ = strconv.Atoi(forecast.Hourly[i-24].Temp)
		input.Wind[i], _ = strconv.ParseFloat(forecast.Hourly[i-24].Temp, 64)
		input.Load[i] = 0
		input.Radiation[i] = int(radiation[i%24])
		input.RoomRate[i] = roomRate[i%24]
		input.OccupancyRate[i] = occupancyRate[i%24]
	}

	output, _ := json.Marshal(&input)
	return output
}

//访问办公网数据库
func GetData(index string, base int, zutuan string) []float64 {
	var array []float64
	array = make([]float64, 24)
	limit0, limit1 := FindIntervalDay(base)

	type Update struct {
		Id        int
		UpdatedAt string
		TableName string
	}

	switch index {
	case "roomRate":
		var devices []Kekong
		data, _ := Db.Collection("keKong").Find(context.TODO(), bson.M{"time": bson.M{"$gte": limit0, "$lt": limit1}})
		err = data.All(context.TODO(), &devices)
		if zutuan == "D4组团" {
			for i := 0; i < 24; i++ {
				array[i] = devices[i].D4
			}
		} else if zutuan == "D5组团" {
			for i := 0; i < 24; i++ {
				array[i] = devices[i].D5
			}
		} else if zutuan == "D6组团" {
			for i := 0; i < 24; i++ {
				array[i] = devices[i].D6
			}
		} else {
			for i := 0; i < 24; i++ {
				array[i] = devices[i].D4
			}
		}
	default:
		var devices []Atmosphere
		data, _ := Db.Collection("atmosphere").Find(context.TODO(), bson.M{"time": bson.M{"$gte": limit0, "$lt": limit1}})
		err = data.All(context.TODO(), &devices)
		if index == "temperature" {
			for i := 0; i < 24; i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Result.Temperature.Value, 64)
			}
		} else if index == "humidity" {
			for i := 0; i < 24; i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Result.Humidity.Value, 64)
			}
		} else if index == "radiation" {
			for i := 0; i < 24; i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Result.Radiation.Value, 64)
			}
		} else if index == "wind" {
			for i := 0; i < 24; i++ {
				array[i], _ = strconv.ParseFloat(devices[i].Result.Wind.Value, 64)
			}
		}
	}

	return array
}

func GetLoad(index string, flag string) []float64 {
	var load []float64

	if flag == "today" {
		//str := GetToday()
		str := "2023/02/20"
		switch index {
		case "D1组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay1, str)
		case "D2组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay2, str)
		case "D3组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay3, str)
		case "D4组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay4, str)
		case "D5组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay5, str)
		case "D6组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay6, str)
		case "公共组团南区":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDayPubS, str)
		case "公共组团北区":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDayPubS, str)
		}
	} else if flag == "yesterday" {
		switch index {
		case "D1组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay1, GetYesterday())
		case "D2组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay2, GetYesterday())
		case "D3组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay3, GetYesterday())
		case "D4组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay4, GetYesterday())
		case "D5组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay5, GetYesterday())
		case "D6组团":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDay6, GetYesterday())
		case "公共组团南区":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDayPubS, GetYesterday())
		case "公共组团北区":
			load, _ = GetResultFloatList(defs.GroupHeatConsumptionDayPubS, GetYesterday())
		}
	}
	return load
}

func GetForecast() Forecast {
	data := Forecast{}
	resp, err := http.Get(getWeatherForecastUrl)
	if err != nil {
		log.Println(err)
		return Forecast{}
	}
	defer resp.Body.Close()
	n, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(n, &data)
	return data
}

func UnixToString(unix int) string {
	timeLayout := "2006-01-02 15:04:05"
	timeStr := time.Unix(int64(unix), 0).Format(timeLayout)
	return timeStr
}

func FindStart(value int) int {
	Time := time.Unix(int64(value), 0)
	Time2 := time.Date(Time.Year(), Time.Month(), Time.Day(), 0, 0, 0, 0, Time.Location())
	return int(Time2.Unix())
}

func GetDay(v int64) string {
	timeLayout := "2006-01-02 15:04:05"
	timeStr := time.Unix(v, 0).Format(timeLayout)
	a := strings.Split(timeStr, " ")
	return a[0]
}

func GetToday() string {
	timeLayout := "2006-01-02 15:04:05"
	timeStr := time.Unix(time.Now().Unix(), 0).Format(timeLayout)
	a := strings.Split(timeStr, " ")
	return a[0]
}

func GetYesterday() string {
	timeLayout := "2006-01-02 15:04:05"
	timeStr := time.Unix(time.Now().Unix()-86400, 0).Format(timeLayout)
	a := strings.Split(timeStr, " ")
	return a[0]
}

func FindIntervalDay(value int) (int, int) {
	Time := time.Unix(int64(value), 0)
	Time1 := time.Date(Time.Year(), Time.Month(), Time.Day(), 0, 0, 0, 0, Time.Location())
	Time2 := time.Date(Time.Year(), Time.Month(), Time.Day(), 24, 0, 0, 0, Time.Location())
	return int(Time1.Unix()), int(Time2.Unix())
}
