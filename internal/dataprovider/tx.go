package dataprovider

import (
	"context"

	"github.com/gxravel/bus-routes/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Tx struct {
	*sqlx.Tx
}

type Txer interface {
	New() (*Tx, error)
}

func EndTransaction(tx *Tx, logger logger.Logger, err error) error {
	if err == nil {
		if cerr := tx.Commit(); cerr != nil {
			return errors.Wrap(cerr, "can not commit transaction")
		}
		logger.Debug("commit OK")

		return nil
	}

	logger.WithErr(err).Debug("rolling back transaction")
	if rerr := tx.Rollback(); rerr != nil {
		logger.WithErr(rerr).Error("can't roll back transaction")
	}

	return err
}

func BeginAutoCommitedTx(ctx context.Context, txer Txer, f func(*Tx) error) error {
	logger := logger.FromContext(ctx)
	logger.Debug("begin transaction")
	tx, err := txer.New()
	if err != nil {
		return errors.Wrap(err, "can't create new transaction")
	}
	err = f(tx)
	return EndTransaction(tx, logger, err)
}
