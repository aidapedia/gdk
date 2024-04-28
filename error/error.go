package error

import (
	"github.com/google/uuid"
)

// ErrorDetail type represent the error detail
type ErrorDetail string

// errorKit type represent the error message
type errorKit struct {
	ID         string      `json:"id,omitempty"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Detail     ErrorDetail `json:"details,omitempty"`
	RawRequest interface{} `json:"raw_request,omitempty"`
	RawError   error       `json:"error"`
}

// Error function used to return error message
func (e errorKit) Error() string {
	return e.RawError.Error()
}

// ErrorWithDetail function used to return error message with detail
func NewError(code int, err error, args ...interface{}) *errorKit {
	result := &errorKit{
		ID:       uuid.New().String(),
		Code:     code,
		RawError: err,
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

// WithRequest function used to add request to error
func WithRequest(err *errorKit, request interface{}) *errorKit {
	err.RawRequest = request
	return err
}

// CastError function used to convert error to errorKit
func CastError(err error) *errorKit {
	if err == nil {
		return nil
	}

	if e, ok := err.(*errorKit); ok {
		return e
	}

	return NewError(0, err, err.Error())
}
