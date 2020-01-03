package main

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
	pf, ok := getPromotedFileByChecksum(checksum)
	if !ok {
		r.log.Errorf("File not found by checksum: %s", checksum)
		http.Error(r.w, "File not found by checksum.", http.StatusNotFound)
		return
	}

	openFile, err := os.Open(pf.localPath)
	defer openFile.Close()
	if err != nil {
		r.log.Errorf("Failed to open file: %s", err.Error())
		http.Error(r.w, "File not found.", 404)
		return
	}

	//Send the headers
	r.w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(pf.CanonicalName))
	r.w.Header().Set("Content-Length", strconv.FormatInt(pf.size, 10))
	r.w.Header().Set("X-target-path", pf.CanonicalName)

	_, err = io.Copy(r.w, openFile)
	if err != nil {
		r.log.Errorf("failed to write out the file to the client: %v", err)
	}
	r.log.Infof("download promoted file: %s - %s", checksum, pf.CanonicalName)
}
