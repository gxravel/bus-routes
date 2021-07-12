package v1

// Response describes amqp range items response for api v1.
type RangeItemsResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

// Response describes amqp response for api v1.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

// APIReason describes amqp model of error reason for api v1.
type APIReason struct {
	Err     string `json:"error"`
	Message string `json:"message,omitempty"`
}

// APIError describes amqp model of error for api v1.
type APIError struct {
	Code   int        `json:"code"`
	Reason *APIReason `json:"reason"`
}

// Response describes amqp response of amqp model for api v1.
type Bus struct {
	Num  string `json:"num"`
	City string `json:"city"`
}

// RoutePoint describes a unit of route for a bus for api v1.
type RoutePoint struct {
	Step    int8   `json:"step"`
	Address string `json:"address"`
}

// RouteDetailed describes amqp model of detailed route for api v1.
type RouteDetailed struct {
	City   string       `json:"city"`
	Bus    string       `json:"bus"`
	Points []RoutePoint `json:"points"`
}

// RangeBusesResponse describes response for range of routes for api v1.
type RangeRoutesResponse struct {
	Routes []*RouteDetailed `json:"items"`
	Total  int64            `json:"total"`
}
