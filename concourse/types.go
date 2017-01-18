package concourse

import (
	"crypto/sha1"
	"fmt"
)

type Source struct {
	Deployment   string `json:"deployment"`
	Client       string `json:"client"`
	ClientSecret string `json:"client_secret"`
	Target       string `json:"target"`
	CACert       string `json:"ca_cert"`
}

type Version struct {
	ManifestSha1 string `json:"manifest_sha1"`
	Target       string `json:"target"`
}

type OutParams struct {
	Manifest   string `json:"manifest"`
	TargetFile string `json:"target_file,omitempty"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

func NewVersion(bytesToSha1 []byte, target string) Version {
	return Version{
		ManifestSha1: fmt.Sprintf("%x", sha1.Sum(bytesToSha1)),
		Target:       target,
	}
}
