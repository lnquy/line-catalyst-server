package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const lineSignatureHeader = "X-Line-Signature"

func ValidateLineSignature(secret string, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		sign := r.Header.Get(lineSignatureHeader)
		if sign == "" {
			log.Errorf("middleware: request signature must be specified")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			r.Body.Close()
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(sign)
		if err != nil {
			log.Errorf("middleware: failed to decode request signature: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("middleware: failed to read request body: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		hash := hmac.New(sha256.New, []byte(secret))
		hash.Write(body)
		if !hmac.Equal(decoded, hash.Sum(nil)) {
			log.Errorf("middleware: invalid request signature: %s", decoded)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // Pass body data to next handler
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Errorf("panic recovered: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
