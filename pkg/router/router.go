package router

import (
	"net/http"
	"strings"
)

// Middleware represents a middleware function
type Middleware func(http.Handler) http.Handler

// Router is a simple HTTP router with middleware support
type Router struct {
	routes      map[string]http.Handler
	middlewares []Middleware
}

// New creates a new Router
func New() *Router {
	return &Router{
		routes:      make(map[string]http.Handler),
		middlewares: make([]Middleware, 0),
	}
}

// Use adds a middleware to the router
func (r *Router) Use(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

// Handle registers a handler for the given pattern
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.routes[pattern] = handler
}

// HandleFunc registers a handler function for the given pattern
func (r *Router) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	r.routes[pattern] = handlerFunc
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Find matching route
	handler := r.findHandler(req.URL.Path)

	if handler == nil {
		http.NotFound(w, req)
		return
	}

	// Apply middleware chain
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	handler.ServeHTTP(w, req)
}

// findHandler finds the appropriate handler for the given path
func (r *Router) findHandler(path string) http.Handler {
	// Exact match
	if handler, ok := r.routes[path]; ok {
		return handler
	}

	// Prefix match for wildcard routes
	for pattern, handler := range r.routes {
		if strings.HasSuffix(pattern, "/*") {
			prefix := strings.TrimSuffix(pattern, "/*")
			if strings.HasPrefix(path, prefix) {
				return handler
			}
		}
	}

	return nil
}
