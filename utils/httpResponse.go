package utils

import (
	"net/http"
)

type HTTPResponse struct {
	StatusCode int
	Content    []byte
	Headers    http.Header
}
