package concourse

type Source struct {
	Deployment   string `json:"deployment"`
	Client       string `json:"client"`
	ClientSecret string `json:"client_secret"`
	Target       string `json:"target"`
	CACert       string `json:"ca_cert"`
}

type OutParams struct {
	Manifest string `json:"manifest"`
}

type OutRequest struct {
	Params OutParams `json:"params"`
	Source Source    `json:"source"`
}