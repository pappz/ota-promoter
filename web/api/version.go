package api

import (
	"github.com/gorilla/mux"

	"bitbucket.org/pzoli/ota-promoter/promoter"
	"bitbucket.org/pzoli/ota-promoter/web/middleware"
)

// RegisterVersionHandler sets up the routing of the HTTP handlers.
func RegisterVersionHandler(router *mux.Router, service promoter.Promoter) {
	m := middleware.JsonParser{}
	h := versionHandler{service}
	router.HandleFunc("/files/version", m.Handle(h.handle)).Methods("GET")
}

type ResponseVersion struct {
	Version string `json:"version"`
}

type versionHandler struct {
	service promoter.Promoter
}

func (req versionHandler) handle(r *middleware.Request) (middleware.ResponseData, error) {
	v := ResponseVersion{
		Version: req.service.Version(),
	}

	r.Log.Infof("get version information (%s)", v.Version)
	return v, nil
}
