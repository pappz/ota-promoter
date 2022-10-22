package web

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"bitbucket.org/pzoli/ota-promoter/promoter"
)

type Server struct {
	httpServer *http.Server
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
	}

	middleware := NewMiddleware(service)

	router.HandleFunc("/files", middleware.Handle(fileList)).Methods("GET")
	router.HandleFunc("/files/version", middleware.Handle(getVersion)).Methods("GET")
	router.HandleFunc("/files/{checksum}", middleware.Handle(downloadFile)).Methods("GET")

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
