package middleware

import (
	"github.com/aidapedia/gdk/context"
	"github.com/aidapedia/gdk/log"
	"github.com/gofiber/fiber/v3"
)

func WithContextLog() Middleware {
	return func(c fiber.Ctx) error {
		logID := c.Get(context.ContextKeyLogID)
		if logID == "" {
			logID = log.GenerateLogID()
			c.Set(context.ContextKeyLogID, logID)
		}
		c.Response().Header.Set(context.ContextKeyLogID, logID)
		c.Set(context.ContextKeyLogID, logID)
		return c.Next()
	}
}
