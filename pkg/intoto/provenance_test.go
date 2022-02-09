package intoto

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/vcs"
	"github.com/stretchr/testify/assert"
)

func TestCreateProvenanceArtifact(t *testing.T) {
	deps := toDeps()
	artDeps := toArtDeps(deps)

	provenanceArtifact := CreateProvenanceArtifact("artifact", artDeps)
	assert.Equal(t, "artifact", provenanceArtifact.Name)
	assert.Equal(t, "todoType", provenanceArtifact.BuildType)
	assert.Equal(t, deps, provenanceArtifact.Dependencies.RuntimeDeps)
	assert.True(t, time.Now().UTC().After(provenanceArtifact.BuildStartedOn))
}

func TestCreateProvenanceArtifact_withContext(t *testing.T) {
	deps := toDeps()
	artDeps := toArtDeps(deps)
	anyContext := toAnyContext()
	err := os.Setenv("GITHUB_ACTIONS", "true")
	assert.NoError(t, err)

	provenanceArtifact := CreateProvenanceArtifact("artifact", artDeps).
		WithRunnerContext(anyContext)
	slsaPredicate := GenerateSlsaPredicate(provenanceArtifact)

	// VCS Context
	assert.Equal(t, "https://github.com/nais/salsa/actions/runs/1234", slsaPredicate.Metadata.BuildInvocationID)
	assert.Equal(t, "todoType", slsaPredicate.BuildType)
	assert.Equal(t, toExpectedInvocation(), slsaPredicate.Invocation)
	assert.Equal(t, "https://github.com/nais/salsa/Attestations/GitHubHostedActions@v1", slsaPredicate.Builder.ID)
	assert.Equal(t, toExpectedMaterials(), slsaPredicate.Materials)

	// completeness
	assert.Equal(t, false, slsaPredicate.Metadata.Completeness.Environment)
	assert.Equal(t, true, slsaPredicate.Metadata.Completeness.Materials)
	assert.Equal(t, true, slsaPredicate.Metadata.Completeness.Parameters)
}

func TestProvenanceArtifact_GenerateSlsaPredicate(t *testing.T) {
	deps := toDeps()
	artDeps := toArtDeps(deps)

	provenanceArtifact := CreateProvenanceArtifact("artifact", artDeps)
	slsaPredicate := GenerateSlsaPredicate(provenanceArtifact)

	// Predicate
	expectedConfigSource := slsa.ConfigSource{
		URI:        "",
		Digest:     slsa.DigestSet(nil),
		EntryPoint: "",
	}

	expectedMaterials := []slsa.ProvenanceMaterial{
		{
			URI: "pkg:groupId:artifactId:v01",
			Digest: slsa.DigestSet{
				"todo": "todo",
			},
		},
	}

	assert.Equal(t, "todoType", slsaPredicate.BuildType)
	assert.Equal(t, "", slsaPredicate.Builder.ID)
	assert.Equal(t, nil, slsaPredicate.BuildConfig)
	assert.Equal(t, expectedConfigSource, slsaPredicate.Invocation.ConfigSource)
	assert.Equal(t, nil, slsaPredicate.Invocation.Parameters)
	assert.Equal(t, nil, slsaPredicate.Invocation.Environment)

	// metadata
	assert.Equal(t, "", slsaPredicate.Metadata.BuildInvocationID)
	assert.True(t, time.Now().UTC().After(*slsaPredicate.Metadata.BuildStartedOn))
	assert.True(t, time.Now().UTC().After(*slsaPredicate.Metadata.BuildFinishedOn))
	assert.Equal(t, false, slsaPredicate.Metadata.Reproducible)

	// completeness
	assert.Equal(t, false, slsaPredicate.Metadata.Completeness.Environment)
	assert.Equal(t, false, slsaPredicate.Metadata.Completeness.Materials)
	assert.Equal(t, false, slsaPredicate.Metadata.Completeness.Parameters)

	// materials
	assert.Equal(t, expectedMaterials, slsaPredicate.Materials)
}

func toDeps() []build.Dependency {
	return []build.Dependency{
		{
			Coordinates: fmt.Sprintf("%s:%s", "groupId", "artifactId"),
			Version:     "v01",
			CheckSum: build.CheckSum{
				Algorithm: "todo",
				Digest:    "todo",
			},
		},
	}
}

func toArtDeps(deps []build.Dependency) *build.ArtifactDependencies {
	return &build.ArtifactDependencies{
		Cmd:         "lang list:deps",
		RuntimeDeps: deps,
	}
}

func toAnyContext() *vcs.AnyContext {
	return &vcs.AnyContext{
		GitHubContext: vcs.GitHubContext{
			Repository: "nais/salsa",
			RunId:      "1234",
			SHA:        "4321",
			Workflow:   "Create a provenance",
		},
		AnyEvent: vcs.AnyEvent{
			Inputs: nil,
		},
	}
}

func toExpectedMaterials() []slsa.ProvenanceMaterial {
	return []slsa.ProvenanceMaterial{
		{
			URI: "pkg:groupId:artifactId:v01",
			Digest: slsa.DigestSet{
				"todo": "todo",
			},
		},
		{
			URI: "git+https://github.com/nais/salsa",
			Digest: slsa.DigestSet{
				"sha1": "4321",
			},
		},
	}
}

func toExpectedInvocation() slsa.ProvenanceInvocation {
	return slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI:        "https://github.com/Attestations/GitHubActionsWorkflow@v1",
			Digest:     slsa.DigestSet(nil),
			EntryPoint: "Create a provenance",
		},
		Parameters:  json.RawMessage(nil),
		Environment: interface{}(nil)}
}
