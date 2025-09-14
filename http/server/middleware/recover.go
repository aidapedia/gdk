package middleware

import (
	"github.com/aidapedia/gdk/log"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap/zapcore"
)

// WithRecover is the middleware that recovers from panics.
func WithRecover() Middleware {
	return func(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.ErrorCtx(c.Context(), "Panic recovered", zapcore.Field{
					Key:       "recover",
					Interface: r,
				})
			}
		}()
		return c.Next()
	}
}
