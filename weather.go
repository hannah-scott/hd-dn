package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"math"
)

type Forecast struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	HourlyUnits          struct {
		Time                     string `json:"time"`
		Temperature2M            string `json:"temperature_2m"`
		PrecipitationProbability string `json:"precipitation_probability"`
	} `json:"hourly_units"`
	Hourly struct {
		Time                     []string  `json:"time"`
		Temperature2M            []float64 `json:"temperature_2m"`
		PrecipitationProbability []float64     `json:"precipitation_probability"`
	} `json:"hourly"`
}

type DailyForecast struct {
	Date			string
	Forecast 	struct {
		TempAvg				int
		TempMax				int
		TempMin				int
		ChanceOfRain 	int
	}
}

type ForecastCondensed struct {
	Forecasts []DailyForecast
}

func getDays(times []string) []string {
	dates := []string{}

	for _, v := range times {
		if v[len(v) - 5:len(v)] == "00:00" {
			dates = append(dates, v[0:10])
		}
	}

	return dates
}

func mean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	var sum float64
	for _, d := range data {
		sum += d
	}
	return sum / float64(len(data))
}

func max(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	max := math.Inf(-1)
	for _, d := range data {
		max = math.Max(max, d)
	}
	return max
}

func min(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	min := math.Inf(1)
	for _, d := range data {
		min = math.Min(min, d)
	}
	return min
}

func main() {
	url := "https://api.open-meteo.com/v1/forecast?latitude=51.5085&longitude=-0.1257&hourly=temperature_2m,precipitation_probability"
	// url := "https://api.open-meteo.com/v1/forecast?latitude=51.5085&longitude=-0.1257&hourly=temperature_2m"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	forecast := Forecast{}
	
	err = json.Unmarshal(body, &forecast)
	if err != nil {
		log.Fatal(err)
	}

	hourly_forecast := forecast.Hourly

	fc := ForecastCondensed{}

	for i, d := range getDays(hourly_forecast.Time) {
		df := DailyForecast{}
		df.Date = d

		temps := hourly_forecast.Temperature2M[i*24 : i*24 + 23]
		precps := hourly_forecast.PrecipitationProbability[i*24 : i*24 + 23]

		df.Forecast.TempAvg = int(math.Round(mean(temps)))
		df.Forecast.TempMax = int(math.Round(max(temps)))
		df.Forecast.TempMin = int(math.Round(min(temps)))
		df.Forecast.ChanceOfRain = int(math.Round(mean(precps)))
		fc.Forecasts = append(fc.Forecasts, df)
	}

	fmt.Println(fc)
}
