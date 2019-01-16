package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	_ "github.com/heroku/x/hmetrics/onload"
	log "github.com/sirupsen/logrus"

	"github.com/lnquy/line-catalyst-server/internal/bot"
	"github.com/lnquy/line-catalyst-server/pkg/middleware"
)

const (
	secretEnv string = "LINE_BOT_CHANNEL_SECRET"
	tokenEnv  string = "LINE_BOT_CHANNEL_TOKEN"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	token, secret := "", ""
	secret = os.Getenv(secretEnv)
	token = os.Getenv(tokenEnv)
	if secret == "" || token == "" {
		log.Errorf("main: failed to load bot's secret and token from environment")
	}

	catalyst, err := bot.NewCatalyst(secret, token)
	if err != nil {
		log.Panicf("main: failed to create Catalyst bot: %v", err)
	}

	r := chi.NewRouter()
	r.Get("/line", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	r.Post("/line/callback", middleware.
		ValidateLineSignature(secret, http.HandlerFunc(catalyst.MessageHandler)).
		ServeHTTP,
	)

	log.Infof("main: server starting at %s", port)
	if err := http.ListenAndServe(":" + port, r); err != nil && err != http.ErrServerClosed {
		log.Panicf("main: failed to start server on %s: %v", port, err)
	}
	log.Info("main: exit")
}
