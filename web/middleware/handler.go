package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Request will be passed to all controllers.
type Request struct {
	W   http.ResponseWriter
	R   *http.Request
	Log *log.Entry
}

// ResponseData will be sent out by the middleware to the http
// Request as a http response. It could be nil if the handler want
// to send out nothing
type ResponseData interface{}

// Handler interfaces used by the middleware.
type Handler func(r *Request) (ResponseData, error)
