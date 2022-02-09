package intoto

import (
	"github.com/nais/salsa/pkg/vcs"
	"os"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/scan"
)

type ProvenanceArtifact struct {
	Name              string
	BuilderId         string
	BuildType         string
	Dependencies      *scan.ArtifactDependencies
	BuildStartedOn    time.Time
	BuildInvocationId string
	BuilderRepoDigest *slsa.ProvenanceMaterial
	Invocation        slsa.ProvenanceInvocation
}

func CreateProvenanceArtifact(name string, deps *scan.ArtifactDependencies) *ProvenanceArtifact {
	return &ProvenanceArtifact{
		Name:           name,
		BuildType:      "todoType",
		Dependencies:   deps,
		BuildStartedOn: time.Now().UTC(),
	}
}

func (in *ProvenanceArtifact) GenerateSlsaPredicate() slsa.ProvenancePredicate {
	return in.withPredicate()
}

func (in *ProvenanceArtifact) withPredicate() slsa.ProvenancePredicate {
	return slsa.ProvenancePredicate{
		Builder: slsa.ProvenanceBuilder{
			ID: in.BuilderId,
		},
		BuildType:   in.BuildType,
		Invocation:  in.Invocation,
		BuildConfig: nil,
		Metadata:    in.withMetadata(false, time.Now().UTC()),
		Materials:   in.withMaterials(),
	}
}

// TODO: use other type of materials aswell, e.g. github actions run in the build
func (in *ProvenanceArtifact) withMaterials() []slsa.ProvenanceMaterial {
	materials := make([]slsa.ProvenanceMaterial, 0)
	in.AppendRuntimeDependencies(&materials)
	in.AppendBuildContext(&materials)
	return materials
}

func (in *ProvenanceArtifact) AppendRuntimeDependencies(materials *[]slsa.ProvenanceMaterial) {
	for _, dep := range in.Dependencies.RuntimeDeps {
		m := slsa.ProvenanceMaterial{
			URI:    dep.ToUri(),
			Digest: dep.ToDigestSet(),
		}
		*materials = append(*materials, m)
	}
}

func (in *ProvenanceArtifact) AppendBuildContext(materials *[]slsa.ProvenanceMaterial) {
	if in.BuilderRepoDigest != nil {
		*materials = append(*materials, *in.BuilderRepoDigest)
	}
}

func (in *ProvenanceArtifact) withMetadata(rp bool, buildFinished time.Time) *slsa.ProvenanceMetadata {
	return &slsa.ProvenanceMetadata{
		BuildInvocationID: in.BuildInvocationId,
		BuildStartedOn:    &in.BuildStartedOn,
		BuildFinishedOn:   &buildFinished,
		Completeness:      withCompleteness(false, false),
		Reproducible:      rp,
	}
}

func withCompleteness(environment, materials bool) slsa.ProvenanceComplete {
	return slsa.ProvenanceComplete{
		Environment: environment,
		Materials:   materials,
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
		WithBuildType(context)
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

func (in *ProvenanceArtifact) WithBuildType(context *vcs.AnyContext) *ProvenanceArtifact {
	in.Invocation = slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI:        vcs.BuildType,
			Digest:     nil,
			EntryPoint: context.Workflow,
		},
		Parameters:  nil,
		Environment: nil,
	}
	return in
}
