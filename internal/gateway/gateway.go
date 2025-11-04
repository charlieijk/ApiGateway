package gateway

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"apigateway/internal/config"
	"apigateway/internal/middleware"
	"apigateway/pkg/router"
)

// Gateway represents the API Gateway
type Gateway struct {
	config *config.Config
	router *router.Router
}

// New creates a new API Gateway instance
func New(cfg *config.Config) (*Gateway, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	gw := &Gateway{
		config: cfg,
		router: router.New(),
	}

	// Setup routes
	if err := gw.setupRoutes(); err != nil {
		return nil, fmt.Errorf("failed to setup routes: %w", err)
	}

	return gw, nil
}

// Router returns the HTTP handler
func (gw *Gateway) Router() http.Handler {
	return gw.router
}

// setupRoutes configures all the routes and their proxies
func (gw *Gateway) setupRoutes() error {
	// Apply global middleware
	gw.router.Use(middleware.Logger)
	gw.router.Use(middleware.CORS)
	gw.router.Use(middleware.RateLimiter(gw.config.RateLimit))

	// Health check endpoint
	gw.router.HandleFunc("/health", gw.healthCheckHandler)

	// Setup service routes
	for _, service := range gw.config.Services {
		if err := gw.addServiceRoute(service); err != nil {
			return fmt.Errorf("failed to add route for service %s: %w", service.Name, err)
		}
	}

	return nil
}

// addServiceRoute adds a route for a backend service
func (gw *Gateway) addServiceRoute(service config.Service) error {
	targetURL, err := url.Parse(service.Target)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Customize proxy behavior
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, fmt.Sprintf("Service unavailable: %v", err), http.StatusServiceUnavailable)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = targetURL.Host
		proxy.ServeHTTP(w, r)
	})

	// Register the route
	gw.router.Handle(service.Path, handler)

	return nil
}

// healthCheckHandler handles health check requests
func (gw *Gateway) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"api-gateway"}`))
}
