package middleware

import (
	"github.com/aidapedia/gdk/util"
	"github.com/gofiber/fiber/v3"
)

// WithIPAllowList is the option that allowlist the IP.
// If the IP is in the allowlist, the request will be passed.
// Otherwise, the request will be dropped.
//
// Example:
//
//	WithIPAllowList([]string{"192.168.1.0/24"})
func WithIPAllowList(allowlist []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()
		for _, b := range allowlist {
			if util.CheckSubnet(ip, b) {
				return c.Next()
			}
		}
		return c.Drop()
	}
}
