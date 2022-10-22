package web

type ResponseVersion struct {
	Version string `json:"version"`
}

func getVersion(r *request) (ResponseData, error) {
	v := ResponseVersion{
		Version: r.service.Version(),
	}

	r.log.Infof("get version information (%s)", v.Version)
	return v, nil
}
