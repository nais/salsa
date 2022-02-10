package intoto

import (
	"github.com/mitchellh/mapstructure"
	"github.com/nais/salsa/pkg/digest"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/vcs"
)

type ProvenanceArtifact struct {
	BuildConfig       string
	BuilderId         string
	BuilderRepoDigest *slsa.ProvenanceMaterial
	BuildInvocationId string
	BuildStartedOn    time.Time
	BuildType         string
	Dependencies      *build.ArtifactDependencies
	Environment       *vcs.Environment
	Invocation        slsa.ProvenanceInvocation
	Name              string
}

func CreateProvenanceArtifact(name string, deps *build.ArtifactDependencies, env *vcs.Environment) *ProvenanceArtifact {
	if env == nil {
		return &ProvenanceArtifact{
			BuildConfig:    "Some commands that made this build",
			BuilderId:      vcs.DefaultBuildId,
			BuildStartedOn: time.Now().UTC(),
			BuildType:      vcs.AdHocBuildType,
			Dependencies:   deps,
			Name:           name,
		}
	}

	pa := &ProvenanceArtifact{
		BuildType:      vcs.BuildType,
		BuildStartedOn: time.Now().UTC(),
		Dependencies:   deps,
		Environment:    env,
		Name:           name,
	}

	pa.BuildInvocationId = pa.Environment.BuildInvocationId()
	pa.BuilderId = pa.Environment.BuilderId()

	return pa.withBuilderRepoDigest().withBuilderInvocation()
}

func (in *ProvenanceArtifact) withBuilderRepoDigest() *ProvenanceArtifact {
	in.BuilderRepoDigest = &slsa.ProvenanceMaterial{
		URI: "git+" + in.Environment.RepoUri(),
		Digest: slsa.DigestSet{
			digest.AlgorithmSHA1: in.Environment.GithubSha(),
		},
	}
	return in
}

func (in *ProvenanceArtifact) withBuilderInvocation() *ProvenanceArtifact {
	in.Invocation = slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI: "git+" + in.Environment.RepoUri(),
			Digest: slsa.DigestSet{
				digest.AlgorithmSHA1: in.Environment.GithubSha(),
			},
			EntryPoint: in.Environment.Workflow,
		},
		Parameters: in.Environment.EventBytes(),
		// Should contain the architecture of the runner.
		Environment: in.Environment.RunnerContext,
	}
	return in
}

func (in *ProvenanceArtifact) HasLegitBuilderRepoDigest() bool {
	if in.BuilderRepoDigest == nil {
		return false
	}

	return in.BuilderRepoDigest.Digest != nil && in.BuilderRepoDigest.URI != ""

}

func (in *ProvenanceArtifact) HasLegitDependencies() bool {
	if in.Dependencies == nil {
		return false
	}

	return len(in.Dependencies.RuntimeDeps) > 0
}

func (in *ProvenanceArtifact) HasLegitParameters() bool {
	if in.Invocation.Parameters == nil {
		return false
	}

	event := &vcs.Event{}
	output := mapstructure.Decode(in.Invocation.Parameters, event.Inputs)

	return output != nil
}
