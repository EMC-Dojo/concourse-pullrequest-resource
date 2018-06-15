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
	Ref string `json:"ref"`
}

// Metadata is
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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

// InParams is
type InParams struct {
	Globs                []string `json:"globs"`
	IncludeSourceTarball bool     `json:"include_source_tarball"`
	IncludeSourceZip     bool     `json:"include_source_zip"`
}

// InRequest is
type InRequest struct {
	Source   `json:"source"`
	Version  *Version `json:"version"`
	InParams `json:"params"`
}

// InResponse is
type InResponse struct {
	Version  `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

// NewInRequest is
func NewInRequest() InRequest {
	return InRequest{}
}

// OutParams is
type OutParams struct {
	NamePath string `json:"name"`
	BodyPath string `json:"body"`
}

// OutRequest is
type OutRequest struct {
	Source    `json:"source"`
	OutParams `json:"params"`
}

// OutResponse is
type OutResponse struct {
	Version  `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

// NewOutRequest is
func NewOutRequest() OutRequest {
	return OutRequest{}
}
