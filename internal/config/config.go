package config

import (
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gxravel/bus-routes/internal/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Environment     string        `mapstructure:"environment"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`

	API api `mapstructure:"api"`
	DB  DB  `mapstructure:"db"`
	Log Log `mapstructure:"logger"`
}

type api struct {
	ServeSwagger bool          `mapstructure:"serve_swagger"`
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DB struct {
	URL          string `mapstructure:"url"`
	SchemaName   string `mapstructure:"schema_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

var defaults = map[string]interface{}{
	"environment":      "development",
	"shutdown_timeout": time.Second * 5,

	"db.url":            "gxravel:gxravel@tcp(localhost:3308)",
	"db.schema_name":    "bus_routes",
	"db.max_open_conns": 2,
	"db.max_idle_conns": 2,

	"api.serve_swagger": true,
	"api.address":       ":8090",
	"api.read_timeout":  time.Second * 5,
	"api.write_timeout": time.Second * 5,

	"logger.level":  "debug",
	"logger.format": "json",
}

func New(dst string) (*Config, error) {
	dir, basename := filepath.Split(dst)
	name := strings.TrimSuffix(basename, path.Ext(basename))

	viper.AddConfigPath(dir)
	viper.SetConfigName(name)

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	var c Config

	if err := viper.ReadInConfig(); err != nil {
		logger.Default().WithErr(err).Errorf("could not read config, using defaults: %v", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
