package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Raw     string
	Headers map[string]string
	Body    string
}

type HTTPResponse struct {
	StatusCode int
	StatusText string
	Body       string
	Headers    map[string]string
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("TCP server listening on port 8000")

	// keep loop and send them to go routine
	for {

		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("New Connection From %s\n", clientAddr)

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error Reading from %s\n", clientAddr)
		return
	}

	rawRequest := string(buffer[:n])
	fmt.Printf("Recevied rawRequest: %s\n", rawRequest)

	request, err := parseHTTPRequest(rawRequest)
	if err != nil {
		fmt.Printf("error from parsing the request: %s", err)
		badRequest := "HTTP 1.1 400 Bad Request\r\n\r\nBad Request"
		conn.Write([]byte(badRequest))
		return
	}

	// Send HTTP response with proper headers
	response := routeRequest(request)

	err = sendResponse(conn, response)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Sent HTTP Response to %s\n", clientAddr)
}

func routeRequest(req *HTTPRequest) *HTTPResponse {
	switch req.Path {
	case "/":
		return handleHome(req)
	default:
		return createErrResponse(404, "Not Found")
	}
}

func handleHome(req *HTTPRequest) *HTTPResponse {
	if req.Method != "GET" {
		return createErrResponse(405, "Method Not Allowed")
	}

	body := "Welcome to our HTTP Server!\nAvailable endpoints:\n- GET /\n- GET /users\n- POST /users\n- GET /health"
	return createResponse(200, "OK", "text/plain", body)
}

func createResponse(statusCode int, statusText, contentType, body string) *HTTPResponse {
	headers := map[string]string{
		"Content-Type":   contentType,
		"Content-Length": strconv.Itoa(len(body)),
		"Date":           time.Now().Format(time.RFC1123),
		"Server":         "CustomHTTPServer/1.0",
	}

	return &HTTPResponse{
		StatusCode: statusCode,
		StatusText: statusText,
		Headers:    headers,
		Body:       body,
	}
}

func createErrResponse(statusCode int, statusText string) *HTTPResponse {
	errBody := map[string]any{
		"error":   statusText,
		"code":    statusCode,
		"message": getErrorMessage(statusCode),
	}

	jsonBody, _ := json.Marshal(errBody)
	return createResponse(statusCode, statusText, "application/json", string(jsonBody))
}

func getErrorMessage(statusCode int) string {
	switch statusCode {
	case 400:
		return "The request was invalid"
	case 404:
		return "The requested resource was not found"
	case 405:
		return "The HTTP method is not allowed for this resource"
	case 500:
		return "Internal server error"
	default:
		return "An error occurred"
	}
}

func sendResponse(conn net.Conn, response *HTTPResponse) error {
	var responseStr strings.Builder

	responseStr.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", response.StatusCode, response.StatusText))

	for kay, value := range response.Headers {
		responseStr.WriteString(fmt.Sprintf("%s:%s\r\n", kay, value))
	}

	responseStr.WriteString("\r\n")
	responseStr.WriteString(response.Body)

	_, err := conn.Write([]byte(responseStr.String()))
	return err
}

func parseHTTPRequest(rawRequest string) (*HTTPRequest, error) {
	lines := strings.Split(rawRequest, "\r\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty request")
	}

	requestLine := lines[0]

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid requestLine format: %s", requestLine)
	}

	method := parts[0]
	path := parts[1]
	version := parts[2]

	if method == "" || path == "" || version == "" {
		return nil, fmt.Errorf("invalid http request")
	}

	headers := make(map[string]string)
	headerEndIndex := 1

	for i := 1; i <= len(lines); i++ {
		line := lines[i]

		if line == "" {
			headerEndIndex = i
			break
		}

		colonindex := strings.Index(lines[i], ":")
		if colonindex == -1 {
			// skip maiformed header line
			continue
		}

		key := strings.TrimSpace(line[:colonindex])
		value := strings.TrimSpace(line[colonindex+1:])

		headers[key] = value
	}

	var body string
	if headerEndIndex+1 < len(lines) {
		bodylines := lines[headerEndIndex+1:]
		body = strings.Join(bodylines, " ")
	}

	contentLengthStr := headers["Content-Length"]
	contentLength, err := strconv.Atoi(contentLengthStr)
	if err == nil && contentLength > 0 {
		// to get exact body content length
		if contentLength > len(body) {
			body = body[:contentLength]
		}
	}

	return &HTTPRequest{
		Method:  method,
		Path:    path,
		Version: version,
		Raw:     rawRequest,
		Headers: headers,
		Body:    body,
	}, nil

}
