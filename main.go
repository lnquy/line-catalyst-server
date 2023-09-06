package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/go-chi/chi"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/lnquy/line-catalyst-server/internal/bot"
	"github.com/lnquy/line-catalyst-server/internal/config"
	"github.com/lnquy/line-catalyst-server/internal/repo"
	"github.com/lnquy/line-catalyst-server/pkg/middleware"
	"github.com/lnquy/line-catalyst-server/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadEnvConfig()
	logPanic(err, "main: failed to load configurations")

	lvl, err := log.ParseLevel(cfg.LogLevel)
	logPanic(err, "main: failed to parse log level")
	log.SetLevel(lvl)
	log.Infof("main: configuration: %s", utils.ToJSON(cfg))

	glbCtx, glbCtxCancel := context.WithCancel(context.Background())
	_ = glbCtx

	var messageRepo repo.MessageRepository
	var userRepo repo.UserRepository
	var schedRepo repo.ScheduleRepository
	switch strings.ToLower(cfg.Database.Type) {
	case "mongodb":
		// di, err := mgo.ParseURL(cfg.Database.MongoDB.URI)
		// logPanic(err, "main: invalid MongoDB URI")
		// tlsConfig := &tls.Config{}
		// di.Timeout = 10 * time.Second
		// di.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		// 	conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		// 	return conn, err
		// }
		// session, err := mgo.DialWithInfo(di)
		session, err := mgo.Dial(cfg.Database.MongoDB.URI)
		logPanic(err, "main: failed to dial mongodb")

		// client, err := mongo.Connect(glbCtx, options.Client().ApplyURI(cfg.Database.MongoDB.URI))
		// logPanic(err, "main: failed to dial mongodb")
		// defer func() {
		// 	logPanic(client.Disconnect(glbCtx), "main: failed to gracefully shutdown MongoDB connection")
		// }()
		// db := client.Database("catalyst")

		messageRepo = repo.NewMessageMongoDBRepo(session)
		err = messageRepo.EnsureIndex()
		logPanic(err, "main: failed to ensure database index")

		userRepo = repo.NewUserMongoDBRepo(session)
		err = userRepo.EnsureIndex()
		logPanic(err, "main: failed to ensure database index")
		schedRepo = repo.NewScheduleMongoDBRepo(session)
		err = schedRepo.EnsureIndex()
		logPanic(err, "main: failed to ensure database index")
	default:
		log.Panicf("main: unsupported database type: %s", cfg.Database.Type)
	}

	catalyst, err := bot.NewCatalyst(cfg.Bot, messageRepo, userRepo, schedRepo)
	logPanic(err, "main: failed to create Catalyst bot")

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello world!"))
	})
	r.Post("/line/callback",
		middleware.Recovery(
			middleware.ValidateLineSignature(cfg.Bot.Secret, http.HandlerFunc(catalyst.MessageHandler)),
		).ServeHTTP,
	)

	log.Infof("main: server starting at %s", cfg.Port)
	if err := http.ListenAndServe("0.0.0.0:"+cfg.Port, r); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Panicf("main: failed to start server on %s: %v", cfg.Port, err)
	}
	glbCtxCancel()
	log.Info("main: exit")
}

func logPanic(err error, msg string) {
	if err != nil {
		log.Panicf(msg+": %v", err)
	}
}
