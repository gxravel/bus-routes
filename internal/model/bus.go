package model

type Bus struct {
	ID   int64  `db:"id"`
	City string `db:"city"`
	Num  string `db:"num"`
}
