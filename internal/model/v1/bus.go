package v1

type Bus struct {
	ID   int64  `json:"id,omitempty"`
	City string `json:"city"`
	Num  string `json:"num"`
}
