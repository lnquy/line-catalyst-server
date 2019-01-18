package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const replyTmpl = `%s Current Weather Report
-----
Temperature: %.1f°C (%.1f°C - %.1f°C)
Humidity: %.0f%%
Pressure: %.0fhPa
Wind speed: %.1fm/s (%.1f°)
Cloudiness: %.0f%%
Sunrise: %v
Sunset: %v`

type openWeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure float64 `json:"pressure"`
		Humidity float64 `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All float64 `json:"all"`
	} `json:"clouds"`
	Dt  float64 `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int64   `json:"sunrise"`
		Sunset  int64   `json:"sunset"`
	} `json:"sys"`
	ID   int     `json:"id"`
	Name string  `json:"name"`
	Cod  float64 `json:"cod"`
}

func GetWeatherInfo(city, token string) (string, error) {
	u := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, token)
	resp, err := http.Get(u)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get weather info from OpenWeatherMap")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("failed to get weather info from OpenWeatherMap: %s", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read response from OpenWeatherMap")
	}

	var data openWeatherResponse
	if err := json.Unmarshal(b, &data); err != nil {
		return "", errors.Wrapf(err, "failed to decode OpenWeatherMap response data")
	}

	w := fmt.Sprintf(replyTmpl,
		strings.Title(city),
		kelvinToCelsius(data.Main.Temp), kelvinToCelsius(data.Main.TempMin), kelvinToCelsius(data.Main.TempMax),
		data.Main.Humidity,
		data.Main.Pressure,
		data.Wind.Speed, data.Wind.Deg,
		data.Clouds.All,
		time.Unix(data.Sys.Sunrise, 0).Format("15:04"),
		time.Unix(data.Sys.Sunset, 0).Format("15:04"),
	)
	return w, nil
}

func kelvinToCelsius(k float64) float64 {
	return k - 273.15
}
