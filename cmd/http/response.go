package http

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type HTTPResponse struct {
	StatusCode int
	StatusText string
	Body       string
	Headers    map[string]string
}

func (s *HTTPServer) sendResponse(conn net.Conn, response *HTTPResponse) error {
	// Set write timeout
	conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))

	// Build the HTTP response string
	var responseStr strings.Builder

	// Status line
	responseStr.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", response.StatusCode, response.StatusText))

	// Headers
	for key, value := range response.Headers {
		responseStr.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Empty line to separate headers from body
	responseStr.WriteString("\r\n")

	// Body
	responseStr.WriteString(response.Body)

	// Send the response
	_, err := conn.Write([]byte(responseStr.String()))
	if err != nil {
		return fmt.Errorf("failed to write response: %v", err)
	}

	return nil
}

// helper
func (s *HTTPServer) createResponse(statusCode int, statusText, contentType, body string) *HTTPResponse {
	headers := map[string]string{
		"Content-Type":   contentType,
		"Content-Length": strconv.Itoa(len(body)),
		"Date":           time.Now().Format(time.RFC1123),
		"Server":         "CustomHTTPServer/1.0",
		"Connection":     "close",
	}

	return &HTTPResponse{
		StatusCode: statusCode,
		StatusText: statusText,
		Headers:    headers,
		Body:       body,
	}
}

func (s *HTTPServer) createJSONResponse(statusCode int, statusText string, data interface{}) *HTTPResponse {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		s.logger.Error("Failed to marshal JSON response: %v", err)
		return s.createErrorResponse(500, "Internal Server Error")
	}
	return s.createResponse(statusCode, statusText, "application/json", string(jsonBody))
}

func (s *HTTPServer) createErrorResponse(statusCode int, statusText string) *HTTPResponse {
	errorBody := map[string]interface{}{
		"error":     statusText,
		"code":      statusCode,
		"message":   getErrorMessage(statusCode),
		"timestamp": time.Now().Format(time.RFC3339),
	}
	return s.createJSONResponse(statusCode, statusText, errorBody)
}

func createErrorResponse(statusCode int, statusText string) *HTTPResponse {
	errorBody := map[string]interface{}{
		"error":     statusText,
		"code":      statusCode,
		"message":   getErrorMessage(statusCode),
		"timestamp": time.Now().Format(time.RFC3339),
	}
	jsonBody, _ := json.Marshal(errorBody)
	headers := map[string]string{
		"Content-Type":   "application/json",
		"Content-Length": strconv.Itoa(len(jsonBody)),
		"Date":           time.Now().Format(time.RFC1123),
		"Server":         "CustomHTTPServer/1.0",
	}
	return &HTTPResponse{
		StatusCode: statusCode,
		StatusText: statusText,
		Headers:    headers,
		Body:       string(jsonBody),
	}
}
