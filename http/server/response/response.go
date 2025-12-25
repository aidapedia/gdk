package response

import (
	"github.com/gofiber/fiber/v3"
)

// Use this struct to create a response from usecase level
type BaseResponse struct {
	Code    int
	Message string
}

type HTTPResponse struct {
	BaseResponse
	Data  interface{}
	Error error
}

// Use this function to send a JSON response from usecase level
func JSONResponse(c fiber.Ctx, rawResponse HTTPResponse) error {
	if rawResponse.Error != nil {
		resp := map[string]interface{}{
			"success": false,
			"message": "Internal Server Error",
		}
		statusCode := fiber.StatusInternalServerError
		if rawResponse.Code != 0 {
			statusCode = rawResponse.Code
		}
		if rawResponse.Message != "" {
			resp["message"] = rawResponse.Message
		}
		if rawResponse.Error != nil {
			resp["error"] = rawResponse.Error
		}
		c.Status(statusCode).JSON(resp)
		return rawResponse.Error
	}
	resp := map[string]interface{}{
		"success": true,
		"message": "Success",
	}
	statusCode := fiber.StatusOK
	if rawResponse.Code != 0 {
		statusCode = rawResponse.Code
	}
	if rawResponse.Message != "" {
		resp["message"] = rawResponse.Message
	}
	if rawResponse.Data != nil {
		resp["data"] = rawResponse.Data
	}
	c.Status(statusCode).JSON(resp)
	return nil
}
