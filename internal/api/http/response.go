package api

import (
	"context"
	"encoding/json"
	"net/http"

	httpv1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	ierr "github.com/gxravel/bus-routes/internal/errors"
	log "github.com/gxravel/bus-routes/internal/logger"
)

type MIME string

func (m MIME) String() string { return string(m) }

const (
	MIMEApplicationJSON MIME = "application/json"
)

const (
	HeaderContentType = "Content-Type"
)

func RespondJSON(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	if data == nil {
		w.WriteHeader(code)
		return
	}

	w.Header().Set(HeaderContentType, MIMEApplicationJSON.String())

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		log.
			FromContext(ctx).
			WithErr(err).
			Error("encode data to respond with json")
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

// RespondEmptyItems responds with empty items and 200 status code.
func RespondEmptyItems(ctx context.Context, w http.ResponseWriter) {
	RespondData(ctx, w, http.StatusOK, httpv1.RangeItemsResponse{})
}

// RespondData responds with custom status code and JSON in format: {"data": <val>}.
func RespondData(ctx context.Context, w http.ResponseWriter, code int, val interface{}) {
	RespondJSON(ctx, w, code, &httpv1.Response{
		Data: val,
	})
}

// RespondError converts error to Reason, resolves http status code and responds with APIError.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) {
	reason := ierr.ConvertToReason(err)
	code := ierr.ResolveStatusCode(reason.Err)

	RespondJSON(ctx, w, code, &httpv1.Response{
		Error: &httpv1.APIError{
			Reason: &httpv1.APIReason{
				RType:   string(reason.RType),
				Err:     reason.Error(),
				Message: reason.Message,
			},
		},
	})
}
