package web

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/pappz/ota-promoter/promoter"
	"github.com/pappz/ota-promoter/web/api"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(listenAddress string, service *promoter.Promoter) Server {
	router := mux.NewRouter()
	api.RegisterVersionHandler(router, service)
	api.RegisterFileListHandler(router, service)
	api.RegisterDownloadHandler(router, service)

	httpServer := &http.Server{
		Addr:         listenAddress,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s := Server{
		httpServer: httpServer,
	}
	return s
}

func (s *Server) Listen() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen error: '%s'", err)
		}
	}()
}

func (s *Server) ShutDownServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
