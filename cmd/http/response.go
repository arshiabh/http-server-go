package http

import (
	"fmt"
	"net"
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
