package database

import (
	"github.com/gxravel/bus-routes/internal/config"
	log "github.com/gxravel/bus-routes/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Client struct {
	*sqlx.DB
	logger log.Logger

	schemaName string
}

func NewClient(cfg config.Config, logger log.Logger) (*Client, error) {
	db, err := sqlx.Open("mysql", cfg.DB.URL+"/"+cfg.DB.SchemaName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	return &Client{
		db,
		logger,
		cfg.DB.SchemaName,
	}, nil
}

func (db *Client) Migrate() error {
	if _, err := db.Exec(`CREATE SCHEMA IF NOT EXISTS ` + db.schemaName); err != nil {
		return errors.Wrap(err, "can't create schema")
	}

	m, err := migrations(db.schemaName, "migrations")
	if err != nil {
		return errors.Wrap(err, "can't create a new migrator instance")
	}

	if err := m.Migrate(db.DB.DB); err != nil {
		return errors.Wrap(err, "can't migrate the db")
	}

	return nil
}
