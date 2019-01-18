package aqi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const replyTmpl = `%s Air Quality Index
-----
AQI: %.0f
Level: %s (%d)
Health Implication: %s
Precaution: %s`


var (
	aqiLevels = []aqiLevel{
		{
			level:       "Good",
			emoji:       ":blush:",
			implication: "Air quality is considered satisfactory, and air pollution poses little or no risk.",
			cautionary:  "Everyone can continue outdoor activities normally.",
		},
		{
			level:       "Moderate",
			emoji:       ":neutral_face:",
			implication: "Air quality is acceptable; however, for some pollutants there may be a moderate health concern for a very small number of people who are unusually sensitive to air pollution.",
			cautionary:  "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion.",
		},
		{
			level:       "Unhealthy for Sensitive Groups",
			emoji:       ":fearful:",
			implication: "Members of sensitive groups may experience health effects. The general public is not likely to be affected.",
			cautionary:  "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion.",
		},
		{
			level:       "Unhealthy",
			emoji:       ":mask:",
			implication: "Everyone may begin to experience health effects; members of sensitive groups may experience more serious health effects",
			cautionary:  "Active children and adults, and people with respiratory disease, such as asthma, should avoid prolonged outdoor exertion; everyone else, especially children, should limit prolonged outdoor exertion.",
		},
		{
			level:       "Very Unhealthy",
			emoji:       ":dizzy_face:",
			implication: "Health warnings of emergency conditions. The entire population is more likely to be affected.",
			cautionary:  "Active children and adults, and people with respiratory disease, such as asthma, should avoid all outdoor exertion; everyone else, especially children, should limit outdoor exertion.",
		},
		{
			level:       "Hazardous",
			emoji:       ":skull_crossbones:",
			implication: "Health alert: everyone may experience more serious health effects.",
			cautionary:  "Everyone should avoid all outdoor exertion.",
		},
	}

	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			DisableKeepAlives:     false,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
)

// https://mholt.github.io/json-to-go/
type (
	aqiLevel struct {
		level       string
		emoji       string
		implication string
		cautionary  string
	}
	aqicnResponse struct {
		Status string `json:"status"`
		Data   struct {
			Aqi          float64 `json:"aqi"`
			Idx          int     `json:"idx"`
			Attributions []struct {
				URL  string `json:"url"`
				Name string `json:"name"`
			} `json:"attributions"`
			City struct {
				Geo  []float64 `json:"geo"`
				Name string    `json:"name"`
				URL  string    `json:"url"`
			} `json:"city"`
			Dominentpol string `json:"dominentpol"`
			Iaqi        struct {
				Co struct {
					V float64 `json:"v"`
				} `json:"co"`
				H struct {
					V float64 `json:"v"`
				} `json:"h"`
				No2 struct {
					V float64 `json:"v"`
				} `json:"no2"`
				O3 struct {
					V float64 `json:"v"`
				} `json:"o3"`
				P struct {
					V float64 `json:"v"`
				} `json:"p"`
				Pm10 struct {
					V float64 `json:"v"`
				} `json:"pm10"`
				Pm25 struct {
					V float64 `json:"v"`
				} `json:"pm25"`
				R struct {
					V float64 `json:"v"`
				} `json:"r"`
				So2 struct {
					V float64 `json:"v"`
				} `json:"so2"`
				T struct {
					V float64 `json:"v"`
				} `json:"t"`
				W struct {
					V float64 `json:"v"`
				} `json:"w"`
			} `json:"iaqi"`
			Time struct {
				S  string  `json:"s"`
				Tz string  `json:"tz"`
				V  float64 `json:"v"`
			} `json:"time"`
			Debug struct {
				Sync time.Time `json:"sync"`
			} `json:"debug"`
		} `json:"data"`
	}
)

func GetAQIInfo(city, token string) (string, error) {
	u := fmt.Sprintf("https://api.waqi.info/feed/%s/?token=%s", city, token)
	resp, err := client.Get(u)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get AQI info from AQICN")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("failed to get AQI info from AQICN: %s", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read response from AQICN")
	}
	status := new(struct {
		Status string `json:"status"`
	})
	_ = json.Unmarshal(b, status)
	if status.Status != "ok" {
		return "", fmt.Errorf("failed to get AQI info from AQICN: status=%s", status.Status)
	}

	var data aqicnResponse
	if err := json.Unmarshal(b, &data); err != nil {
		return "", errors.Wrapf(err, "failed to decode AQICN response data")
	}

	city = strings.Title(city)
	lvl, aqiLvl := getAQILevel(data.Data.Aqi)
	return fmt.Sprintf(replyTmpl,
		city,
		data.Data.Aqi,
		aqiLvl.level, lvl,
		aqiLvl.implication,
		aqiLvl.cautionary,
	), nil
}

func getAQILevel(value float64) (int, *aqiLevel) {
	lvlIdx := 0
	switch {
	case value < 0:
		return -1, nil
	case value <= 50:
		lvlIdx = 0
	case value <= 100:
		lvlIdx = 1
	case value <= 150:
		lvlIdx = 2
	case value <= 200:
		lvlIdx = 3
	case value <= 300:
		lvlIdx = 4
	case value > 300:
		lvlIdx = 5
	default:
		return -1, nil
	}
	return lvlIdx + 1, &aqiLevels[lvlIdx]
}
