package server

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

var (
	ErrAppNil = errors.New("app is nil")
)

type Option interface {
	Apply(svc *Server) error
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

func (o *withMiddlewares) Apply(svc *Server) error {
	if svc.App == nil {
		return ErrAppNil
	}
	for _, m := range o.middlewares {
		svc.App.Use(m)
	}
	return nil
}

func WithPostShutdown(hook ...fiber.OnPostShutdownHandler) Option {
	return &withPostShutdown{hook: hook}
}

type withPostShutdown struct {
	hook []fiber.OnPostShutdownHandler
}

func (o *withPostShutdown) Apply(svc *Server) error {
	if svc.App == nil {
		return ErrAppNil
	}
	svc.App.Hooks().OnPostShutdown(o.hook...)
	return nil
}

func WithAppConfig(config fiber.Config) Option {
	return &withAppConfig{config: config}
}

type withAppConfig struct {
	config fiber.Config
}

func (o *withAppConfig) Apply(svc *Server) error {
	svc.config = &o.config
	return nil
}
