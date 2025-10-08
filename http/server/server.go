package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aidapedia/gdk/http/server/middleware"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
)

// Server is a struct to handle server
type Server struct {
	app *fiber.App
}

// New creates a new server
func New(opt ...Option) (*Server, error) {
	svc := &Server{
		app: fiber.New(
			fiber.Config{
				JSONEncoder: sonic.Marshal,
				JSONDecoder: sonic.Unmarshal,
			},
		),
	}
	for _, o := range opt {
		o.Apply(svc)
	}
	return svc, nil
}

// NewWithDefaultConfig creates a new server with default config
// This config choosen by the author of the package
func NewWithDefaultConfig(opt ...Option) (*Server, error) {
	return New(
		WithMiddlewares(
			middleware.WithContextLog(),
			middleware.WithRecover(),
			middleware.WithRequestLog(),
		),
	)
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

	return s.app.Listen(address, config...)
}

// Shutdown shuts down the server
func (s *Server) Shutdown() {
	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down triggered")
		s.app.Shutdown()
	}()
}

// ListenGracefully starts the server with the given address and config
// It will shutdown the server gracefully when os.Interrupt, syscall.SIGTERM, or syscall.SIGQUIT signal is received
// It will return error if the server failed to start
func (s *Server) ListenGracefully(address string, config ...fiber.ListenConfig) error {
	s.Shutdown()
	return s.Listen(address, config...)
}

func (s *Server) GetFiberApp() *fiber.App {
	return s.app
}
