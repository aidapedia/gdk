package error

import (
	"errors"

	goerrors "github.com/go-errors/errors"
)

// UserMessage type represent the error detail for user
type UserMessage string

type Code int

// Error type represent the error message
type Error struct {
	// code represent the error code, can be used to http error, grpc error, etc
	code Code
	// userMessage represent the error userMessage message that can be reader by user
	// this is used to give information about the error and what the user can do for the error
	userMessage UserMessage
	// param represent the request that cause the error
	// this helping debugging the error
	param map[string]interface{}
	// fork from github.com/go-errors/errors
	// this will be helpful to get the stack trace error
	err *goerrors.Error
}

// New function used to create new error message
func New(args ...interface{}) *Error {
	result := &Error{}
	// apply the arguments
	for _, arg := range args {
		switch arg := arg.(type) {
		case Code:
			result.code = arg
		case UserMessage:
			result.userMessage = arg
		case error:
			result.err = goerrors.New(arg)
		case string:
			result.err = goerrors.New(errors.New(arg))
		case map[string]interface{}:
			result.param = arg
		}
	}
	return result
}

// Error function used to return error message
func (e Error) Error() string {
	return e.err.Error()
}

// Error function used to return error message
func (e Error) IsError() bool {
	return e.err != nil
}

// Return the wrapped error (implements api for As function).
func (e *Error) GetError() *goerrors.Error {
	return e.err
}

// SetParam function used to add request to error
func (e *Error) SetParam(key string, request interface{}) {
	if e.param == nil {
		e.param = make(map[string]interface{})
	}
	e.param[key] = request
}

// GetParam function used to get request from error
func (e Error) GetParam() interface{} {
	return e.param
}

// SetCode function used to set error code
func (e *Error) SetCode(code Code) {
	e.code = code
}

// GetCode function used to get error code
func (e Error) GetCode() Code {
	return e.code
}

// SetUserMessage function used to set user message
func (e *Error) SetUserMessage(details UserMessage) {
	e.userMessage = details
}

// GetUserMessage function used to get user message
func (e Error) GetUserMessage() string {
	return string(e.userMessage)
}
