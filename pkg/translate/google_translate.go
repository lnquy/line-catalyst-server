package translate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// [[[\"Financial summary \",\"สรุปฐานะการเงิน\",null,null,3],[\"Industry performance\",\"ผลการดำเนินงานรายอุตสาหกรรม\",null,null,3]],null,\"th\"]
func Translate(sourceLang, targetLang, text string) (string, error) {
	resp, err := http.DefaultClient.Do(getGoogleTranslateAPIRequest(sourceLang, targetLang, text))
	defer resp.Body.Close()
	if err != nil {
		return "", errors.Wrapf(err, "failed to query Google Translate")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("failed to query Google Translate: %s", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read from Google Translate response body")
	}
	logrus.Tracef("trans: %s", string(b))

	gj := gjson.Parse(string(b))
	detectedLang := gj.Get("2").String()
	_ = detectedLang

	items := gj.Get("0").Array()
	translated := ""
	for _, item := range items {
		translated += item.Get("0").String()
	}
	return translated, nil
}

func getGoogleTranslateAPIRequest(sl, tl, text string) *http.Request {
	u, _ := url.Parse("https://translate.googleapis.com/translate_a/single")
	values := url.Values{}
	values.Add("client", "gtx")
	values.Add("q", text)
	values.Add("dt", "t")
	values.Add("hl", "en")
	values.Add("sl", sl)
	values.Add("tl", tl)
	values.Add("ie", "UTF-8")
	values.Add("oe", "UTF-8")
	// values.Add("multires", "1")
	// values.Add("otf", "1")
	// values.Add("pc", "1")
	// values.Add("trs", "1")
	// values.Add("ssel", "3")
	// values.Add("tsel", "6")
	// values.Add("sc", "1")
	u.RawQuery = values.Encode()
	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Charset", "UTF-8")
	return req
}
