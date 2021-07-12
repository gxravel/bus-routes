package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/gxravel/bus-routes/internal/api/http/handler"
	"github.com/gxravel/bus-routes/internal/busroutes"
	"github.com/gxravel/bus-routes/internal/config"
	"github.com/gxravel/bus-routes/internal/database"
	"github.com/gxravel/bus-routes/internal/dataprovider/mysql"
	"github.com/gxravel/bus-routes/internal/jwt"
	log "github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/storage"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	defaultLogger := log.Default()

	var configPath = flag.String("config", "./config.example.json", "path to configuration file")

	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		defaultLogger.WithErr(err).Fatal("can't create config")
	}

	logger, err := log.New(
		cfg.Log.Level,
		log.DefaultOutput,
	)
	if err != nil {
		defaultLogger.WithErr(err).Error("can not init logger with specified params, defaults are used")
		logger = defaultLogger
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.CtxWithLogger(ctx, logger)

	db, err := database.NewClient(*cfg, logger)
	if err != nil {
		logger.WithErr(err).Fatal("construct database")
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.WithErr(err).Error("close database connection")
		}
	}()

	if err := db.Migrate(); err != nil {
		logger.WithErr(err).Fatal("can't migrate the db")
	}

	storage, err := storage.NewClient(*cfg)
	if err != nil {
		logger.WithErr(err).Fatal("connect to storage client")
	}

	txer := mysql.NewTxManager(db)

	busroutes := busroutes.New(
		cfg,
		db,
		logger,
		mysql.NewBusStore(db, txer),
		mysql.NewCityStore(db, txer),
		mysql.NewStopStore(db, txer),
		mysql.NewRouteStore(db, txer),
		mysql.NewUserStore(db, txer),
		txer,
		jwt.New(storage, *cfg),
	)

	apiServer := handler.NewServer(
		cfg,
		busroutes,
		logger,
	)

	var (
		shutdown     = make(chan os.Signal, 1)
		serverErrors = make(chan error, 1)
	)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		serverErrors <- apiServer.ListenAndServe()
	}()

	go func() {
		amqpServer.ListenAndServe()
	}()

	logger.Info("started")

	defer logger.Info("stopped")

	select {
	case err = <-serverErrors:
		logger.WithErr(err).Error("api server stopped")

	case sig := <-shutdown:
		logger.WithField("signal", sig.String()).Info("gracefully shutdown application")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err = apiServer.Shutdown(ctx); err != nil {
			logger.WithErr(err).Error("api server shutdown error")
			err = apiServer.Close()
		}

		if err != nil {
			logger.WithErr(err).Error("could not stopped api server gracefully")
		}
	}
}
