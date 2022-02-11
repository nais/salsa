package intoto

import (
	"fmt"
	"github.com/nais/salsa/pkg/digest"
	"testing"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/vcs"
	"github.com/stretchr/testify/assert"
)

func TestCreateProvenanceArtifact(t *testing.T) {
	deps := ExpectedDeps()
	artDeps := ExpectedArtDeps(deps)

	for _, test := range []struct {
		name              string
		buildType         string
		buildInvocationId string
		builderId         string
		buildConfig       string
		builderRepoDigest *slsa.ProvenanceMaterial
		configSource      slsa.ConfigSource
		buildTimerIsSet   bool
		runnerContext     bool
	}{
		{
			name:              "create provenance artifact with default values",
			buildType:         vcs.AdHocBuildType,
			buildInvocationId: "",
			builderId:         vcs.DefaultBuildId,
			buildConfig:       "Some commands that made this build",
			builderRepoDigest: (*slsa.ProvenanceMaterial)(nil),
			configSource: slsa.ConfigSource{
				URI:        "",
				Digest:     slsa.DigestSet(nil),
				EntryPoint: "",
			},
			buildTimerIsSet: true,
			runnerContext:   false,
		},
		{
			name:              "create provenance artifact with runner context",
			buildType:         vcs.BuildType,
			buildInvocationId: "https://github.com/nais/salsa/actions/runs/1234",
			builderId:         "https://github.com/nais/salsa/Attestations/GitHubHostedActions@v1",
			buildConfig:       "",
			builderRepoDigest: ExpectedBuilderRepoDigestMaterial(),
			configSource:      ExpectedConfigSource(),
			buildTimerIsSet:   true,
			runnerContext:     true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.runnerContext {
				env := Environment()
				provenanceArtifact := CreateProvenanceArtifact("artifact", artDeps, env)
				assert.Equal(t, "artifact", provenanceArtifact.Name)
				assert.Equal(t, test.buildType, provenanceArtifact.BuildType)
				assert.Equal(t, deps, provenanceArtifact.Dependencies.RuntimeDeps)
				assert.Equal(t, test.buildTimerIsSet, time.Now().UTC().After(provenanceArtifact.BuildStartedOn))
				assert.Equal(t, test.buildInvocationId, provenanceArtifact.BuildInvocationId)
				assert.Equal(t, test.buildConfig, provenanceArtifact.BuildConfig)
				assert.NotEmpty(t, provenanceArtifact.Invocation)
				assert.NotEmpty(t, provenanceArtifact.Invocation.Parameters)
				assert.NotEmpty(t, provenanceArtifact.Invocation.Environment)
				assert.Equal(t, test.builderId, provenanceArtifact.BuilderId)
				assert.Equal(t, test.builderRepoDigest, provenanceArtifact.BuilderRepoDigest)
				assert.Equal(t, test.configSource, provenanceArtifact.Invocation.ConfigSource)

			} else {

				provenanceArtifact := CreateProvenanceArtifact("artifact", artDeps, nil)
				assert.Equal(t, "artifact", provenanceArtifact.Name)
				assert.Equal(t, test.buildType, provenanceArtifact.BuildType)
				assert.Equal(t, deps, provenanceArtifact.Dependencies.RuntimeDeps)
				assert.Equal(t, test.buildTimerIsSet, time.Now().UTC().After(provenanceArtifact.BuildStartedOn))
				assert.Equal(t, test.buildInvocationId, provenanceArtifact.BuildInvocationId)
				assert.Equal(t, test.buildConfig, provenanceArtifact.BuildConfig)
				assert.Equal(t, slsa.ProvenanceInvocation{}, provenanceArtifact.Invocation)
				assert.Empty(t, provenanceArtifact.Invocation.Parameters)
				assert.Empty(t, provenanceArtifact.Invocation.Environment)
				assert.Equal(t, test.builderId, provenanceArtifact.BuilderId)
				assert.Equal(t, test.builderRepoDigest, provenanceArtifact.BuilderRepoDigest)
				assert.Equal(t, test.configSource, provenanceArtifact.Invocation.ConfigSource)
			}
		})
	}
}

func ExpectedBuilderRepoDigestMaterial() *slsa.ProvenanceMaterial {
	return &slsa.ProvenanceMaterial{
		URI: "git+https://github.com/nais/salsa",
		Digest: slsa.DigestSet{
			digest.AlgorithmSHA1: "4321",
		},
	}
}

func ExpectedDeps() []build.Dependency {
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

func ExpectedArtDeps(deps []build.Dependency) *build.ArtifactDependencies {
	return &build.ArtifactDependencies{
		Cmd:         "lang list:deps",
		RuntimeDeps: deps,
	}
}

func Environment() *vcs.Environment {
	return &vcs.Environment{
		GitHubContext: vcs.GitHubContext{
			Repository: "nais/salsa",
			RunId:      "1234",
			SHA:        "4321",
			Workflow:   "Create a provenance",
			ServerUrl:  "https://github.com",
		},
		Event: vcs.Event{
			Inputs: []byte("some vents"),
		},
		RunnerContext: vcs.RunnerContext{
			OS:        "Linux",
			Temp:      "/home/runner/work/_temp",
			ToolCache: "/opt/hostedtoolcache",
		},
	}
}

func ExpectedConfigSource() slsa.ConfigSource {
	return slsa.ConfigSource{
		URI: "git+https://github.com/nais/salsa",
		Digest: slsa.DigestSet{
			digest.AlgorithmSHA1: "4321",
		},
		EntryPoint: "Create a provenance",
	}
}
