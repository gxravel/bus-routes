package model

// Route describes route in bus_routes.route.
type Route struct {
	BusID  int64 `db:"bus_id"`
	StopID int64 `db:"stop_id"`
	Step   int8  `db:"step"`
}
