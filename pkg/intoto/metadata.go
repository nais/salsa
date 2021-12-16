package intoto

type Metadata struct {
	BuildInvocationId string `json:"buildInvocationId"`
	Completeness      `json:"completeness"`
	Reproducible      bool `json:"reproducible"`
	// BuildStartedOn not defined as it's not available from a GitHub Action.
	BuildFinishedOn string `json:"buildFinishedOn"`
}

type Completeness struct {
	Arguments   bool `json:"arguments"`
	Environment bool `json:"environment"`
	Materials   bool `json:"materials"`
}

func (m Metadata) withCompleteness(arguments, environment, materials bool) Metadata {
	m.Completeness = Completeness{
		Arguments:   arguments,
		Environment: environment,
		Materials:   materials,
	}
	return m
}
