package concourse

type CompiledRelease struct {
	Name string   `json:"name"`
	Jobs []string `json:"jobs"`
}

type InParams struct {
	CompiledReleases []CompiledRelease `json:"compiled_releases,omitempty"`
}
