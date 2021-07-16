package busroutes

import (
	"github.com/gxravel/bus-routes/internal/config"
	"github.com/gxravel/bus-routes/internal/database"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/jwt"
	log "github.com/gxravel/bus-routes/internal/logger"
)

type Busroutes struct {
	config       *config.Config
	db           *database.Client
	logger       log.Logger
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
	logger log.Logger,
	busStore dataprovider.BusStore,
	cityStore dataprovider.CityStore,
	stopStore dataprovider.StopStore,
	routeStore dataprovider.RouteStore,
	userStore dataprovider.UserStore,
	txer dataprovider.Txer,
	jwtManager jwt.Manager,
) *Busroutes {
	return &Busroutes{
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
