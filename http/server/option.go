package server

import "github.com/aidapedia/gdk/http/server/middleware"

type Option interface {
	Apply(svc *Server)
}

// WithMiddlewares is the option that adds middlewares to the server.
//
// Example:
//
//	WithMiddlewares(middleware.WithRecover())
func WithMiddlewares(middlewares ...middleware.Middleware) Option {
	return &withMiddlewares{middlewares: middlewares}
}

type withMiddlewares struct {
	middlewares []middleware.Middleware
}

func (o *withMiddlewares) Apply(svc *Server) {
	for _, m := range o.middlewares {
		svc.app.Use(m)
	}
}
