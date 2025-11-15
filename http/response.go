package http

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	gerr "github.com/aidapedia/gdk/error"
)

type SuccessResponse struct {
	StatusCode int
	Message    string
	Data       interface{}
}

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
	// Error Response Check
	if valError != nil {
		err, ok := valError.(*gerr.Error)
		if ok {
			msg := err.GetMetadataValue(ErrorMetadataUserMessage)
			if msg == nil || msg == "" {
				msg = err.Error()
			}

			code := err.GetMetadataValue(ErrorMetadataCode)
			if code == nil {
				code = http.StatusInternalServerError
			}
			c.Status(code.(int)).JSON(&fiber.Map{
				"success": false,
				"message": msg,
			})
			return err
		} else {
			c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"message": valError.Error(),
			})
			return valError
		}
	}
	// Success Response
	successResponse := map[string]interface{}{
		"success": true,
		"message": "Success",
	}
	statusCode := fiber.StatusOK
	if valSuccess != nil {
		if valSuccess.StatusCode != 0 {
			statusCode = valSuccess.StatusCode
		}
		if valSuccess.Message == "" {
			successResponse["message"] = "Success"
		}
		if valSuccess.Data != nil {
			successResponse["data"] = valSuccess.Data
		}
	}
	c.Status(statusCode).JSON(successResponse)
	return nil
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
