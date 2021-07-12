package v1

import (
	"github.com/gxravel/bus-routes/internal/model"
)

// Response describes http range itmes response for api v1.
type RangeItemsResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

// Response describes http response for api v1.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

// APIReason describes http model of error reason for api v1.
type APIReason struct {
	RType   string `json:"type"`
	Err     string `json:"error"`
	Message string `json:"message,omitempty"`
}

// APIError describes http model of error for api v1.
type APIError struct {
	Reason *APIReason `json:"reason"`
}

// Bus describes http model of bus for api v1.
type Bus struct {
	ID     int64  `json:"id,omitempty"`
	CityID int    `json:"city_id,omitempty"`
	Num    string `json:"num"`

	City string `json:"city,omitempty"`
}

// City describes http model of city for api v1.
type City struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

// Stop describes http model of bus stop for api v1.
type Stop struct {
	ID      int64  `json:"id,omitempty"`
	CityID  int    `json:"city_id,omitempty"`
	Address string `json:"address"`

	City string `json:"city,omitempty"`
}

// Route describes http model of route for api v1.
type Route struct {
	BusID  int64 `json:"bus_id"`
	StopID int64 `json:"stop_id"`
	Step   int8  `json:"step"`
}

// RoutePoint describes a unit of route for a bus.
type RoutePoint struct {
	Step    int8   `json:"step"`
	Address string `json:"address"`
}

// RouteDetailed describes http model of detailed route for api v1.
type RouteDetailed struct {
	City   string       `json:"city"`
	Bus    string       `json:"bus"`
	Points []RoutePoint `json:"points"`
}

// User describes http model of user for api v1.
type User struct {
	ID       int64          `json:"id,omitempty"`
	Email    string         `json:"email"`
	Password string         `json:"password,omitempty"`
	Type     model.UserType `json:"type,omitempty"`
}

// Token describes http model of JWT token for api v1.
type Token struct {
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}
