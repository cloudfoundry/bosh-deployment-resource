package concourse

type OutParams struct {
	Manifest  string                 `json:"manifest"`
	NoRedact  bool                   `json:"no_redact,omitempty"`
	Cleanup   bool                   `json:"cleanup,omitempty"`
	Releases  []string               `json:"releases,omitempty"`
	Stemcells []string               `json:"stemcells,omitempty"`
	Vars      map[string]interface{} `json:"vars,omitempty"`
	VarsFiles []string               `json:"vars_files,omitempty"`
	OpsFiles  []string               `json:"ops_files,omitempty"`
}
