package concourse

type VarsStore struct {
	Provider string `json:"storage"`
	Config   []byte `json:"config"`
}
