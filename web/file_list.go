package web

import (
	"bitbucket.org/pzoli/ota-promoter/promoter"
)

type Response struct {
	Files   []*promoter.File `json:"files"`
	Version string           `json:"version"`
}

func fileList(r *request) (ResponseData, error) {
	resp := Response{
		Version: r.service.Version(),
		Files:   r.service.PromotedFiles(),
	}
	return resp, nil
}
