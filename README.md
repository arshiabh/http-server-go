# Custom HTTP Server in Go

A fully functional HTTP server built from scratch in Go without using the standard `net/http` package. This project demonstrates a deep understanding of the HTTP protocol by implementing all core features manually using only TCP connections.

## 🚀 Features

- **Pure TCP Implementation** - Built using only `net.Conn` and raw TCP sockets
- **Full HTTP Protocol Support** - Complete HTTP/1.1 request/response handling
- **RESTful API** - CRUD operations with proper status codes
- **JSON Support** - Automatic JSON marshaling/unmarshaling
- **Concurrent Connections** - Goroutine-based connection handling
- **Structured Logging** - Comprehensive logging with different levels
- **Error Recovery** - Panic recovery and graceful error handling
- **Request Routing** - Pattern-based URL routing system
- **Configuration Management** - Environment variable support
- **Timeout Handling** - Read/write timeouts for connection safety

## 🏗️ Architecture

The server is organized into several key components:

- **HTTPServer** - Main server struct managing connections and configuration
- **Router** - Handles URL routing and pattern matching
- **Logger** - Structured logging with multiple levels
- **Config** - Configuration management with environment variable support
- **Request/Response** - HTTP message parsing and formatting

## 📋 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/`      | Welcome page with available endpoints |
| GET    | `/users` | Get all users (JSON) |
| POST   | `/users` | Create a new user |
| GET    | `/users/{id}` | Get user by ID |
| PUT    | `/users/{id}` | Update user by ID |
| DELETE | `/users/{id}` | Delete user by ID |
| GET    | `/health` | Health check endpoint |
| GET    | `/error` | Error testing endpoint |

## 🚀 Quick Start

### Prerequisites

- Go 1.16 or higher

### Installation & Running

1. Clone the repository:
```bash
git clone https://github.com/arshiabh/http-server-go.git
cd http-server-go
```

2. Run the server:
```bash
go run main.go
```

3. The server will start on `localhost:8000` by default.

### Configuration

Configure the server using environment variables:

```bash
# Custom port and host
HTTP_PORT=9000 HTTP_HOST=0.0.0.0 go run main.go

# Enable debug logging
LOG_LEVEL=DEBUG go run main.go
```

Available configuration options:
- `HTTP_HOST` - Host to bind to (default: localhost)
- `HTTP_PORT` - Port to listen on (default: 8080)
- `LOG_LEVEL` - Logging level: INFO, DEBUG, ERROR (default: INFO)

## 📖 Usage Examples

### Basic GET Request
```bash
curl http://localhost:8000/
```

### Get All Users
```bash
curl http://localhost:8000/users
```

### Create a User
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}' \
  http://localhost:8080/users
```

### Get User by ID
```bash
curl http://localhost:8000/users/123
```

### Update User
```bash
curl -X PUT \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}' \
  http://localhost:8080/users/123
```

### Delete User
```bash
curl -X DELETE http://localhost:8000/users/123
```

### Health Check
```bash
curl http://localhost:8000/health
```

## 🔧 Technical Implementation

### HTTP Request Parsing

The server manually parses HTTP requests by:
1. Reading raw TCP data
2. Splitting request lines by `\r\n`
3. Parsing the request line (Method, Path, Version)
4. Extracting headers (key-value pairs)
5. Reading the request body based on Content-Length

### Response Generation

HTTP responses are manually formatted with:
- Status line (HTTP/1.1 200 OK)
- Headers (Content-Type, Content-Length, Date, etc.)
- Empty line separator
- Response body

### Routing System

The router supports:
- Exact path matching (`/users`)
- Pattern matching (`/users/{id}`)
- HTTP method validation
- Handler function registration

## 🐛 Error Handling

The server includes comprehensive error handling:

- **400 Bad Request** - Invalid HTTP syntax
- **404 Not Found** - Unknown endpoints
- **405 Method Not Allowed** - Wrong HTTP method
- **408 Request Timeout** - Slow clients
- **500 Internal Server Error** - Server panics

All errors return structured JSON responses:

```json
{
  "error": "Not Found",
  "code": 404,
  "message": "The requested resource was not found",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 📊 Logging

The server provides structured logging with different levels:

```
[HTTP-SERVER] [INFO] HTTP Server starting on localhost:8000
[HTTP-SERVER] [INFO] New connection from: 127.0.0.1:54321
[HTTP-SERVER] [INFO] Request: GET /users from 127.0.0.1:54321
[HTTP-SERVER] [INFO] Response: 200 OK to 127.0.0.1:54321
```

Enable debug logging to see raw HTTP requests:
```bash
LOG_LEVEL=DEBUG go run main.go
```

## 🧪 Testing

Test the server with various tools:

### Using curl
```bash
# Test all endpoints
curl http://localhost:8000/
curl http://localhost:8000/users
curl -X POST -d '{"name":"Test"}' http://localhost:8000/users
```

### Using a web browser
Navigate to `http://localhost:8000` in your browser.

### Testing error scenarios
```bash
# Test 404
curl http://localhost:8000/nonexistent

# Test 405
curl -X POST http://localhost:8000/health

# Test invalid JSON
curl -X POST -d 'invalid json' http://localhost:8000/users
```

## 🔄 Development Steps

This server was built incrementally through these steps:

1. **Basic TCP Server** - Accept connections and echo data
2. **HTTP Request Line Parsing** - Parse Method, Path, Version
3. **HTTP Header Parsing** - Extract request headers
4. **Request Body Handling** - Read POST data and JSON
5. **Structured HTTP Responses** - Proper status codes and headers
6. **Error Handling & Logging** - Comprehensive error management
7. **Code Organization** - Professional structure and configuration

## 🚧 Limitations

This is an educational project. For production use, consider:

- HTTP/2 and HTTP/3 support
- TLS/SSL encryption
- Performance optimizations
- More robust request parsing
- Security features (rate limiting, input validation)
- Static file serving
- WebSocket support

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🎓 Learning Outcomes

By building this server, you'll understand:

- How HTTP protocol works at the TCP level
- Request/response message formatting
- Connection handling and concurrency
- Error handling in networked applications
- Go's networking capabilities
- Server architecture and design patterns

## 🔗 Resources

- [HTTP/1.1 RFC 7230](https://tools.ietf.org/html/rfc7230)
- [Go net package documentation](https://golang.org/pkg/net/)
- [HTTP Status Codes](https://httpstatuses.com/)

---

**Built with ❤️ and Go** - A deep dive into HTTP protocol implementation