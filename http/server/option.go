package server

import "github.com/aidapedia/gdk/http/server/middleware"

type Option interface {
	Apply(svc *Server)
}

func WithMiddlewares(middlewares ...middleware.Middleware) Option {
	return &withMiddlewares{middlewares: middlewares}
}

type withMiddlewares struct {
	middlewares []middleware.Middleware
}

func (o *withMiddlewares) Apply(svc *Server) {
	svc.middlewares = append(svc.middlewares, o.middlewares...)
}
