package errors

import "net/http"

// APIError describes http model of error.
type APIError struct {
	Reason *Reason `json:"reason"`
}

type HTTPStatusCoder interface {
	HTTPStatusCode() int
}

type BadRequestError string

func (e BadRequestError) Type() ReasonType    { return ReasonProcessingError }
func (e BadRequestError) Error() string       { return string(e) }
func (e BadRequestError) HTTPStatusCode() int { return http.StatusBadRequest }

const (
	ErrBadRequest  BadRequestError = "bad request"
	ErrMustProvide BadRequestError = "must provide"
)

type ValidationError string

func (e ValidationError) Type() ReasonType    { return ReasonValidationError }
func (e ValidationError) Error() string       { return string(e) }
func (e ValidationError) HTTPStatusCode() int { return http.StatusBadRequest }

const (
	ErrValidationFailed ValidationError = "validation failed"
)

type AuthorizationError string

func (e AuthorizationError) Type() ReasonType    { return ReasonAuthError }
func (e AuthorizationError) Error() string       { return string(e) }
func (e AuthorizationError) HTTPStatusCode() int { return http.StatusUnauthorized }

const (
	ErrUnauthorized     AuthorizationError = "unauthorized"
	ErrWrongCredentials AuthorizationError = "wrong credentials"
	ErrInvalidToken     AuthorizationError = "invalid token"
	ErrInvalidJWT       AuthorizationError = "invalid JWT format"
	ErrTokenExpired     AuthorizationError = "token expired"
)

type ForbiddenError string

func (e ForbiddenError) Type() ReasonType    { return ReasonAuthError }
func (e ForbiddenError) Error() string       { return string(e) }
func (e ForbiddenError) HTTPStatusCode() int { return http.StatusForbidden }

const (
	ErrPermissionDenied ForbiddenError = "user does not have the permission"
)

type ConflictError string

func (e ConflictError) Type() ReasonType    { return ReasonProcessingError }
func (e ConflictError) Error() string       { return string(e) }
func (e ConflictError) HTTPStatusCode() int { return http.StatusConflict }

const (
	ErrConflict ConflictError = "conflict"
)

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
