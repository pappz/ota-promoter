package web

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/pzoli/ota-promoter/promoter"
)

var (
	// ErrRespInternalError is generic error for unexpected cases
	ErrRespInternalError = errors.New("internal error")
)

// ErrorResponse is the generic Json format for http error responses
type ErrorResponse struct {
	Message string
}

// Middleware prepare all necessary data for the handlers and
// manage the json responses and errors
type Middleware struct {
	service promoter.Promoter
}

// NewMiddleware instantiate a new Middleware
func NewMiddleware(service promoter.Promoter) Middleware {
	return Middleware{
		service: service,
	}
}

// Handle manage the http headers and status codes in the response.
// If a handler response with a non nil struct then the middleware do the
// Json marshal step and send it out. In case the handler response with
// error the middleware send out the proper status code with the Json error
// message.
func (m Middleware) Handle(h Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &request{
			w:       w,
			r:       r,
			service: m.service,
			log:     log.WithFields(log.Fields{"tag": "web", "address": r.RemoteAddr}),
		}

		v, err := h(request)
		if err != nil {
			m.responseError(w, err)
			return
		}

		if v == nil {
			return
		}

		if err := m.responseJson(w, v); err != nil {
			request.log.Debug("failed to send out json response: %s", err.Error())
		}

	}
}

// responseError response with error to the request. It set the proper http headers
// and based on the error type send it out the required error message.
func (m Middleware) responseError(w http.ResponseWriter, e error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if e == ErrRespInternalError {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	resp := ErrorResponse{
		e.Error(),
	}

	// json marshal error never should happen so ignore it
	j, _ := json.Marshal(resp)
	_, _ = w.Write(j)
	return
}

// responseJson marshal the response content and send out to the http request with
// the proper headers.
func (m Middleware) responseJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(j)
	return err
}
