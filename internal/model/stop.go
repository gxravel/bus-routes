package model

// Stop describes stop in bus_routes.stop.
type Stop struct {
	ID      int64  `db:"id"`
	CityID  int    `db:"city_id"`
	Address string `db:"address"`

	// implicitly
	City string `db:"city"`
}
