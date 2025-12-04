package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	gers "github.com/aidapedia/gdk/error"
	ghttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/log"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// WithRecover is the middleware that recovers from panics.
func WithRecover() fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				var ok bool
				var recErr error
				if recErr, ok = r.(error); !ok {
					recErr = fmt.Errorf("%v", r)
				}
				log.ErrorCtx(c.Context(), "Panic recovered", zap.Any("error", recErr), zap.ByteString("stack", debug.Stack()))
				// We ignore the error returned by JSONResponse and return nil to the fiber handler
				// because we have already written the response. Returning an error here would
				// trigger Fiber's default error handler, which might overwrite our response.
				_ = ghttp.JSONResponse(c, nil, gers.NewWithMetadata(recErr, ghttp.Metadata(http.StatusInternalServerError, "Internal Server Error")))
				err = nil
				return
			}
		}()
		return c.Next()
	}
}
