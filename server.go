package main

import (
	"fmt"
	"net"
)

type HTTPServer struct {
	router *Router
	logger *Logger
	config *Config
}

var logger *Logger

func NewHTTPServer(config *Config) *HTTPServer {
	logger = NewLogger(config.LogLevel)
	router := NewRouter(logger)

	return &HTTPServer{
		logger: logger,
		router: router,
		config: config,
	}
}

func (s *HTTPServer) Start() error {
	listener, err := net.Listen("tcp", s.config.Address())
	if err != nil {
		return fmt.Errorf("failed to start server %s", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			s.logger.Error("error in closing listener: %v", err)
		}
	}()

	s.logger.Info("HTTP Server starting on %s", s.config.Address())
	s.logger.Info("Configuration:")
	s.logger.Info("  Read Timeout: %v", s.config.ReadTimeout)
	s.logger.Info("  Write Timeout: %v", s.config.WriteTimeout)
	s.logger.Info("  Log Level: %s", s.config.LogLevel)

	s.setupRoutes()

	return nil
}

func (s *HTTPServer) setupRoutes() {
	s.logger.Info("Setting up routes...")

	// Register route handlers
	s.router.HandleFunc("GET", "/", s.handleHome)
	s.router.HandleFunc("GET", "/users", s.handleGetUsers)
	s.router.HandleFunc("POST", "/users", s.handleCreateUser)
	s.router.HandleFunc("GET", "/users/{id}", s.handleGetUser)
	s.router.HandleFunc("PUT", "/users/{id}", s.handleUpdateUser)
	s.router.HandleFunc("DELETE", "/users/{id}", s.handleDeleteUser)
	s.router.HandleFunc("GET", "/health", s.handleHealth)
	s.router.HandleFunc("GET", "/error", s.handleError)

	s.logger.Info("Routes registered successfully")

}
