package server

import "github.com/gofiber/fiber/v3"

type Context interface {
	fiber.Ctx
	// JSONResponse is the function that used to send json response.
	JSONResponse(data interface{}, val error) error
}
