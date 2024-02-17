package error

import (
	"github.com/google/uuid"
)

// ErrorDetail type represent the error detail
type ErrorDetail string

// errorKit type represent the error message
type errorKit struct {
	ID      string      `json:"id,omitempty"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  ErrorDetail `json:"details,omitempty"`
}

// Error function used to return error message
func (e errorKit) Error() string {
	return e.Message
}

// ErrorWithDetail function used to return error message with detail
func NewError(code int, args ...interface{}) *errorKit {
	result := &errorKit{
		ID:   uuid.New().String(),
		Code: code,
	}
	// apply the arguments
	for _, arg := range args {
		switch arg := arg.(type) {
		case ErrorDetail:
			result.Detail = arg
		case string:
			result.Message = arg
		}
	}

	return result
}
