package main

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
