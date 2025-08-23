package server

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
)

// Server is a struct to handle server
type Server struct {
	*fiber.App
}

// New creates a new server
func New() *Server {
	return &Server{
		App: fiber.New(
			fiber.Config{
				JSONEncoder: sonic.Marshal,
				JSONDecoder: sonic.Unmarshal,
			},
		),
	}
}

// Listen starts the server with the given address and config
func (s *Server) Listen(address string, config ...fiber.ListenConfig) error {
	if len(config) > 0 {
		config[0].DisableStartupMessage = true
		config[0].EnablePrefork = true
	} else {
		config = append(config, fiber.ListenConfig{
			DisableStartupMessage: true,
			EnablePrefork:         true,
		})
	}
	return s.App.Listen(address, config...)
}
