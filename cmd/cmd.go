package main

import (
	"flag"
	"fmt"

	"github.com/gxravel/bus-routes/internal/config"
	"github.com/gxravel/bus-routes/internal/logger"
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

	fmt.Println("starting here")
}
