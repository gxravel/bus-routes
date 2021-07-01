package model

// Bus describes bus in bus_routes.bus.
type Bus struct {
	ID     int64  `db:"id"`
	CityID int    `db:"city_id"`
	Number string `db:"num"`

	// implicitly
	City string `db:"city"`
}
