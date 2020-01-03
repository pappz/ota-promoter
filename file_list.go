package main

import (
	"encoding/json"
)

type Response struct {
	Files   []*PromotedFile `json:"files"`
	Version string          `json:"version"`
}

func fileList(r *request) {
	var resp = &Response{
		Version: version,
		Files:   promotedFiles,
	}

	j, err := json.Marshal(resp)
	if err != nil {
		r.log.Infof("failed to marshal response: %v", err)
	}

	responseJson(r.w, j)
}
