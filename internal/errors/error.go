package errors

import (
	"strings"
)

type ReasonType string

const (
	ReasonUnknownError    ReasonType = "unknown_error"
	ReasonProcessingError ReasonType = "processing_error"
	ReasonAuthError       ReasonType = "authorization_error"
	ReasonValidationError ReasonType = "validation_error"
)

// Reason describes error reason.
type Reason struct {
	RType   ReasonType `json:"type"`
	Err     TypedError `json:"error"`
	Message string     `json:"message,omitempty"`
}

func (r *Reason) Type() ReasonType { return r.RType }

func (r *Reason) Error() string {
	cause := r.Cause()
	result := strings.Builder{}
	result.Grow(len(r.Message) + len(cause.Error()))
	if r.Message != "" {
		result.WriteString(r.Message + ": ")
	}
	result.WriteString(cause.Error())
	return result.String()
}

func (r *Reason) Cause() error { return r.Err }

func (e *Reason) WithMessage(message string) *Reason {
	e.Message = message
	return e
}

func NewReason(err TypedError) *Reason {
	return &Reason{
		Err:   err,
		RType: err.Type(),
	}
}

func ConvertToReason(err error) *Reason {
	switch val := err.(type) {
	case *Reason:
		return val
	case TypedError:
		return NewReason(val)
	default:
		typedError := newTypedError(ReasonUnknownError, val)
		return NewReason(typedError)
	}
}

type TypedError interface {
	error
	Type() ReasonType
}

type typedError struct {
	reasonType ReasonType
	err        error
	Message    string `json:"message"`
}

func (e *typedError) Type() ReasonType {
	return e.reasonType
}

func (e *typedError) Error() string {
	cause := e.Cause()
	if cause == nil {
		return ""
	}
	return cause.Error()
}

func (e *typedError) Cause() error {
	return e.err
}

func newTypedError(reasonType ReasonType, err error) TypedError {
	return &typedError{
		reasonType: reasonType,
		err:        err,
		Message:    err.Error(),
	}
}

// CheckDuplicate checks if the database error contains "Duplicate" and return updated error.
func CheckDuplicate(err error, field string) error {
	if strings.Contains(err.Error(), "Duplicate") {
		return NewReason(ErrConflict).WithMessage("the " + field + " is already in use")
	}
	return err
}

type causer interface {
	Cause() error
}

// Cause returns true cause error of err.
func Cause(err error) error {
	for err != nil {
		cause, okCause := err.(causer)
		if okCause {
			err = cause.Cause()
		} else {
			break
		}
	}

	return err
}
