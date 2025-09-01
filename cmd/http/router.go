package http

import "strings"

type RouteHandler func(*HTTPRequest) *HTTPResponse

type Route struct {
	Method  string
	Path    string
	Handler RouteHandler
}

type Router struct {
	routes []Route
	logger *Logger
}

func NewRouter(logger *Logger) *Router {
	return &Router{
		routes: make([]Route, 0),
		logger: logger,
	}
}

func (r *Router) HandleFunc(method, path string, handler RouteHandler) {
	route := Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
	r.routes = append(r.routes, route)
	r.logger.Debug("Registered route: %s %s", method, path)
}

// Route finds and executes the appropriate handler
func (r *Router) Route(request *HTTPRequest) *HTTPResponse {
	defer func() {
		if rec := recover(); rec != nil {
			r.logger.Error("Panic in route handler for %s %s: %v", request.Method, request.Path, rec)
		}
	}()

	// Find matching route
	for _, route := range r.routes {
		if r.matchRoute(route, request) {
			return route.Handler(request)
		}
	}

	// No route found
	r.logger.Info("Route not found: %s %s", request.Method, request.Path)
	return createErrorResponse(404, "Not Found")
}

// matchRoute checks if a route matches the request
func (r *Router) matchRoute(route Route, request *HTTPRequest) bool {
	// Check method
	if route.Method != request.Method {
		return false
	}

	// Exact path match
	if route.Path == request.Path {
		return true
	}

	// Pattern matching for paths like /users/{id}
	if strings.Contains(route.Path, "{id}") {
		pattern := strings.Replace(route.Path, "{id}", "", 1)
		return strings.HasPrefix(request.Path, pattern) && len(request.Path) > len(pattern)
	}

	return false
}
