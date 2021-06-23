package v1

type Stop struct {
	ID      int64  `json:"id,omitempty"`
	City    string `json:"city"`
	Address string `json:"address"`
}
