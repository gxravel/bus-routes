package api

import (
	"context"
	"encoding/json"
	"net/http"

	ierr "github.com/gxravel/bus-routes/internal/errors"
	"github.com/gxravel/bus-routes/internal/logger"
)

const (
	headerContentType   = "Content-Type"
	mimeApplicationJSON = "application/json"
)

type RangeItemsResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

// Response describes http response for api v1.
type Response struct {
	Data  interface{}    `json:"data,omitempty"`
	Error *ierr.APIError `json:"error,omitempty"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	if data == nil {
		w.WriteHeader(code)
		return
	}

	w.Header().Set(headerContentType, mimeApplicationJSON)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		logger.FromContext(ctx).WithErr(err).Error("encoding data to respond with json")
	}
}

func RespondEmpty(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func RespondCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// RespondDataOK responds with 200 status code and JSON in format: {"data": <val>}.
func RespondDataOK(ctx context.Context, w http.ResponseWriter, val interface{}) {
	RespondData(ctx, w, http.StatusOK, val)
}

// RespondData responds with custom status code and JSON in format: {"data": <val>}.
func RespondData(ctx context.Context, w http.ResponseWriter, code int, val interface{}) {
	RespondJSON(ctx, w, code, &Response{
		Data: val,
	})
}

// RespondError converts error to Reason, resolves http status code and responds with APIError.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) {
	reason := ierr.ConvertToReason(err)
	code := ierr.ResolveStatusCode(ierr.Cause(reason.Err))

	RespondJSON(ctx, w, code, &Response{
		Error: &ierr.APIError{
			Reason: reason,
		},
	})
}
