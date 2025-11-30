package middleware

import (
	"net/http"

	gdkErr "github.com/aidapedia/gdk/error"
	gdkHttp "github.com/aidapedia/gdk/http"
	"github.com/aidapedia/gdk/log"
	"github.com/aidapedia/gdk/mask"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WithRequestLog is the middleware that logs the request.
// If you want to masked your request. You can configure it by setting on log agent.
// Some agent have feature like redaction rule to mask the request body.
func WithRequestLog(mask *mask.Mask) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Preparing Request Data
		request := []zapcore.Field{
			zap.String("method", string(c.Request().Header.Method())),
			zap.String("uri", string(c.Request().RequestURI())),
		}
		reqBody := make(map[string]interface{})
		if c.Request().BodyStream() != http.NoBody {
			switch string(c.Request().Header.Peek(fiber.HeaderContentType)) {
			case fiber.MIMEApplicationJSON:
				reqBodyByte := c.Request().Body()
				if reqBodyByte != nil && len(reqBodyByte) > 0 {
					errx := sonic.Unmarshal(reqBodyByte, &reqBody)
					if errx != nil {
						log.ErrorCtx(c.Context(), "Unmarshal Error", zap.Error(errx))
					}
				}
			case fiber.MIMEMultipartForm, fiber.MIMEApplicationForm:
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
			case fiber.MIMETextXML, fiber.MIMEApplicationXML, fiber.MIMEApplicationCBOR:
				reqBody["body"] = string(c.Request().Body())
			}
			// Masking Request Body
			if mask != nil {
				reqBody, _ = mask.MaskMap(reqBody)
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
				request = append(request, zap.Any("error", err.Error()))
				meta := ers.GetMetadata()
				for k, v := range meta {
					if k == gdkHttp.ErrorMetadataUserMessage || k == gdkHttp.ErrorMetadataCode {
						continue
					}
					request = append(request, zap.Any(k, v))
				}
			}
			log.ErrorCtx(c.Context(), "Incoming HTTP Request", request...)
			return nil
		}
		log.InfoCtx(c.Context(), "Incoming HTTP Request", request...)
		return nil
	}
}
