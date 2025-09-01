package http

import (
	"fmt"
	"strconv"
	"strings"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Raw     string
	Headers map[string]string
	Body    string
}

func parseHTTPRequest(rawRequest string) (*HTTPRequest, error) {
	// Split the request into lines
	lines := strings.Split(rawRequest, "\r\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty request")
	}

	// Parse request line
	request, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}

	// Parse headers
	headers, headerEndIndex, err := parseHeaders(lines)
	if err != nil {
		return nil, err
	}
	request.Headers = headers

	// Parse body
	body, err := parseBody(lines, headerEndIndex, headers)
	if err != nil {
		return nil, err
	}
	request.Body = body
	request.Raw = rawRequest

	return request, nil
}

func parseRequestLine(requestLine string) (*HTTPRequest, error) {
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line format: %s", requestLine)
	}

	method, path, version := parts[0], parts[1], parts[2]

	// Validate components
	if method == "" || path == "" || version == "" {
		return nil, fmt.Errorf("invalid request line components")
	}

	// Validate HTTP method
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}
	isValidMethod := false
	for _, validMethod := range validMethods {
		if method == validMethod {
			isValidMethod = true
			break
		}
	}
	if !isValidMethod {
		return nil, fmt.Errorf("invalid HTTP method: %s", method)
	}

	return &HTTPRequest{
		Method:  method,
		Path:    path,
		Version: version,
	}, nil
}

func parseHeaders(lines []string) (map[string]string, int, error) {
	headers := make(map[string]string)
	headerEndIndex := 1

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if line == "" {
			headerEndIndex = i
			break
		}

		colonIndex := strings.Index(line, ":")
		if colonIndex == -1 {
			continue
		}

		key := strings.TrimSpace(line[:colonIndex])
		value := strings.TrimSpace(line[colonIndex+1:])

		if key == "" {
			continue
		}

		headers[key] = value
	}

	return headers, headerEndIndex, nil
}

func parseBody(lines []string, headerEndIndex int, headers map[string]string) (string, error) {
	var body string
	if headerEndIndex+1 < len(lines) {
		bodyLines := lines[headerEndIndex+1:]
		body = strings.Join(bodyLines, "\r\n")
	}

	// Handle Content-Length
	if contentLengthStr := headers["Content-Length"]; contentLengthStr != "" {
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return "", fmt.Errorf("invalid Content-Length header: %s", contentLengthStr)
		}
		if contentLength > 0 && len(body) > contentLength {
			body = body[:contentLength]
		}
	}

	return body, nil
}
