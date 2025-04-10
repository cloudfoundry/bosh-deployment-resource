package bosh

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

type DeploymentManifest struct {
	manifest map[interface{}]interface{}
}

func NewDeploymentManifest(manifest []byte) (DeploymentManifest, error) {
	d := DeploymentManifest{}

	err := yaml.Unmarshal(manifest, &d.manifest)
	if err != nil {
		return d, fmt.Errorf("Failed to unmarshal manifest: %s", err) //nolint:staticcheck
	}

	return d, nil
}

func (d DeploymentManifest) UseReleaseVersion(releaseName, version string) error {
	releases, ok := d.manifest["releases"].([]interface{})
	if !ok {
		return errors.New("No releases section in deployment manifest") //nolint:staticcheck
	}

	for i := range releases {
		release := releases[i].(map[interface{}]interface{})
		if release["name"] == releaseName {
			release["version"] = version
			return nil
		}
	}

	return fmt.Errorf("Release %s not defined in deployment manifest", releaseName) //nolint:staticcheck
}

func (d DeploymentManifest) UseStemcellVersion(stemcellName, os, version string) error {
	stemcells, ok := d.manifest["stemcells"].([]interface{})
	if !ok {
		return errors.New("No stemcells section in deployment manifest") //nolint:staticcheck
	}

	matchingStemcells := []map[interface{}]interface{}{}
	for i := range stemcells {
		stemcell := stemcells[i].(map[interface{}]interface{})
		if stemcell["name"] == stemcellName || stemcell["os"] == os {
			matchingStemcells = append(matchingStemcells, stemcell)
		}
	}

	if len(matchingStemcells) == 0 {
		return fmt.Errorf("Stemcell %s not defined in deployment manifest", stemcellName) //nolint:staticcheck
	}

	foundMatch := false
	for _, stemcell := range matchingStemcells {
		if stemcell["version"] == "latest" {
			if !foundMatch {
				stemcell["version"] = version
				foundMatch = true
			} else {
				return fmt.Errorf("Multiple matches for stemcell %s", stemcellName) //nolint:staticcheck
			}
		}
	}

	return nil
}

func (d DeploymentManifest) Manifest() []byte {
	bytes, _ := yaml.Marshal(d.manifest) //nolint:errcheck

	return bytes
}

func (d DeploymentManifest) Stemcells() ([]Stemcell, error) {
	stemcells, ok := d.manifest["stemcells"].([]interface{})
	if !ok {
		return nil, errors.New("No stemcells section in deployment manifest") //nolint:staticcheck
	}

	out := make([]Stemcell, 0)
	for i := range stemcells {
		stemcell := stemcells[i].(map[interface{}]interface{})
		os, ok := stemcell["os"].(string)
		if !ok {
			return nil, errors.New("expected os key for stemcell")
		}

		version, ok := stemcell["version"].(string)
		if !ok {
			return nil, errors.New("expected version key for stemcell")
		}

		out = append(out, Stemcell{OperatingSystem: os, Version: version})
	}
	return out, nil
}
