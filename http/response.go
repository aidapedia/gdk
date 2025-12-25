package http

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	gerr "github.com/aidapedia/gdk/error"
	"github.com/aidapedia/gdk/http/server/response"
)

type SuccessResponse struct {
	StatusCode int
	Message    string
	Data       interface{}
}

// DEPRECATED: Use response.JSONResponse instead
// JSONResponse is the function that will be used to send a JSON response.
// For custome message, you can use gerr.Metadata to store the message.
// Raw error will show on your log message.
// Example:
//
//	 err := gerr.New("Connection Lost")
//		return gerr.NewWithMetadata(err, http.Metadata(http.StatusBadRequest, "Internal Server Error"))
//
// JSON Response:
//
//	 {
//		"success": false,
//		"message": "Internal Server Error",
//	 }
//
// But if you want to show the raw error message to the user, you can leave ErrorMetadataUserMessage empty.
// Example:
//
//	 err := gerr.New("Connection Lost")
//			return gerr.NewWithMetadata(err, http.Metadata(http.StatusBadRequest, ""))
//
// JSON Response:
//
//	 {
//		"success": false,
//		"message": "Connection Lost",
//	 }
func JSONResponse(c fiber.Ctx, valSuccess *SuccessResponse, valError error) error {
	var (
		statusCode int
		message    string
	)

	if valError != nil {
		statusCode = http.StatusInternalServerError
		message = "Internal Server Error"

		err, ok := valError.(*gerr.Error)
		if ok {
			msg := err.GetMetadataValue(ErrorMetadataUserMessage)
			if msg == nil || msg == "" {
				message = err.Error()
			}

			code := err.GetMetadataValue(ErrorMetadataCode)
			if code == nil {
				statusCode = http.StatusInternalServerError
			}
		} else {
			message = valError.Error()
		}
		return response.JSONResponse(c, response.HTTPResponse{
			BaseResponse: response.BaseResponse{
				Code:    statusCode,
				Message: message,
			},
			Error: valError,
		})
	}

	if valSuccess != nil {
		statusCode = valSuccess.StatusCode
		message = valSuccess.Message
	} else {
		statusCode = http.StatusOK
		message = "Success"
	}

	return response.JSONResponse(c, response.HTTPResponse{
		BaseResponse: response.BaseResponse{
			Code:    statusCode,
			Message: message,
		},
		Data: valSuccess.Data,
	})
}

// HTTPMetadata is the function that will be used to create a metadata for HTTP response.
//
// keyPairs must be in the format of key, value, key, value, ...
func Metadata(code int, message string, keyPairs ...interface{}) gerr.Metadata {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	metadata := gerr.Metadata{
		ErrorMetadataCode:        code,
		ErrorMetadataUserMessage: message,
	}
	if len(keyPairs) > 0 && len(keyPairs)%2 != 0 {
		return metadata
	}
	for i := 0; i < len(keyPairs); i += 2 {
		metadata[keyPairs[i].(string)] = keyPairs[i+1]
	}
	return metadata
}
