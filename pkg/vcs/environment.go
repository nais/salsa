package vcs

import "time"

type ContextEnvironment interface {
	BuilderId() string
	BuildInvocationId() string
	BuildType() string
	Context() string
	CurrentFilteredEnvironment() map[string]string
	GetBuildStartedOn() time.Time
	NonReproducibleMetadata() *Metadata
	UserDefinedParameters() *Event
	RepoUri() string
	Sha() string
}
