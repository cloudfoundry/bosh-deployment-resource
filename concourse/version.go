package concourse

import (
	"crypto/sha1"
	"fmt"
)

type Version struct {
	ManifestSha1 string `json:"manifest_sha1"`
	Target       string `json:"target"`
}

func NewVersion(bytesToSha1 []byte, target string) Version {
	return Version{
		ManifestSha1: fmt.Sprintf("%x", sha1.Sum(bytesToSha1)),
		Target:       target,
	}
}
