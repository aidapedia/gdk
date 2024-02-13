package error

import (
	"net/http"

	"github.com/aidapedia/devkit/util"
	"github.com/google/uuid"
)

// ErrorKit type represent the error message
type ErrorKit struct {
	ID      string `json:"id,omitempty"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"details,omitempty"`
}

// Error function used to return error message
func (e ErrorKit) Error() string {
	return e.Detail
}

var (
	ErrorInternalServer *ErrorKit = &ErrorKit{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
		Detail:  "Internal Server Error.",
	}

	ErrorBadRequest *ErrorKit = &ErrorKit{
		Code:    http.StatusBadRequest,
		Message: "Bad Request",
		Detail:  "Bad Request.",
	}

	ErrorUnauthorized *ErrorKit = &ErrorKit{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
		Detail:  "Unauthorized.",
	}

	ErrorForbidden *ErrorKit = &ErrorKit{
		Code:    http.StatusForbidden,
		Message: "Forbidden",
		Detail:  "Forbidden.",
	}

	ErrorNotFound *ErrorKit = &ErrorKit{
		Code:    http.StatusNotFound,
		Message: "Not Found",
		Detail:  "Not Found.",
	}
)

// ErrorWithDetail function used to return error message with detail
func ErrorWithDetail(err *ErrorKit, detail string) *ErrorKit {
	result := &ErrorKit{
		ID:      err.ID,
		Code:    err.Code,
		Message: err.Message,
		Detail:  util.TernaryEqualString(err.Detail, "", detail),
	}
	if result.ID == "" {
		result.ID = uuid.New().String()
	}
	return result
}
