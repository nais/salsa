package intoto

import (
	"github.com/nais/salsa/pkg/digest"
	"os"
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

	return pa.WithBuildInvocationId().
		WithBuilderRepoDigest().
		WithBuilderId().
		WithBuilderInvocation()
}

func (in *ProvenanceArtifact) repoUri() string {
	return "https://github.com/" + in.Environment.GitHubContext.Repository
}

func (in *ProvenanceArtifact) WithBuildInvocationId() *ProvenanceArtifact {
	in.BuildInvocationId = in.repoUri() + "/actions/runs/" + in.Environment.GitHubContext.RunId
	return in
}

func (in *ProvenanceArtifact) WithBuilderRepoDigest() *ProvenanceArtifact {
	in.BuilderRepoDigest = &slsa.ProvenanceMaterial{
		URI: "git+" + in.repoUri(),
		Digest: slsa.DigestSet{
			digest.AlgorithmSHA1: in.Environment.GitHubContext.SHA,
		},
	}
	return in
}

func (in *ProvenanceArtifact) WithBuilderId() *ProvenanceArtifact {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		in.BuilderId = in.repoUri() + vcs.GitHubHostedIdSuffix
	} else {
		in.BuilderId = in.repoUri() + vcs.GitHubHostedIdSuffix
	}
	return in
}

func (in *ProvenanceArtifact) WithBuilderInvocation() *ProvenanceArtifact {
	in.Invocation = slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI: "git+" + in.repoUri(),
			Digest: slsa.DigestSet{
				digest.AlgorithmSHA1: in.Environment.GitHubContext.SHA,
			},
			EntryPoint: in.Environment.Workflow,
		},
		Parameters: in.Environment.Inputs,
		// Should contain the architecture of the runner.
		Environment: nil,
	}
	return in
}
