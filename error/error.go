package error

import (
	"github.com/go-errors/errors"
)

// ErrorDetail type represent the error detail
type ErrorDetail string

// Error type represent the error message
type Error struct {
	code    int
	details ErrorDetail
	param   interface{}
	// Fork from go-errors/errors
	err *errors.Error
}

// Error function used to return error message
func (e Error) Error() string {
	return e.err.Error()
}

// New function used to create new error message
func New(code int, args ...interface{}) *Error {
	result := &Error{
		code: code,
	}
	// apply the arguments
	for _, arg := range args {
		switch arg := arg.(type) {
		case ErrorDetail:
			result.details = arg
		case error:
		case string:
			result.err = errors.New(arg)
		}
	}
	return result
}

// Return the wrapped error (implements api for As function).
func (e *Error) GetError() *errors.Error {
	return e.err
}

// SetParam function used to add request to error
func (e *Error) SetParam(request interface{}) {
	e.param = request
}

// GetParam function used to get request from error
func (e Error) GetParam() interface{} {
	return e.param
}

// SetCode function used to set error code
func (e *Error) SetCode(code int) {
	e.code = code
}

// GetCode function used to get error code
func (e Error) GetCode() int {
	return e.code
}

// SetDetails function used to set error details
func (e *Error) SetDetails(details ErrorDetail) {
	e.details = details
}

// GetDetails function used to get error details
func (e Error) GetDetails() string {
	return string(e.details)
}
