package vcs

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

func (in *Environment) NonReproducibleMetadata() *Metadata {
	// Other variables that are required to reproduce the build and that cannot be
	// recomputed using existing information.
	//(Documentation would explain how to recompute the rest of the fields.)
	return &Metadata{
		Arch: in.RunnerContext.Arch,
		Env:  in.GetCurrentFilteredEnvironment(),
		Context: Context{
			Github: Github{
				RunId: in.GitHubContext.RunId,
			},
			Runner: Runner{
				Os:   in.RunnerContext.OS,
				Temp: in.RunnerContext.Temp,
			},
		},
	}
}
