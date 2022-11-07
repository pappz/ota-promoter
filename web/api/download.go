package api

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"

	"bitbucket.org/pzoli/ota-promoter/promoter"
	"bitbucket.org/pzoli/ota-promoter/web/middleware"
)

var (
	errFileNotFound = errors.New("file not found")
)

// RegisterDownloadHandler sets up the routing of the HTTP handlers.
func RegisterDownloadHandler(router *mux.Router, service *promoter.Promoter) {
	m := middleware.JsonParser{}
	h := downloadHandler{service}
	router.HandleFunc("/files/{checksum}", m.Handle(h.handle)).Methods("GET")
}

type downloadHandler struct {
	service *promoter.Promoter
}

func (req downloadHandler) handle(r *middleware.Request) (middleware.ResponseData, error) {
	params := mux.Vars(r.R)
	checksum := params["checksum"]
	pf, ok := req.service.PromotedFileByChecksum(checksum)
	if !ok {
		r.Log.Errorf("file not found by checksum: %s", checksum)
		return nil, errFileNotFound
	}

	openFile, err := os.Open(pf.LocalPath)
	defer func(openFile *os.File) {
		_ = openFile.Close()
	}(openFile)

	if err != nil {
		r.Log.Errorf("Failed to open file: %s", err.Error())
		return nil, middleware.ErrRespInternalError
	}

	//Send the headers
	r.W.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(pf.PromotedPath))
	r.W.Header().Set("Content-Length", strconv.FormatInt(pf.Size, 10))
	r.W.Header().Set("X-target-path", pf.PromotedPath)

	_, err = io.Copy(r.W, openFile)
	if err != nil {
		r.Log.Errorf("failed to write out the file to the client: %v", err)
	}
	r.Log.Infof("download promoted file: %s - %s", checksum, pf.PromotedPath)
	return nil, nil
}
