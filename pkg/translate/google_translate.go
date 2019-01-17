package translate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const googleTranslateURL = "https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s"

func Translate(sourceLang, targetLang, text string) (string, error) {
	gtURL := fmt.Sprintf(googleTranslateURL, sourceLang, targetLang, url.QueryEscape(text))
	resp, err := http.Get(gtURL)
	defer resp.Body.Close()
	if err != nil {
		return "", errors.Wrapf(err, "failed to query Google Translate")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read from Google Translate response body")
	}
	return string(b), nil
}
