package error

import (
	"fmt"
	"runtime"
)

// Error is a struct to handle error
type Error struct {
	// raw error message
	message string
	// caller of the error. this very helpful where the error come from.
	caller string
	// metadata can store anything that you need pass from error.
	// for example you want to different message from raw error and user error from backend.
	metadata map[string]interface{}
}

// New function used to create new error
func New() *Error {
	error := &Error{
		metadata: make(map[string]interface{}),
	}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		error.caller = fmt.Sprintf("%s:%d", file, line)
	}
	return error
}

// Error function used to return error message
func (e *Error) Error() string {
	return e.message
}

// Caller function used to return caller of the error
func (e *Error) Caller() string {
	return e.caller
}

// SetMetadata function used to set metadata of the error
func (e *Error) SetMetadata(key string, value interface{}) {
	e.metadata[key] = value
}

// GetMetadata function used to get metadata of the error
func (e *Error) GetMetadata(key string) interface{} {
	return e.metadata[key]
}
