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
	*fiber.App
	config *fiber.Config
}

// New creates a new server
// Basically, it creates a new fiber app with the given config
func New(serverName string, opt ...Option) (*Server, error) {
	svc := &Server{}
	optPostInit := []Option{}
	for _, o := range opt {
		err := o.Apply(svc)
		if err == ErrAppNil {
			optPostInit = append(optPostInit, o)
		}
	}
	if svc.config == nil {
		svc.config = &fiber.Config{
			JSONEncoder:   sonic.Marshal,
			JSONDecoder:   sonic.Unmarshal,
			ServerHeader:  serverName,
			StrictRouting: true,
			CaseSensitive: true,
			Immutable:     true,
		}
	} else {
		svc.config.ServerHeader = serverName
	}
	svc.App = fiber.New(*svc.config)
	for _, o := range optPostInit {
		o.Apply(svc)
	}
	return svc, nil
}

// NewWithDefaultConfig creates a new server with default config
// This config choosen by the author of the package
func NewWithDefaultConfig(serverName string, opt ...Option) (*Server, error) {
	opt = append(opt,
		WithMiddlewares(
			middleware.WithContextLog(),
			middleware.WithRecover(),
			middleware.WithRequestLog(),
		),
	)
	return New(serverName, opt...)
}

// shutdown shuts down the server
func (s *Server) shutdown() {
	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down triggered")
		s.App.Shutdown()
	}()
}

// Listen starts the server with the given address and config
// It will shutdown the server gracefully when os.Interrupt, syscall.SIGTERM, or syscall.SIGQUIT signal is received
// It will return error if the server failed to start
func (s *Server) Listen(address string, config ...fiber.ListenConfig) error {
	s.shutdown()
	// Handle path not found
	s.App.Use(func(c fiber.Ctx) error {
		return c.SendStatus(404)
	})
	return s.App.Listen(address, config...)
}
