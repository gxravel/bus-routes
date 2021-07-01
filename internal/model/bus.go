package model

// Bus describes bus in bus_routes.bus.
type Bus struct {
	ID     int64  `db:"id"`
	City   string `db:"city"`
	CityID int    `db:"city_id"`
	Num    string `db:"num"`
}
