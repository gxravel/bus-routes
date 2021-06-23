package database

import (
	"context"

	"github.com/gxravel/bus-routes/internal/config"
	"github.com/gxravel/bus-routes/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Client struct {
	*sqlx.DB
	log logger.Logger

	schemaName string
}

func NewClient(cfg config.Config, logger logger.Logger) (*Client, error) {
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

	// Migrate up
	if err := m.Migrate(db.DB.DB); err != nil {
		return errors.Wrap(err, "can't migrate the db")
	}

	return nil
}

const statusCheckQuery = `SELECT true`

func (db *Client) StatusCheck(ctx context.Context) error {
	var tmp bool
	return db.QueryRowContext(ctx, statusCheckQuery).Scan(&tmp)
}
