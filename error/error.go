package error

import (
	"fmt"
	"runtime"
)

const (
	MetadataKeyCaller = "caller"
)

type Metadata map[string]interface{}

// Error is a struct to handle error
type Error struct {
	// raw error message
	message string
	// metadata can store anything that you need pass from error.
	// for example you want to different message from raw error and user error from backend.
	metadata Metadata
}

// New function used to create new error
func New(err error) *Error {
	ers, ok := err.(*Error)
	if !ok {
		error := &Error{
			metadata: make(Metadata),
			message:  err.Error(),
		}
		_, file, line, ok := runtime.Caller(1)
		if ok {
			error.metadata[MetadataKeyCaller] = fmt.Sprintf("%s:%d", file, line)
		}
		return error
	}
	if ers.Caller() == "" {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			ers.metadata[MetadataKeyCaller] = fmt.Sprintf("%s:%d", file, line)
		}
	}
	return ers
}

// NewWithMetadata function used to create new error with metadata
func NewWithMetadata(err error, metadata map[string]interface{}) *Error {
	ers, ok := err.(*Error)
	if !ok {
		error := &Error{
			metadata: metadata,
			message:  err.Error(),
		}
		_, file, line, ok := runtime.Caller(1)
		if ok {
			error.metadata[MetadataKeyCaller] = fmt.Sprintf("%s:%d", file, line)
		}
		return error
	}
	// Apply metadata
	for i := range metadata {
		ers.SetMetadata(i, metadata[i])
	}
	// Aplly Caller
	if ers.Caller() == "" {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			ers.metadata[MetadataKeyCaller] = fmt.Sprintf("%s:%d", file, line)
		}
	}
	return ers
}

// Error function used to return error message
func (e *Error) Error() string {
	return e.message
}

// Caller function used to return caller of the error
func (e *Error) Caller() string {
	return e.GetMetadataValue(MetadataKeyCaller).(string)
}

// SetMetadata function used to set metadata of the error
func (e *Error) SetMetadata(key string, value interface{}) {
	e.metadata[key] = value
}

// GetMetadata function used to get metadata of the error
func (e *Error) GetMetadataValue(key string) interface{} {
	if e.metadata[key] == nil {
		return nil
	}
	return e.metadata[key]
}

func (e *Error) GetMetadata() Metadata {
	return e.metadata
}
