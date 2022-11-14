package api

import (
	"github.com/gorilla/mux"

	"github.com/pappz/ota-promoter/promoter"
	"github.com/pappz/ota-promoter/web/middleware"
)

// RegisterFileListHandler sets up the routing of the HTTP handlers.
func RegisterFileListHandler(router *mux.Router, service *promoter.Promoter) {
	m := middleware.JsonParser{}
	h := fileListHandler{service}
	router.HandleFunc("/files", m.Handle(h.handle)).Methods("GET")
}

type fileListHandler struct {
	service *promoter.Promoter
}

type ResponseFileList struct {
	Files   []*promoter.File `json:"files"`
	Version string           `json:"version"`
}

func (req fileListHandler) handle(r *middleware.Request) (middleware.ResponseData, error) {
	resp := ResponseFileList{
		Version: req.service.Version(),
		Files:   req.service.PromotedFiles(),
	}
	return resp, nil
}
