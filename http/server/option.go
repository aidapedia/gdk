package server

import (
	"github.com/gofiber/fiber/v3"
)

type Option interface {
	Apply(svc *Server)
}

// WithMiddlewares is the option that adds middlewares to the server.
//
// Example:
//
//	WithMiddlewares(middleware.WithRecover())
func WithMiddlewares(middlewares ...fiber.Handler) Option {
	return &withMiddlewares{middlewares: middlewares}
}

type withMiddlewares struct {
	middlewares []fiber.Handler
}

func (o *withMiddlewares) Apply(svc *Server) {
	for _, m := range o.middlewares {
		svc.App.Use(m)
	}
}
