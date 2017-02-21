package concourse

type VarsStore struct {
	Provider string                 `json:"provider"`
	Config   map[string]interface{} `json:"config"`
}
