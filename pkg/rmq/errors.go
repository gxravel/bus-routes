package rmq

import "net/http"

// Reason describes error reason.
type Reason struct {
	Err     error
	Message string
}

func (r *Reason) Error() string { return r.Err.Error() }

func (e *Reason) WithMessage(message string) *Reason {
	e.Message = message
	return e
}

func NewReason(err error) *Reason {
	return &Reason{
		Err: err,
	}
}

func ConvertToReason(err error) *Reason {
	switch val := err.(type) {
	case *Reason:
		return val

	default:
		return NewReason(err)
	}
}

type HTTPStatusCoder interface {
	HTTPStatusCode() int
}

func ResolveStatusCode(err error) int {
	var code int
	if val, ok := err.(HTTPStatusCoder); ok {
		code = val.HTTPStatusCode()
	}

	if code == 0 {
		code = http.StatusInternalServerError
	}

	return code
}
