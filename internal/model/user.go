package model

import "time"

type User struct {
	ID             int64     `db:"id"`
	Email          string    `db:"email"`
	HashedPassword []byte    `db:"hashed_password"`
	Type           string    `db:"type"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type UserType string

func (t UserType) String() string { return string(t) }

const (
	UserAdmin UserType = "admin"
	UserGuest UserType = "guest"
)

var (
	V1BusroutesUserTypes = []UserType{UserAdmin, UserGuest}
)
