package model

// City describes city in bus_routes.city.
type City struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
