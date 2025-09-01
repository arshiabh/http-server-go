package main

type RouteHandler func(*HTTPRequest) *HTTPResponse

type Route struct {
	Method  string
	Path    string
	Handler RouteHandler
}

type Router struct {
	Routes []Route
	Logger *Logger
}

func NewRouter(logger *Logger) *Router {
	return &Router{
		Routes: make([]Route, 0),
		Logger: logger,
	}
}

func (s *Router) HandleFunc(method string, path string, handler RouteHandler) {
	
}
