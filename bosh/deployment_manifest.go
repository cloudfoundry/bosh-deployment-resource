package bosh

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"errors"
)

type DeploymentManifest struct {
	manifest map[interface{}]interface{}
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

func (d DeploymentManifest) UseStemcellVersion(stemcellName, os, version string) error {
	stemcells, ok := d.manifest["stemcells"].([]interface{})
	if !ok {
		return errors.New("No stemcells section in deployment manifest")
	}

	matchingStemcells := []map[interface{}]interface{}{}
	for i := range stemcells {
		stemcell := stemcells[i].(map[interface{}]interface{})
		if stemcell["name"] == stemcellName || stemcell["os"] == os {
			matchingStemcells = append(matchingStemcells, stemcell)
		}
	}

	if len(matchingStemcells) == 0 {
		return fmt.Errorf("Stemcell %s not defined in deployment manifest", stemcellName)
	}

	foundMatch := false
	for _, stemcell := range matchingStemcells {
		if stemcell["version"] == "latest" {
			if !foundMatch {
				stemcell["version"] = version
				foundMatch = true
			} else {
				return fmt.Errorf("Multiple matches for stemcell %s", stemcellName)
			}
		}
	}

	return nil
}

func (d DeploymentManifest) Manifest() []byte {
	bytes, _ := yaml.Marshal(d.manifest)

	return bytes
}