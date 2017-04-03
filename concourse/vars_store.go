package concourse

type VarsStore struct {
	Provider string                 `json:"provider,omitempty" yaml:"provider"`
	Config   map[string]interface{} `json:"config,omitempty" yaml:"config"`
}
