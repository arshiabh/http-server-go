package http

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/arshiabh/http-server-go/config"
)

type HTTPServer struct {
	router *Router
	logger *Logger
	config *config.Config
}

var logger *Logger

func NewHTTPServer(config *config.Config) *HTTPServer {
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

	// Accept connections in a loop
	for {
		// Accept a connection
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("Error accepting connection: %v", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go s.handleConnection(conn)
	}

}

func (s *HTTPServer) handleConnection(conn net.Conn) {
	// Ensure connection is always closed
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Error("Error closing connection: %v", err)
		}
	}()

	// Panic recovery
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("Panic in connection handler: %v", r)
			response := s.createErrorResponse(500, "Internal Server Error")
			if sendErr := s.sendResponse(conn, response); sendErr != nil {
				s.logger.Error("Failed to send error response after panic: %v", sendErr)
			}
		}
	}()

	// Get client address
	clientAddr := conn.RemoteAddr().String()
	s.logger.Info("New connection from: %s", clientAddr)

	// Set read timeout
	conn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))

	// Read data from the connection
	buffer := make([]byte, 8192)
	n, err := conn.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			s.logger.Error("Read timeout from %s", clientAddr)
			response := s.createErrorResponse(408, "Request Timeout")
			s.sendResponse(conn, response)
		} else {
			s.logger.Error("Error reading from %s: %v", clientAddr, err)
		}
		return
	}

	// Parse the HTTP request
	rawRequest := string(buffer[:n])
	s.logger.Debug("Raw request from %s: %s", clientAddr, strings.ReplaceAll(rawRequest, "\r\n", "\\r\\n"))

	request, err := parseHTTPRequest(rawRequest)
	if err != nil {
		s.logger.Error("Error parsing request from %s: %v", clientAddr, err)
		response := s.createErrorResponse(400, "Bad Request")
		if sendErr := s.sendResponse(conn, response); sendErr != nil {
			s.logger.Error("Failed to send bad request response: %v", sendErr)
		}
		return
	}

	// Log the parsed request
	s.logger.Info("Request: %s %s from %s", request.Method, request.Path, clientAddr)

	// Route the request and create response
	response := s.router.Route(request)

	// Send the response
	if err := s.sendResponse(conn, response); err != nil {
		s.logger.Error("Failed to send response to %s: %v", clientAddr, err)
		return
	}

	s.logger.Info("Response: %d %s to %s", response.StatusCode, response.StatusText, clientAddr)
}
