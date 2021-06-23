package model

type Bus struct {
	ID     int64  `db:"id"`
	City   string `db:"city"`
	CityID int    `db:"city_id"`
	Num    string `db:"num"`
}
