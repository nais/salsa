package intoto

import (
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

	return pa.WithBuilderRepoDigest().WithBuilderInvocation()
}

func (in *ProvenanceArtifact) WithBuilderRepoDigest() *ProvenanceArtifact {
	in.BuilderRepoDigest = &slsa.ProvenanceMaterial{
		URI: "git+" + in.Environment.RepoUri(),
		Digest: slsa.DigestSet{
			digest.AlgorithmSHA1: in.Environment.GithubSha(),
		},
	}
	return in
}

func (in *ProvenanceArtifact) WithBuilderInvocation() *ProvenanceArtifact {
	in.Invocation = slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI: "git+" + in.Environment.RepoUri(),
			Digest: slsa.DigestSet{
				digest.AlgorithmSHA1: in.Environment.GithubSha(),
			},
			EntryPoint: in.Environment.Workflow,
		},
		Parameters: in.Environment.Inputs,
		// Should contain the architecture of the runner.
		Environment: nil,
	}
	return in
}
