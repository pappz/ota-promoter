package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	Router = mux.NewRouter()
	server http.Server
)

type request struct {
	w   http.ResponseWriter
	r   *http.Request
	log *logrus.Entry
}

type handler func(w http.ResponseWriter, r *http.Request)
type extendedHandler func(r *request)

func responseJson(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func shutDownServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown the server:%+v", err)
	}
}

func listen() {
	Router.HandleFunc("/files", requestMiddleware(fileList)).Methods("GET")
	Router.HandleFunc("/files/version", requestMiddleware(getVersion)).Methods("GET")
	Router.HandleFunc("/files/{checksum}", requestMiddleware(downloadFile)).Methods("GET")

	server = http.Server{
		Addr:         fmt.Sprintf("%s:%d", listenAddress, listenPort),
		Handler:      Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Infof("promoter server is listening on http://%s ...", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func requestMiddleware(h extendedHandler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &request{
			w:   w,
			r:   r,
			log: log.WithFields(logrus.Fields{"address": r.RemoteAddr}),
		}
		h(request)
	}
}
