package main

import (
	"flag"
	"fmt"

	"github.com/gxravel/bus-routes/internal/config"
	log "github.com/gxravel/bus-routes/internal/logger"
)

func main() {
	defaultLog := log.Default()

	var configPath = flag.String("config", "./config.example.json", "path to configuration file")

	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		defaultLog.WithErr(err).Fatal("can't create config")
	}

	fmt.Println("starting here")
}
