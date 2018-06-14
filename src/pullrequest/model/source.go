package model

type Input struct {
	Source `json:"source"`
	Version `json:"version"`
}

type Source struct {
	URI string `json:"uri"`
}

type Version struct {
	REF string `json:"ref"`
}
