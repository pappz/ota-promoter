package main

import (
	"encoding/json"
)

type ResponseVersion struct {
	Version string `json:"version"`
}

func getVersion(r *request) {
	v := &ResponseVersion{
		Version: version,
	}

	if j, err := json.Marshal(v); err == nil {
		responseJson(r.w, j)
	}
	r.log.Infof("get version information (%s)", version)
}
