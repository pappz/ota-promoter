package web

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"bitbucket.org/pzoli/ota-promoter/promoter"
)

type request struct {
	w       http.ResponseWriter
	r       *http.Request
	service promoter.Promoter
	log     *log.Entry
}

func responseJson(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

type Server struct {
	httpServer *http.Server
	service    promoter.Promoter
}

func NewServer(listenAddress string, service promoter.Promoter) Server {
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:         listenAddress,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s := Server{
		httpServer: httpServer,
		service:    service,
	}

	router.HandleFunc("/files", s.handler(fileList)).Methods("GET")
	router.HandleFunc("/files/version", s.handler(getVersion)).Methods("GET")
	router.HandleFunc("/files/{checksum}", s.handler(downloadFile)).Methods("GET")

	return s
}

func (s *Server) Listen() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("listen error: '%s'", err)
		}
	}()
}

func (s *Server) ShutDownServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handler(h func(r *request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &request{
			w:       w,
			r:       r,
			service: s.service,
			log:     log.WithFields(log.Fields{"tag": "web", "address": r.RemoteAddr}),
		}
		h(request)
	}
}
