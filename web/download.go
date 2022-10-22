package web

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func downloadFile(r *request) {
	params := mux.Vars(r.r)
	checksum := params["checksum"]
	pf, ok := r.service.PromotedFileByChecksum(checksum)
	if !ok {
		r.log.Errorf("file not found by checksum: %s", checksum)
		http.Error(r.w, "file not found by checksum.", http.StatusNotFound)
		return
	}

	openFile, err := os.Open(pf.LocalPath)
	defer func(openFile *os.File) {
		_ = openFile.Close()
	}(openFile)

	if err != nil {
		r.log.Errorf("Failed to open file: %s", err.Error())
		http.Error(r.w, "File not found.", 404)
		return
	}

	//Send the headers
	r.w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(pf.PromotedPath))
	r.w.Header().Set("Content-Length", strconv.FormatInt(pf.Size, 10))
	r.w.Header().Set("X-target-path", pf.PromotedPath)

	_, err = io.Copy(r.w, openFile)
	if err != nil {
		r.log.Errorf("failed to write out the file to the client: %v", err)
	}
	r.log.Infof("download promoted file: %s - %s", checksum, pf.PromotedPath)
}
