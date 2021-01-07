package web

import (
	"net/http"

	"github.com/go-chi/render"

	"shared-go/errs"
)

type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func NewErrorResponse(err error) render.Renderer {
	if dErr, ok := err.(errs.Error); ok {
		return &ErrorResponse{
			StatusCode: dErr.Code,
			Message:    dErr.Message,
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
