package httpserver

import (
	"context"
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/config"
)

type Server struct {
	server *http.Server
	router *Router
}

func NewServer(cfg config.HTTPConfig, middlewares ...Middleware) *Server {
	router := NewRouter(middlewares...)

	return &Server{
		server: &http.Server{
			Addr:              cfg.HTTPAddress(),
			Handler:           router.Handler(),
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
		},
		router: router,
	}
}

func (s *Server) RegisterRoutes(registrars ...RouteRegistrar) {
	s.router.RegisterRoutes(registrars...)
}

func (s *Server) RegisterRouteGroup(
	prefix string,
	middlewares []Middleware,
	registrars ...RouteRegistrar,
) {
	s.router.RegisterRouteGroup(prefix, middlewares, registrars...)
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Close() error {
	return s.server.Close()
}

func (s *Server) Address() string {
	return s.server.Addr
}
