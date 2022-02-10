package intoto

import (
	"github.com/nais/salsa/pkg/digest"
	"os"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/vcs"
)

const (
	DefaultBuildId = "https://github.com/nais/salsa"
)

type ProvenanceArtifact struct {
	Name              string
	BuilderId         string
	BuildType         string
	Dependencies      *build.ArtifactDependencies
	BuildStartedOn    time.Time
	BuildInvocationId string
	BuilderRepoDigest *slsa.ProvenanceMaterial
	Invocation        slsa.ProvenanceInvocation
	BuildConfig       string
}

func CreateProvenanceArtifact(name string, deps *build.ArtifactDependencies, context *vcs.AnyContext) *ProvenanceArtifact {
	if context == nil {
		return &ProvenanceArtifact{
			Name:           name,
			BuildType:      vcs.AdHocBuildType,
			BuildConfig:    "Some commands that made this build",
			BuilderId:      DefaultBuildId,
			Dependencies:   deps,
			BuildStartedOn: time.Now().UTC(),
		}
	}

	repoURI := "https://github.com/" + context.GitHubContext.Repository
	pa := &ProvenanceArtifact{
		Name:           name,
		BuildType:      vcs.BuildType,
		Dependencies:   deps,
		BuildStartedOn: time.Now().UTC(),
	}

	return pa.WithBuildInvocationId(repoURI, context).
		WithBuilderRepoDigest(repoURI, context).
		WithBuilderId(repoURI).
		WithBuilderInvocation(repoURI, context)
}

func (in *ProvenanceArtifact) WithBuildInvocationId(repoURI string, context *vcs.AnyContext) *ProvenanceArtifact {
	in.BuildInvocationId = repoURI + "/actions/runs/" + context.GitHubContext.RunId
	return in
}

func (in *ProvenanceArtifact) WithBuilderRepoDigest(repoURI string, context *vcs.AnyContext) *ProvenanceArtifact {
	in.BuilderRepoDigest = &slsa.ProvenanceMaterial{
		URI: "git+" + repoURI,
		Digest: slsa.DigestSet{
			digest.AlgorithmSHA1: context.GitHubContext.SHA,
		},
	}
	return in
}

func (in *ProvenanceArtifact) WithBuilderId(repoURI string) *ProvenanceArtifact {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		in.BuilderId = repoURI + vcs.GitHubHostedIdSuffix
	} else {
		in.BuilderId = repoURI + vcs.GitHubHostedIdSuffix
	}
	return in
}

func (in *ProvenanceArtifact) WithBuilderInvocation(repoURI string, context *vcs.AnyContext) *ProvenanceArtifact {
	in.Invocation = slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI: "git+" + repoURI,
			Digest: slsa.DigestSet{
				digest.AlgorithmSHA1: context.GitHubContext.SHA,
			},
			EntryPoint: context.Workflow,
		},
		Parameters: context.Inputs,
		// Should contain the architecture of the runner.
		Environment: nil,
	}
	return in
}
