package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	config := DefaultConfig()

	// Override with environment variables if available
	if port := os.Getenv("HTTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Port = p
		}
	}
	if host := os.Getenv("HTTP_HOST"); host != "" {
		config.Host = host
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	// Create and start server
	server := NewHTTPServer(config)

	// Start server (this blocks)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
