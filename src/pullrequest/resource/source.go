package resource

// Source is
type Source struct {
	URI         string `json:"uri"`
	Insecure    bool   `json:"insecure"`
	AccessToken string `json:"access_token"`
	Repo        string `json:"repo"`
	Owner       string `json:"owner"`
}

// Version is
type Version struct {
	REF string `json:"ref"`
}

// CheckRequest is
type CheckRequest struct {
	Source  `json:"source"`
	Version `json:"version"`
}

// NewCheckRequest is
func NewCheckRequest() CheckRequest {
	req := CheckRequest{}
	return req
}
