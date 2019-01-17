package weather

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// https://mholt.github.io/json-to-go/
type aqicnResponse struct {
	Status string `json:"status"`
	Data   struct {
		Aqi          int `json:"aqi"`
		Idx          int `json:"idx"`
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
				V int `json:"v"`
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
				V int `json:"v"`
			} `json:"pm10"`
			Pm25 struct {
				V int `json:"v"`
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
			S  string `json:"s"`
			Tz string `json:"tz"`
			V  int    `json:"v"`
		} `json:"time"`
		Debug struct {
			Sync time.Time `json:"sync"`
		} `json:"debug"`
	} `json:"data"`
}

func AQIInfo(city, token string) (string, error) {
	u := fmt.Sprintf("https://api.waqi.info/feed/%s/?token=%s", city, token)
	resp, err := http.Get(u)
	defer resp.Body.Close()

	if err != nil {
		return "", errors.Wrapf(err, "failed to get weather info from AQICN")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("failed to get weather info from AQICN: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read response from AQICN")
	}

	// var data aqicnResponse
	// if err := json.Unmarshal(b, &data); err != nil{
	// 	return "", errors.Wrapf(err, "failed to decode AQICN response data")
	// }
	return string(b), nil
}

// TODO
const replyTmpl = `
`
