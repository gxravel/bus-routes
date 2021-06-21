package main

import (
	"context"
	"flag"

	"github.com/gxravel/bus-routes/internal/config"
	"github.com/gxravel/bus-routes/internal/database"
	"github.com/gxravel/bus-routes/internal/dataprovider/mysql"
	"github.com/gxravel/bus-routes/internal/logger"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	defaultLog := logger.Default()

	var configPath = flag.String("config", "./config.example.json", "path to configuration file")

	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		defaultLog.WithErr(err).Fatal("can't create config")
	}

	log, err := logger.New(
		cfg.Log.Level,
		logger.DefaultOutput)
	if err != nil {
		defaultLog.WithErr(err).Error("can not init log with specified params, defaults are used")
		log = defaultLog
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.CtxWithLogger(ctx, log)

	db, err := database.NewClient(*cfg, log)
	if err != nil {
		log.WithErr(err).Fatal("constructing database")
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.WithErr(err).Error("closing database connection")
		}
	}()

	if err := db.Migrate(); err != nil {
		log.WithErr(err).Fatal("can't migrate the db")
	}

	txer := mysql.NewTxManager(db)
	busStore := mysql.NewBusStore(db, txer)
}
