package httpserver

import "net/http"

type RouteRegistrar interface {
	RegisterRoutes(*http.ServeMux)
}

type Middleware func(http.Handler) http.Handler

type Router struct {
	mux     *http.ServeMux
	handler http.Handler
}

func NewRouter(middlewares ...Middleware) *Router {
	mux := http.NewServeMux()

	return &Router{
		mux:     mux,
		handler: applyMiddlewares(mux, middlewares),
	}
}

func (r *Router) RegisterRoutes(registrars ...RouteRegistrar) {
	for _, registrar := range registrars {
		registrar.RegisterRoutes(r.mux)
	}
}

func (r *Router) RegisterRouteGroup(
	prefix string,
	middlewares []Middleware,
	registrars ...RouteRegistrar,
) {
	groupMux := http.NewServeMux()
	for _, registrar := range registrars {
		registrar.RegisterRoutes(groupMux)
	}

	r.mux.Handle(prefix, applyMiddlewares(groupMux, middlewares))
}

func (r *Router) Handler() http.Handler {
	return r.handler
}

func applyMiddlewares(handler http.Handler, middlewares []Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
