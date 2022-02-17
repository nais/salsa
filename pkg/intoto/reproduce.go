package intoto

type Metadata struct {
	Arch    string            `json:"arch"`
	Env     map[string]string `json:"env"`
	Context Context           `json:"context"`
}

type Context struct {
	Github Github
	Runner Runner
}

type Runner struct {
	Os   string
	Temp string
}

type Github struct {
	RunId string
}
