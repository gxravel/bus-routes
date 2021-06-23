package model

type Route struct {
	BusID  int64 `db:"bus_id"`
	StopID int64 `db:"stop_id"`
	Step   uint8 `db:"step"`
}
