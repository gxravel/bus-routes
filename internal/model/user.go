package model

import "time"

// User describes user in bus_routes.user.
type User struct {
	ID             int64     `db:"id"`
	Email          string    `db:"email"`
	HashedPassword []byte    `db:"hashed_password"`
	Type           UserType  `db:"type"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type UserType string

func (t UserType) String() string { return string(t) }

const (
	UserAdmin       UserType = "admin"
	UserGuest       UserType = "guest"
	UserService     UserType = "service"
	DefaultUserType UserType = UserGuest
)

var (
	V1BusroutesUserTypes = []UserType{UserAdmin, UserGuest, UserService}
)

type UserTypes []UserType

// Exists returns true if types exist in t.
func (t UserTypes) Exists(types ...UserType) bool {
	uniqueTypes := make(map[UserType]struct{}, len(t))
	for _, ut := range t {
		uniqueTypes[ut] = struct{}{}
	}

	for _, typ := range types {
		if _, ok := uniqueTypes[typ]; ok {
			return true
		}
	}

	return false
}
