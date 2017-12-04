package concourse

type CompiledRelease struct {
	Name string `json:"name"`
}

type InParams struct {
	CompiledReleases []CompiledRelease `json:"compiled_releases,omitempty"`
}
