package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/stackus/errors"
)

type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func NewErrorResponse(err error) render.Renderer {
	if dErr, ok := err.(errors.Error); ok {
		return &ErrorResponse{
			StatusCode: dErr.HTTPCode(),
			Message:    fmt.Sprintf("%s: %s", errors.TypeCode(err), err.Error()),
			Err:        dErr,
		}
	}
	return &ErrorResponse{
		StatusCode: http.StatusBadRequest,
		Message:    err.Error(),
		Err:        err,
	}
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}
