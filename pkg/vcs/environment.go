package vcs

type ContextEnvironment interface {
	BuilderId() string
	BuildInvocationId() string
	BuildType() string
	Context() string
	CurrentFilteredEnvironment() map[string]string
	NonReproducibleMetadata() *Metadata
	UserDefinedParameters() *EventInput
	RepoUri() string
	Sha() string
	GetHeadCommitTime() string
}
