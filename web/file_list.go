package web

import (
	"encoding/json"

	"bitbucket.org/pzoli/ota-promoter/promoter"
)

type Response struct {
	Files   []*promoter.File `json:"files"`
	Version string           `json:"version"`
}

func fileList(r *request) {
	var resp = &Response{
		Version: r.service.Version(),
		Files:   r.service.PromotedFiles(),
	}

	j, err := json.Marshal(resp)
	if err != nil {
		r.log.Infof("failed to marshal response: %v", err)
	}

	responseJson(r.w, j)
}
