package middleware

import (
	"context"

	gctx "github.com/aidapedia/gdk/context"
	"github.com/aidapedia/gdk/log"
	"github.com/gofiber/fiber/v3"
)

func WithContextLog() fiber.Handler {
	return func(c fiber.Ctx) error {
		logID := string(c.Request().Header.Peek(gctx.ContextKeyLogID))
		if logID == "" {
			logID = log.GenerateLogID()
		}
		c.Response().Header.Set(gctx.ContextKeyLogID, logID)
		c.SetContext(context.WithValue(context.Background(), gctx.ContextKeyLogID, logID))
		return c.Next()
	}
}
