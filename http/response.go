package http

import (
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
// It will check if the error is a gerr.Error.
// If it is, it will send a JSON response with the error message.
// If it is not, it will send a JSON response with the data.
func JSONResponse(c fiber.Ctx, data interface{}, val error) error {
	err, ok := val.(*gerr.Error)
	if ok && err != nil {
		msg := err.GetMetadata(ErrorMetadataUserMessage).(string)
		c.Status(err.GetMetadata(ErrorMetadataCode).(int)).JSON(&fiber.Map{
			"success": false,
			"message": msg,
		})
		return err
	}

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
