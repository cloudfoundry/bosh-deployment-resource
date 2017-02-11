package bosh

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"errors"
)

type DeploymentManifest struct {
	manifest map[interface{}]interface{}
}

type ManifestRelease struct {
	Name    string
	Version string

	URL  string
	SHA1 string
}

func NewDeploymentManifest(manifest []byte) (DeploymentManifest, error) {
	d := DeploymentManifest{}

	err := yaml.Unmarshal(manifest, &d.manifest)
	if err != nil {
		return d, fmt.Errorf("Failed to unmarshal manifest: %s", err)
	}

	return d, nil
}

func (d DeploymentManifest) UseReleaseVersion(releaseName, version string) error {
	releases, ok := d.manifest["releases"].([]interface{})
	if !ok {
		return errors.New("No releases section in deployment manifest")
	}

	for i := range releases {
		release := releases[i].(map[interface{}]interface{})
		if release["name"] == releaseName {
			release["version"] = version
			return nil
		}
	}

	return fmt.Errorf("Release %s not defined in deployment manifest", releaseName)
}

func (d DeploymentManifest) Manifest() []byte {
	bytes, _ := yaml.Marshal(d.manifest)

	return bytes
}