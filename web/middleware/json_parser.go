package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrRespInternalError is generic error for unexpected cases
	ErrRespInternalError = errors.New("internal error")
)

// ErrorResponse is the generic Json format for http error responses
type ErrorResponse struct {
	Message string
}

// JsonParser prepare all necessary data for the handlers and
// manage the json responses and errors
type JsonParser struct {
}

// Handle manage the http headers and status codes in the response.
// If a handler response with a non nil struct then the middleware do the
// Json marshal step and send it out. In case the handler response with
// error the middleware send out the proper status code with the Json error
// message.
func (m JsonParser) Handle(h Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &Request{
			W:   w,
			R:   r,
			Log: log.WithFields(log.Fields{"tag": "web", "address": r.RemoteAddr}),
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
			request.Log.Debug("failed to send out json response: %s", err.Error())
		}

	}
}

// responseError response with error to the Request. It set the proper http headers
// and based on the error type send it out the required error message.
func (m JsonParser) responseError(w http.ResponseWriter, e error) {
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

// responseJson marshal the response content and send out to the http Request with
// the proper headers.
func (m JsonParser) responseJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(j)
	return err
}
