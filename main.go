package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/go-chi/chi"
	_ "github.com/heroku/x/hmetrics/onload"
	log "github.com/sirupsen/logrus"

	"github.com/lnquy/line-catalyst-server/internal/bot"
	"github.com/lnquy/line-catalyst-server/internal/config"
	"github.com/lnquy/line-catalyst-server/internal/repo"
	"github.com/lnquy/line-catalyst-server/pkg/middleware"
	"github.com/lnquy/line-catalyst-server/pkg/utils"
)

func main() {
	cfg, err := config.LoadEnvConfig()
	logPanic(err, "main: failed to load configurations")

	lvl, err := log.ParseLevel(cfg.LogLevel)
	logPanic(err, "main: failed to parse log level")
	log.SetLevel(lvl)
	log.Infof("main: configuration: %s", utils.ToJSON(cfg))

	var messageRepo repo.MessageRepository
	var userRepo repo.UserRepository
	switch strings.ToLower(cfg.Database.Type) {
	case "mongodb":
		session, err := mgo.DialWithTimeout(cfg.Database.MongoDB.URI, 30*time.Second)
		logPanic(err, "main: failed to dial mongodb")
		messageRepo = repo.NewMessageMongoDBRepo(session)
		err = messageRepo.EnsureIndex()
		logPanic(err, "main: failed to ensure database index")
		userRepo = repo.NewUserMongoDBRepo(session)
		err = messageRepo.EnsureIndex()
		logPanic(err, "main: failed to ensure database index")
	default:
		log.Panicf("main: unsupported database type: %s", cfg.Database.Type)
	}

	catalyst, err := bot.NewCatalyst(cfg.Bot, messageRepo, userRepo)
	logPanic(err, "main: failed to create Catalyst bot")

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

func logPanic(err error, msg string) {
	if err != nil {
		log.Panicf(msg+": %v", err)
	}
}
