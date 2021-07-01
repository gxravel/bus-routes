package model

// Stop describes stop in bus_routes.stop.
type Stop struct {
	ID      int64  `db:"id"`
	City    string `db:"city"`
	CityID  int    `db:"city_id"`
	Address string `db:"address"`
}
