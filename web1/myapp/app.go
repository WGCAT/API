package myapp

import "net/http"

// NewHandler make a new myapp handler
func NewHandler() http.Handler {
	mux := http.NewServeMux()

	return mux
}
