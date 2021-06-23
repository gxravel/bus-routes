package model

type Stop struct {
	ID      int64  `db:"id"`
	City    string `db:"city"`
	CityID  string `db:"city_id"`
	Address string `db:"address"`
}
