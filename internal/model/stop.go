package model

type Stop struct {
	ID      int64  `db:"id"`
	City    string `db:"city"`
	Address string `db:"address"`
}
