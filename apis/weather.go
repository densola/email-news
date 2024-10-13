package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WeatherResponse struct {
	Forecast Forecast `json:"forecast"`
}

type Forecast struct {
	ForecastDay []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Day Day `json:"day"`
}

type Day struct {
	MinTempF          float64   `json:"mintemp_f"`
	MaxTempF          float64   `json:"maxtemp_f"`
	MinTempC          float64   `json:"mintemp_c"`
	MaxTempC          float64   `json:"maxtemp_c"`
	Condition         Condition `json:"condition"`
	DailyChanceOfRain float64   `json:"daily_chance_of_rain"`
}

type Condition struct {
	Text string `json:"text"`
}

func (n *News) getWeather(apiKey, location string) error {
	r, err := http.Get(fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?key=%s&q=%s", apiKey, location))
	if err != nil {
		return fmt.Errorf("getting api response: %w", err)
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	var w WeatherResponse
	err = json.Unmarshal([]byte(string(body)), &w)
	if err != nil {
		return fmt.Errorf("parsing json body: %w", err)
	}

	return nil
}
