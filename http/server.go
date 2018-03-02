package http

import "net/http"

// ListenAndServe just copies the http package ListenAndServe function
// later on I should refactor this to actually create my own server and handlers for now yolo
func ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}
