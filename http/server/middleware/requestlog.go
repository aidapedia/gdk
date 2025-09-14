package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	gdkErr "github.com/aidapedia/gdk/error"
	"github.com/aidapedia/gdk/log"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WithRequestLog is the middleware that logs the request.
// If you want to masked your request. You can configure it by setting on log agent.
// Some agent have feature like redaction rule to mask the request body.
func WithRequestLog() Middleware {
	return func(c fiber.Ctx) error {
		// Preparing Request Data
		request := []zapcore.Field{
			zap.String("method", string(c.Request().Header.Method())),
			zap.String("uri", string(c.Request().RequestURI())),
		}
		reqBody := make(map[string]interface{})
		if c.Request().BodyStream() != http.NoBody { // Read
			switch string(c.Request().Header.Peek(fiber.HeaderContentType)) {
			case fiber.MIMEApplicationJSON:
				reqBodyByte := c.Request().Body()
				if reqBodyByte == nil {
					fmt.Println("Request Body is nil")
				}
				errx := json.Unmarshal(reqBodyByte, &reqBody)
				if errx != nil {
					log.ErrorCtx(c.Context(), "Unmarshal Error", zap.Error(errx))
				}
				c.Request().SetBody(reqBodyByte) // Reset
			case fiber.MIMEMultipartForm:
				form, err := c.Request().MultipartForm()
				if err != nil {
					log.ErrorCtx(c.Context(), "Parse Form Error", zap.Error(err))
				}
				for k, v := range form.Value {
					reqBody[k] = v
				}
				for k, v := range form.File {
					reqBody[k] = v
				}
			}
			request = append(request, zap.Any("request", reqBody))
		}
		// Call next handler and do logging
		err := c.Next()
		if err != nil {
			ers, ok := err.(*gdkErr.Error)
			if !ok {
				request = append(request, zap.Any("error", err.Error()))
			} else {
				request = append(request, zap.Any("error", ers.Error()))
				request = append(request, zap.Any("stacktrace", ers.Caller()))
			}
			log.ErrorCtx(c.Context(), "Incoming HTTP Request", request...)
			return nil
		}
		log.InfoCtx(c.Context(), "Incoming HTTP Request", request...)

		return nil
	}
}
