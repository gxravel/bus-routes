package model

type Stop struct {
	ID      int64  `db:"id"`
	CityID  int    `db:"city_id"`
	Address string `db:"address"`
}
