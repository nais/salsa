package vcs

type ContextEnvironment interface {
	RepoUri() string
	BuildInvocationId() string
	Sha() string
	BuilderId() string
	UserDefinedParameters() *Event
	CurrentFilteredEnvironment() map[string]string
	NonReproducibleMetadata() *Metadata
	Context() string
}

type BuildEnvironment interface {
	FilterEnvs() map[string]string
	GetEnvs() map[string]string
}
