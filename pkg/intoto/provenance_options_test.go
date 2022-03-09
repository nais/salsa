package intoto

import (
	"fmt"
	"github.com/nais/salsa/pkg/config"
	"github.com/nais/salsa/pkg/digest"
	"github.com/spf13/cobra"
	"testing"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/vcs"
	"github.com/stretchr/testify/assert"
)

func TestCreateProvenanceOptions(t *testing.T) {
	deps := ExpectedDeps()
	artDeps := ExpectedArtDeps(deps)

	for _, test := range []struct {
		name              string
		buildType         string
		buildInvocationId string
		builderId         string
		buildConfig       *BuildConfig
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
			buildConfig:       buildConfig(),
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
			buildConfig:       nil,
			builderRepoDigest: ExpectedBuilderRepoDigestMaterial(),
			configSource:      ExpectedConfigSource(),
			buildTimerIsSet:   true,
			runnerContext:     true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.runnerContext {
				env := Environment()
				scanCfg := &config.ScanConfiguration{
					WorkDir:       "",
					RepoName:      "artifact",
					Dependencies:  artDeps,
					CiEnvironment: env,
					Cmd:           nil,
				}
				provenanceArtifact := CreateProvenanceOptions(scanCfg)
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

				scanCfg := &config.ScanConfiguration{
					WorkDir:       "",
					RepoName:      "artifact",
					Dependencies:  artDeps,
					CiEnvironment: nil,
					Cmd:           &cobra.Command{Use: "salsa"},
				}

				provenanceArtifact := CreateProvenanceOptions(scanCfg)
				assert.Equal(t, "artifact", provenanceArtifact.Name)
				assert.Equal(t, test.buildType, provenanceArtifact.BuildType)
				assert.Equal(t, deps, provenanceArtifact.Dependencies.RuntimeDeps)
				assert.Equal(t, test.buildTimerIsSet, time.Now().UTC().After(provenanceArtifact.BuildStartedOn))
				assert.Equal(t, test.buildInvocationId, provenanceArtifact.BuildInvocationId)
				assert.Equal(t, test.buildConfig, provenanceArtifact.BuildConfig)
				assert.Empty(t, provenanceArtifact.Invocation)
				assert.Equal(t, test.builderId, provenanceArtifact.BuilderId)
				assert.Equal(t, test.builderRepoDigest, provenanceArtifact.BuilderRepoDigest)
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

func ExpectedDeps() map[string]build.Dependency {
	deps := map[string]build.Dependency{}
	checksum := build.Verification("todo", "todo")
	deps[fmt.Sprintf("%s:%s", "groupId", "artifactId")] = build.Dependence(
		fmt.Sprintf("%s:%s", "groupId", "artifactId"),
		"v01",
		checksum,
	)
	return deps
}

func ExpectedArtDeps(deps map[string]build.Dependency) *build.ArtifactDependencies {
	return &build.ArtifactDependencies{
		Cmd: build.Cmd{
			Path:     "lang",
			CmdFlags: "list:deps",
		},
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
			EventName:  "workflow_dispatch",
		},
		Event: &vcs.Event{
			Inputs: []byte("some user inputs"),
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

func buildConfig() *BuildConfig {
	return &BuildConfig{
		Commands: []string{
			"salsa ",
			"lang list:deps",
		},
		Shell: "bash",
	}
}
