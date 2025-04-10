package bosh

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	boshIOAPIURL = "https://bosh.io/api/v1/stemcells/%s"
)

var (
	cpiStemcellMap = map[string]string{
		"google_cpi":     "google-kvm",
		"alicloud_cpi":   "alicloud-kvm",
		"vcloud_cpi":     "vcloud-esxi",
		"aws_cpi":        "aws-xen-hvm",
		"openstack_cpi":  "openstack-kvm",
		"virtualbox_cpi": "vsphere-esxi",
		"docker_cpi":     "warden-boshlite",
		"vsphere_cpi":    "vsphere-esxi",
		"azure_cpi":      "azure-hyperv",
		"warden_cpi":     "warden-boshlite",
	}
)

type BoshIOStemcell struct {
	Name    string
	Version string
	URL     string
	Sha1    string
}

//go:generate counterfeiter . BoshIO
type BoshIO interface {
	Stemcells(name string) ([]byte, error)
}

type BoshIOClient struct {
}

func (c BoshIOClient) Stemcells(name string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf(boshIOAPIURL, name))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() //nolint:errcheck
	return io.ReadAll(resp.Body)
}

func LookupBoshIOStemcell(c BoshIO, cpi, os, version string, light bool) (BoshIOStemcell, error) {
	if version == "latest" {
		return BoshIOStemcell{},
			errors.New("Auto upload of \"latest\" stemcell is not support, please use bosh-io-stemcell-resource") //nolint:staticcheck
	}
	name, err := stemcellName(cpi, os)
	if err != nil {
		return BoshIOStemcell{}, err
	}

	stemcells, err := c.Stemcells(name)
	if err != nil {
		return BoshIOStemcell{}, err
	}

	return filterStemcells(stemcells, version, light)
}

func filterStemcells(raw []byte, version string, light bool) (BoshIOStemcell, error) {
	var stemcells []struct {
		Name    string
		Version string
		Regular struct {
			URL  string
			Sha1 string
		}
		Light struct {
			URL  string
			Sha1 string
		}
	}

	err := json.Unmarshal(raw, &stemcells)
	if err != nil {
		return BoshIOStemcell{}, err
	}

	for _, s := range stemcells {
		if s.Version == version {
			out := BoshIOStemcell{
				Name:    s.Name,
				Version: s.Version,
			}
			if light {
				out.URL = s.Light.URL
				out.Sha1 = s.Light.Sha1
			} else {
				out.URL = s.Regular.URL
				out.Sha1 = s.Regular.Sha1
			}
			return out, nil
		}
	}
	return BoshIOStemcell{}, fmt.Errorf("did not find a suitable stemcell with version: %s", version)
}

func stemcellName(cpi, os string) (string, error) {
	name, ok := cpiStemcellMap[cpi]
	if !ok {
		return "", fmt.Errorf("Failed to determine stemcell name for cpi: %s", cpi) //nolint:staticcheck
	}
	return fmt.Sprintf("bosh-%s-%s-go_agent", name, os), nil
}
