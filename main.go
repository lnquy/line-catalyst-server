package main

import (
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/heroku/x/hmetrics/onload"
	log "github.com/sirupsen/logrus"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	log.Infof("Main: server starting at 8080")
	if err := http.ListenAndServe(":8080", r); err != nil && err != http.ErrServerClosed {
		log.Panicf("Main: failed to start server on 8080: %v", err)
	}
	log.Info("Main: exit")
}
