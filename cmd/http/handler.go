package http

import (
	"encoding/json"
	"strings"
	"time"
)

func (s *HTTPServer) handleHome(request *HTTPRequest) *HTTPResponse {
	body := "Welcome to our Custom HTTP Server!\n\nAvailable endpoints:\n- GET /\n- GET /users\n- POST /users\n- GET /users/{id}\n- PUT /users/{id}\n- DELETE /users/{id}\n- GET /health\n- GET /error"
	return s.createResponse(200, "OK", "text/plain", body)
}

func (s *HTTPServer) handleGetUsers(request *HTTPRequest) *HTTPResponse {
	users := []map[string]interface{}{
		{"id": 1, "name": "John", "email": "john@example.com"},
		{"id": 2, "name": "Jane", "email": "jane@example.com"},
	}
	return s.createJSONResponse(200, "OK", users)
}

func (s *HTTPServer) handleCreateUser(request *HTTPRequest) *HTTPResponse {
	if request.Body == "" {
		return s.createErrorResponse(400, "Bad Request - Body required")
	}

	var userData map[string]interface{}
	if err := json.Unmarshal([]byte(request.Body), &userData); err != nil {
		s.logger.Error("Failed to parse JSON body: %v", err)
		return s.createErrorResponse(400, "Bad Request - Invalid JSON")
	}

	response := map[string]interface{}{
		"message": "User created successfully",
		"data":    userData,
	}
	return s.createJSONResponse(201, "Created", response)
}

func (s *HTTPServer) handleGetUser(request *HTTPRequest) *HTTPResponse {
	userID := strings.TrimPrefix(request.Path, "/users/")
	if userID == "" {
		return s.createErrorResponse(400, "Bad Request - Invalid user ID")
	}

	user := map[string]interface{}{
		"id":    userID,
		"name":  "User " + userID,
		"email": "user" + userID + "@example.com",
	}
	return s.createJSONResponse(200, "OK", user)
}

func (s *HTTPServer) handleUpdateUser(request *HTTPRequest) *HTTPResponse {
	userID := strings.TrimPrefix(request.Path, "/users/")
	if userID == "" {
		return s.createErrorResponse(400, "Bad Request - Invalid user ID")
	}

	if request.Body == "" {
		return s.createErrorResponse(400, "Bad Request - Body required")
	}

	var userData map[string]interface{}
	if err := json.Unmarshal([]byte(request.Body), &userData); err != nil {
		s.logger.Error("Failed to parse JSON body: %v", err)
		return s.createErrorResponse(400, "Bad Request - Invalid JSON")
	}

	response := map[string]interface{}{
		"message": "User " + userID + " updated successfully",
		"data":    userData,
	}
	return s.createJSONResponse(200, "OK", response)
}

func (s *HTTPServer) handleDeleteUser(request *HTTPRequest) *HTTPResponse {
	userID := strings.TrimPrefix(request.Path, "/users/")
	if userID == "" {
		return s.createErrorResponse(400, "Bad Request - Invalid user ID")
	}

	response := map[string]interface{}{
		"message": "User " + userID + " deleted successfully",
	}
	return s.createJSONResponse(200, "OK", response)
}

func (s *HTTPServer) handleHealth(request *HTTPRequest) *HTTPResponse {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"server":    "CustomHTTPServer/1.0",
		"uptime":    "unknown",
	}
	return s.createJSONResponse(200, "OK", health)
}

func (s *HTTPServer) handleError(request *HTTPRequest) *HTTPResponse {
	if strings.Contains(request.Path, "panic") {
		panic("This is a test panic!")
	}
	return s.createErrorResponse(500, "Internal Server Error")
}

// helper

func getErrorMessage(statusCode int) string {
	switch statusCode {
	case 400:
		return "The request was invalid"
	case 404:
		return "The requested resource was not found"
	case 405:
		return "The HTTP method is not allowed for this resource"
	case 408:
		return "The request took too long to complete"
	case 500:
		return "Internal server error"
	default:
		return "An error occurred"
	}
}
