package web

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/pzoli/ota-promoter/promoter"
)

// request will be passed to all controllers.
type request struct {
	w       http.ResponseWriter
	r       *http.Request
	service promoter.Promoter
	log     *log.Entry
}

// ResponseData will be sent out by the middleware to the http
// request as a http response. It could be nil if the handler want
// to send out nothing
type ResponseData interface{}

// Handler interfaces used by the middleware.
type Handler func(r *request) (ResponseData, error)
