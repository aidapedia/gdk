package http

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	gerr "github.com/aidapedia/gdk/error"
)

// We can define a base response struct that will be used by all other response structs.
// This will help us to maintain a consistent response format across the application.
// This will also help us to easily change the response format in the future.
// Example:
//
//	 type TestResponse struct {
//		*BaseResponse `json:"response,omitempty"`
//		Data interface{} `json:"data"`
//	 }
type BaseResponse struct {
	Status  int     `json:"status"`
	Message *string `json:"message,omitempty"`
	Error   string  `json:"error"`
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
func JSONResponse(c fiber.Ctx, data interface{}, val error) error {
	// Error Response Check
	err, ok := val.(*gerr.Error)
	if ok && err != nil {
		if ok {
			msg := err.GetMetadata(ErrorMetadataUserMessage)
			if msg == nil || msg == "" {
				msg = err.Error()
			}

			code := err.GetMetadata(ErrorMetadataCode)
			if code == nil {
				code = http.StatusInternalServerError
			}
			c.Status(code.(int)).JSON(&fiber.Map{
				"success": false,
				"message": msg,
			})
			return err
		}
		c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
		return err
	}
	// Success Response
	if data == nil {
		c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
			"message": "Success",
		})
		return err
	}
	c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "Success",
		"data":    data,
	})
	return nil
}

// HTTPMetadata is the function that will be used to create a metadata for HTTP response.
func Metadata(code int, message string) gerr.Metadata {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	return gerr.Metadata{
		ErrorMetadataCode:        code,
		ErrorMetadataUserMessage: message,
	}
}
