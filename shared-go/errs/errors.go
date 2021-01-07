package errs

// Domain error struct
type Error struct {
	Code    int
	Message string
	Err     error
}

/*

Error{} can be tested using errors.Is() and unwrapped with errors.As()

Update and change as needed to meet your domain needs. Codes are HTTP status codes initially for
their well known status.

Start domain specific errors using the domain.NewError() constructor

Example:
  // creates a thing domain error that wraps ErrNotFound and reuses the code but sets a new message
  ErrThingNotFound = domain.NewError("thing not found", domain.ErrNotFound)

*/

// Base domain errors
var (
	ErrClientError         = Error{400, "unknown client error", nil}
	ErrBadRequest          = Error{400, "bad request", ErrClientError}
	ErrUnauthorized        = Error{401, "unauthorized", ErrClientError}
	ErrForbidden           = Error{403, "forbidden", ErrClientError}
	ErrNotFound            = Error{404, "not found", ErrClientError}
	ErrNotAcceptable       = Error{406, "not acceptable", ErrClientError}
	ErrRequestTimeout      = Error{408, "request timeout", ErrClientError}
	ErrConflict            = Error{409, "conflict", ErrClientError}
	ErrGone                = Error{410, "gone", ErrClientError}
	ErrUnprocessableEntity = Error{422, "unprocessable entity", ErrClientError}
	ErrTooManyRequests     = Error{429, "too many requests", ErrClientError}

	ErrServerError    = Error{500, "unknown server error", nil}
	ErrNotImplemented = Error{501, "not implemented", ErrServerError}
)

// Build a new error with optional code, message, or wrapped error
func NewError(args ...interface{}) Error {
	err := Error{}

	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			err.Code = v
		case string:
			err.Message = v
		case Error:
			if err.Code == 0 {
				err.Code = v.Code
			}
			if err.Message == "" {
				err.Message = v.Message
			}
			err.Err = v
		case error:
			err.Err = v
		}
	}

	return err
}

// implement error interface
func (e Error) Error() string {
	return e.Message
}

// Support unwrapping errors : errors.Is() & errors.As()
func (e Error) Unwrap() error {
	return e.Err
}
