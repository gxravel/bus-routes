package busroutes

import (
	"github.com/gxravel/bus-routes/internal/config"
	"github.com/gxravel/bus-routes/internal/database"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/jwt"
	"github.com/gxravel/bus-routes/internal/logger"
)

type BusRoutes struct {
	config       *config.Config
	db           *database.Client
	logger       logger.Logger
	busStore     dataprovider.BusStore
	cityStore    dataprovider.CityStore
	stopStore    dataprovider.StopStore
	routeStore   dataprovider.RouteStore
	userStore    dataprovider.UserStore
	txer         dataprovider.Txer
	tokenManager jwt.Manager
}

func New(
	config *config.Config,
	db *database.Client,
	logger logger.Logger,
	busStore dataprovider.BusStore,
	cityStore dataprovider.CityStore,
	stopStore dataprovider.StopStore,
	routeStore dataprovider.RouteStore,
	userStore dataprovider.UserStore,
	txer dataprovider.Txer,
	jwtManager jwt.Manager,
) *BusRoutes {
	return &BusRoutes{
		config:       config,
		db:           db,
		logger:       logger,
		busStore:     busStore,
		cityStore:    cityStore,
		stopStore:    stopStore,
		routeStore:   routeStore,
		userStore:    userStore,
		txer:         txer,
		tokenManager: jwtManager,
	}
}
