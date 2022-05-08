package main

import (
	"github.com/ffmoyano/goApi/util"
	"net/http"
)

type response map[string]string
type Method int

const (
	Get Method = iota
	Post
	Put
	Head
	Patch
	Connect
	Delete
	Options
	Trace
)

// Type returns Method string value
func (method Method) Type() string {
	return [...]string{"GET", "POST", "PUT", "HEAD", "PATCH", "CONNECT", "DELETE", "OPTIONS", "TRACE"}[method]
}

// isValidMethod checks if call to a handler is made with one of the valid http methods for that route.
func isValidMethod(next http.HandlerFunc, methods ...Method) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if r.Method == method.Type() {
				next.ServeHTTP(w, r)
				break
			} else {
				util.Response(w, http.StatusMethodNotAllowed, response{"Error": "Method " + r.Method + "  not allowed."})
			}
		}
	}
}
