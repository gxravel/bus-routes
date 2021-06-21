package model

type City struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
