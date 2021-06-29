package v1

type Bus struct {
	ID     int64  `json:"id,omitempty"`
	City   string `json:"city,omitempty"`
	CityID string `json:"city_id,omitempty"`
	Num    string `json:"num"`
}

type City struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

type Stop struct {
	ID      int64  `json:"id,omitempty"`
	City    string `json:"city,omitempty"`
	CityID  string `json:"city_id,omitempty"`
	Address string `json:"address"`
}

type Route struct {
	BusID  int64 `json:"bus_id"`
	StopID int64 `json:"stop_id"`
	Step   int8  `json:"step"`
}

type User struct {
	ID       int64  `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Type     string `json:"type,omitempty"`
}

type Token struct {
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}