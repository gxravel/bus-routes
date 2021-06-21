package busroutes

import (
	"context"

	"github.com/gxravel/bus-routes/internal/config"
	"github.com/gxravel/bus-routes/internal/database"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/logger"
)

type BusRoutes struct {
	config   *config.Config
	db       *database.Client
	logger   logger.Logger
	busStore dataprovider.BusStore
	txer     dataprovider.Txer
}

func New(
	config *config.Config,
	db *database.Client,
	logger logger.Logger,
	busStore dataprovider.BusStore,
	txer dataprovider.Txer,
) *BusRoutes {
	return &BusRoutes{
		config:   config,
		db:       db,
		logger:   logger,
		busStore: busStore,
		txer:     txer,
	}
}

func (r *BusRoutes) IsHealthy(ctx context.Context) error {
	return r.db.StatusCheck(ctx)
}
