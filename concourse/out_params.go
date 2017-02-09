package concourse

type OutParams struct {
	Manifest  string    `json:"manifest"`
	NoRedact  bool      `json:"no_redact,omitempty"`
	Cleanup   bool      `json:"cleanup,omitempty"`
	Releases  []string  `json:"releases,omitempty"`
	Stemcells []string  `json:"stemcells,omitempty"`
}
