package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	_ "github.com/heroku/x/hmetrics/onload"
	log "github.com/sirupsen/logrus"

	"github.com/lnquy/line-catalyst-server/internal/bot"
	"github.com/lnquy/line-catalyst-server/internal/config"
	"github.com/lnquy/line-catalyst-server/internal/repo"
	"github.com/lnquy/line-catalyst-server/pkg/middleware"
)

func main() {
	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Panicf("main: failed to load configurations: %v", err)
	}

	var messageRepo repo.MessageRepository
	switch strings.ToLower(cfg.Database.Type) {
	case "mongodb":
		messageRepo, err = repo.NewMessageMongoDBRepo(cfg.Database.MongoDB)
		if err != nil {
			log.Panicf("main: failed to init mongodb: %v", err)
		}
	default:
		log.Panicf("main: unsupported database type: %s", cfg.Database.Type)
	}

	catalyst, err := bot.NewCatalyst(cfg.Bot, messageRepo)
	if err != nil {
		log.Panicf("main: failed to create Catalyst bot: %v", err)
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	r.Post("/line/callback", middleware.
		ValidateLineSignature(cfg.Bot.Secret, http.HandlerFunc(catalyst.MessageHandler)).
		ServeHTTP,
	)

	log.Infof("main: server starting at %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil && err != http.ErrServerClosed {
		log.Panicf("main: failed to start server on %s: %v", cfg.Port, err)
	}
	log.Info("main: exit")
}
