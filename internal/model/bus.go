package model

type Bus struct {
	ID     int64  `db:"id"`
	CityID int    `db:"city_id"`
	Num    string `db:"num"`
}
