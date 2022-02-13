package intoto

import (
	"github.com/nais/salsa/pkg/digest"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/vcs"
)

type ProvenanceOptions struct {
	BuildConfig       *BuildConfig
	BuilderId         string
	BuilderRepoDigest *slsa.ProvenanceMaterial
	BuildInvocationId string
	BuildStartedOn    time.Time
	BuildType         string
	Dependencies      *build.ArtifactDependencies
	Invocation        *slsa.ProvenanceInvocation
	Name              string
}

type BuildConfig struct {
	Commands []string `json:"commands"`
	// Indicates how to parse the strings in commands.
	Shell string `json:"shell"`
}

func CreateProvenanceOptions(name string, deps *build.ArtifactDependencies, env *vcs.Environment) *ProvenanceOptions {
	opts := &ProvenanceOptions{
		BuildStartedOn: time.Now().UTC(),
		Dependencies:   deps,
		Name:           name,
	}

	if env != nil {
		opts.BuildType = vcs.BuildType
		opts.BuildInvocationId = env.BuildInvocationId()
		opts.BuilderId = env.BuilderId()
		opts.withBuilderRepoDigest(env).withBuilderInvocation(env)
		return opts
	}

	opts.BuildConfig = &BuildConfig{
		Commands: []string{"make salsa"},
		Shell:    "bash",
	}
	opts.BuilderId = vcs.DefaultBuildId
	opts.BuildType = vcs.AdHocBuildType
	opts.Invocation = nil
	return opts
}

func (in *ProvenanceOptions) withBuilderRepoDigest(env *vcs.Environment) *ProvenanceOptions {
	in.BuilderRepoDigest = &slsa.ProvenanceMaterial{
		URI: "git+" + env.RepoUri(),
		Digest: slsa.DigestSet{
			digest.AlgorithmSHA1: env.GithubSha(),
		},
	}
	return in
}

func (in *ProvenanceOptions) withBuilderInvocation(env *vcs.Environment) *ProvenanceOptions {
	in.Invocation = &slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI: "git+" + env.RepoUri(),
			Digest: slsa.DigestSet{
				digest.AlgorithmSHA1: env.GithubSha(),
			},
			EntryPoint: env.Workflow,
		},
		Parameters:  env.EventInputs(),
		Environment: env.FilteredEnvironment(),
	}
	return in
}

func (in *ProvenanceOptions) HasBuilderRepoDigest() bool {
	if in.BuilderRepoDigest == nil {
		return false
	}

	return in.BuilderRepoDigest.Digest != nil && in.BuilderRepoDigest.URI != ""

}

func (in *ProvenanceOptions) HasDependencies() bool {
	if in.Dependencies == nil {
		return false
	}

	return len(in.Dependencies.RuntimeDeps) > 0
}

func (in *ProvenanceOptions) HasParameters() bool {
	if in.Invocation == nil {
		return false
	}

	return in.Invocation.Parameters != nil
}

func (in *ProvenanceOptions) HasEnvironment() bool {
	if in.Invocation == nil {
		return false
	}

	return in.Invocation.Environment != nil
}
