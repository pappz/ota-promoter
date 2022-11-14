package api

import (
	"github.com/gorilla/mux"

	"github.com/pappz/ota-promoter/promoter"
	"github.com/pappz/ota-promoter/web/middleware"
)

// RegisterVersionHandler sets up the routing of the HTTP handlers.
func RegisterVersionHandler(router *mux.Router, service *promoter.Promoter) {
	h := versionHandler{service}
	router.HandleFunc("/files/version", middleware.Handle(h.handle)).Methods("GET")
}

type ResponseVersion struct {
	Version string `json:"version"`
}

type versionHandler struct {
	service *promoter.Promoter
}

func (req versionHandler) handle(r *middleware.Request) (middleware.ResponseData, error) {
	v := ResponseVersion{
		Version: req.service.Version(),
	}

	r.Log.Infof("get version information (%s)", v.Version)
	return v, nil
}
