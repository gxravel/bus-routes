package mysql

import (
	"github.com/gxravel/bus-routes/internal/database"
	"github.com/gxravel/bus-routes/internal/dataprovider"

	"github.com/pkg/errors"
)

func NewTxManager(db *database.Client) dataprovider.Txer {
	return &TxManager{
		db: db,
	}
}

type TxManager struct {
	db *database.Client
}

func (txm *TxManager) New() (*dataprovider.Tx, error) {
	sqltx, err := txm.db.Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "creating tx")
	}

	return &dataprovider.Tx{Tx: sqltx}, nil
}
