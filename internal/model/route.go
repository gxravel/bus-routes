package model

type Route struct {
	BusID  int64 `db:"bus_id"`
	StopID int   `db:"stop_id"`
	Step   uint8 `db:"step"`
}
