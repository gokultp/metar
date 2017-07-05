package metar

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// DataURL is the source of metar data
	DataURL string = "http://tgftp.nws.noaa.gov/data/observations/metar/stations/"
	// DataExtn is the extension of the data in URL
	DataExtn       string = ".TXT"
	unhandledError string = "Something went wrong while parsing the data."
)

// Wind is  wind object encapsulates its Direction and Speed
type Wind struct {
	Direction int `json:"direction"`
	Velocity  int `json:"velocity_in_knot"`
	Gust      int `json:"gust_in_speed"`
}

// Metar is the metar parsed object
type Metar struct {
	Station         string    `json:"station"`
	LastObservation time.Time `json:"last_observation"`
	Wind            *Wind     `json:"wind"`
	Temperature     *int      `json:"atmospheric_temperature"`
	DewPoint        *int      `json:"dewpoint"`
}

// GetData will get metar data in format from string
func GetData(station string) (*Metar, error) {
	data, err := getDataFromSource(station)

	if err != nil {
		return nil, err
	}
	return serializeMetar(data, station)

}

func getDataFromSource(station string) (string, error) {
	res, err := http.Get(DataURL + station + DataExtn)
	if err != nil {
		return "", err
	}
	if res.StatusCode >= 400 {
		return "", errors.New("Station not found")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	return string(body), nil
}

func serializeMetar(data, station string) (*Metar, error) {
	metarObj := Metar{}
	var err error
	data1 := strings.Split(data, "\n")
	metarObj.LastObservation, err = getDate(data1[0])
	if err != nil {
		return nil, errors.New(unhandledError)
	}
	metarObj.Station = station
	data2 := strings.Split(data1[1], " ")
	metarObj.Wind, err = getWind(&data2)
	if err != nil {
		return nil, errors.New(unhandledError)
	}
	metarObj.Temperature, metarObj.DewPoint, err = getTemperature(&data2)
	if err != nil {
		return nil, errors.New(unhandledError)
	}
	return &metarObj, nil
}

func getDate(data string) (time.Time, error) {
	data = strings.Replace(data, "/", "-", 2)
	data = strings.Replace(data, " ", "T", 1)
	data += ":00Z"
	return time.Parse(time.RFC3339, data)
}

func getWind(data *[]string) (*Wind, error) {
	var windString string
	var err error

	for i, val := range *data {
		if strings.Contains(val, "KT") {
			windString = val
			data = removeFromArray(*data, i)
			break
		}
	}

	if windString == "" {
		return nil, nil
	}
	windString = strings.Replace(windString, "KT", "", 1)
	wind := Wind{}
	wind.Direction, err = strconv.Atoi(windString[0:3])
	if strings.Contains(windString, "G") {
		windVelocityStrs := strings.Split(windString[3:], "G")
		wind.Gust, err = strconv.Atoi(windVelocityStrs[0])

		wind.Velocity, err = strconv.Atoi(windVelocityStrs[1])

	} else {
		wind.Velocity, err = strconv.Atoi(windString[3:])
	}

	return &wind, err

}

func removeFromArray(data []string, i int) *[]string {
	data[i] = data[len(data)-1]
	newData := data[:len(data)-1]
	return &newData
}

func getTemperature(data *[]string) (*int, *int, error) {
	var tempString string

	for i, val := range *data {
		if strings.Contains(val, "/") {
			tempString = val
			data = removeFromArray(*data, i)

			break
		}
	}
	if tempString == "" {
		return nil, nil, nil
	}

	temperatures := strings.Split(tempString, "/")
	sign := 1
	if strings.Contains(temperatures[0], "M") {
		sign = -1
	}
	atm, err := strconv.Atoi(strings.Replace(temperatures[0], "M", "", 1))

	if err != nil {
		return nil, nil, err
	}
	atm *= sign
	sign = 1
	if strings.Contains(temperatures[1], "M") {
		sign = -1
	}

	dewPoint, err := strconv.Atoi(strings.Replace(temperatures[1], "M", "", 1))
	if err != nil {
		return nil, nil, err
	}
	dewPoint *= sign

	return &atm, &dewPoint, nil

}
