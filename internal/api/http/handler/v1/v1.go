package v1

import "github.com/gxravel/bus-routes/internal/model"

// Bus describes http model of bus for api v1.
type Bus struct {
	ID     int64  `json:"id,omitempty"`
	City   string `json:"city,omitempty"`
	CityID string `json:"city_id,omitempty"`
	Num    string `json:"num"`
}

// City describes http model of city for api v1.
type City struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

// Stop describes http model of bus stop for api v1.
type Stop struct {
	ID      int64  `json:"id,omitempty"`
	City    string `json:"city,omitempty"`
	CityID  string `json:"city_id,omitempty"`
	Address string `json:"address"`
}

// Route describes http model of route for api v1.
type Route struct {
	BusID  int64 `json:"bus_id"`
	StopID int64 `json:"stop_id"`
	Step   int8  `json:"step"`
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
