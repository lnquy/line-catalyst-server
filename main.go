package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	_ "github.com/heroku/x/hmetrics/onload"
	log "github.com/sirupsen/logrus"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	log.Infof("Main: server starting at %s", port)
	if err := http.ListenAndServe(":" + port, r); err != nil && err != http.ErrServerClosed {
		log.Panicf("Main: failed to start server on %s: %v", port, err)
	}
	log.Info("Main: exit")
}
