package intoto

import (
	"os"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/scan/common"
	"github.com/nais/salsa/pkg/vcs"
)

type ProvenanceArtifact struct {
	Name              string
	BuilderId         string
	BuildType         string
	Dependencies      *common.ArtifactDependencies
	BuildStartedOn    time.Time
	BuildInvocationId string
	BuilderRepoDigest *slsa.ProvenanceMaterial
	Invocation        slsa.ProvenanceInvocation
}

func CreateProvenanceArtifact(name string, deps *common.ArtifactDependencies) *ProvenanceArtifact {
	return &ProvenanceArtifact{
		Name:           name,
		BuildType:      "todoType",
		Dependencies:   deps,
		BuildStartedOn: time.Now().UTC(),
	}
}

func (in *ProvenanceArtifact) WithRunnerContext(context *vcs.AnyContext) *ProvenanceArtifact {
	if context == nil {
		// Required
		in.BuilderId = "default"
		return in
	}

	repoURI := "https://github.com/" + context.GitHubContext.Repository
	return in.WithBuildInvocationId(repoURI, context).
		WithBuilderRepoDigest(repoURI, context).
		WithBuilderId(repoURI).
		WithBuilderInvocation(context)
}

func (in *ProvenanceArtifact) WithBuildInvocationId(repoURI string, context *vcs.AnyContext) *ProvenanceArtifact {
	in.BuildInvocationId = repoURI + "/actions/runs/" + context.GitHubContext.RunId
	return in
}

func (in *ProvenanceArtifact) WithBuilderRepoDigest(repoURI string, context *vcs.AnyContext) *ProvenanceArtifact {
	in.BuilderRepoDigest = &slsa.ProvenanceMaterial{
		URI:    "git+" + repoURI,
		Digest: slsa.DigestSet{"sha1": context.GitHubContext.SHA},
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

func (in *ProvenanceArtifact) WithBuilderInvocation(context *vcs.AnyContext) *ProvenanceArtifact {
	in.Invocation = slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI:        vcs.BuildType,
			Digest:     nil,
			EntryPoint: context.Workflow,
		},
		Parameters:  context.Inputs,
		Environment: nil,
	}
	return in
}
