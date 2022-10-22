package web

import (
	"encoding/json"
)

type ResponseVersion struct {
	Version string `json:"version"`
}

func getVersion(r *request) {
	v := &ResponseVersion{
		Version: r.service.Version(),
	}

	if j, err := json.Marshal(v); err == nil {
		responseJson(r.w, j)
	}
	r.log.Infof("get version information (%s)", v.Version)
}
