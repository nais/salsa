package intoto

import (
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/config"
	"github.com/nais/salsa/pkg/vcs"
	"github.com/spf13/cobra"
	"os"
	"testing"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSlsaPredicate(t *testing.T) {
	deps := ExpectedDeps()
	artDeps := ExpectedArtDeps(deps)
	for _, test := range []struct {
		name                    string
		buildType               string
		buildInvocationId       string
		builderId               string
		buildConfig             *BuildConfig
		materials               []slsa.ProvenanceMaterial
		configSource            slsa.ConfigSource
		buildTimerIsSet         bool
		buildTimerFinishedIsSet bool
		runnerContext           bool
	}{
		{
			name:              "create slsa provenance artifact with default values",
			buildType:         "https://github.com/nais/salsa/ManuallyRunCommands@v1",
			buildInvocationId: "",
			builderId:         "https://github.com/nais/salsa",
			buildConfig:       buildConfig(),
			materials:         ExpectedDependenciesMaterial(),
			configSource: slsa.ConfigSource{
				URI:        "",
				Digest:     slsa.DigestSet(nil),
				EntryPoint: "",
			},
			buildTimerIsSet:         true,
			buildTimerFinishedIsSet: true,
			runnerContext:           false,
		},
		{
			name:                    "create slsa provenance with runner context",
			buildType:               "https://github.com/Attestations/GitHubActionsWorkflow@v1",
			buildInvocationId:       "https://github.com/nais/salsa/actions/runs/1234",
			builderId:               "https://github.com/nais/salsa/Attestations/GitHubHostedActions@v1",
			buildConfig:             nil,
			materials:               ToExpectedMaterials(),
			configSource:            ExpectedConfigSource(),
			buildTimerIsSet:         true,
			buildTimerFinishedIsSet: true,
			runnerContext:           true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.runnerContext {
				env := Environment()

				scanCfg := &config.ScanConfiguration{
					WorkDir:            "",
					RepoName:           "artifact",
					Dependencies:       artDeps,
					ContextEnvironment: env,
					Cmd:                nil,
				}

				opts := CreateProvenanceOptions(scanCfg)
				slsaPredicate := GenerateSlsaPredicate(opts)
				err := os.Setenv("GITHUB_ACTIONS", "true")
				assert.NoError(t, err)

				// VCS Context
				assert.Equal(t, test.buildType, slsaPredicate.BuildType)
				assert.NotEmpty(t, slsaPredicate.Invocation)
				i, err := slsaPredicate.Invocation.Parameters.(*vcs.Event).Inputs.MarshalJSON()
				assert.NoError(t, err)
				assert.Equal(t, "some user inputs", string(i))
				e := slsaPredicate.Invocation.Environment.(*vcs.Metadata)
				assert.NoError(t, err)
				assert.Equal(t, expectedMetadata(), e)
				assert.NotEmpty(t, slsaPredicate.Invocation.Environment)
				assert.Equal(t, test.configSource, slsaPredicate.Invocation.ConfigSource)
				assert.Equal(t, test.builderId, slsaPredicate.Builder.ID)

				// metadata
				assert.Equal(t, test.buildInvocationId, slsaPredicate.Metadata.BuildInvocationID)
				assert.Equal(t, test.buildTimerIsSet, time.Now().UTC().After(*slsaPredicate.Metadata.BuildStartedOn))
				assert.Equal(t, test.buildTimerFinishedIsSet, time.Now().UTC().After(*slsaPredicate.Metadata.BuildFinishedOn))
				assert.Equal(t, true, slsaPredicate.Metadata.Reproducible)

				// completeness
				assert.Equal(t, true, slsaPredicate.Metadata.Completeness.Environment)
				assert.Equal(t, true, slsaPredicate.Metadata.Completeness.Materials)
				assert.Equal(t, true, slsaPredicate.Metadata.Completeness.Parameters)

				// materials
				assert.Equal(t, 2, len(slsaPredicate.Materials))
				assert.Equal(t, test.materials, slsaPredicate.Materials)

			} else {

				scanCfg := &config.ScanConfiguration{
					WorkDir:            "",
					RepoName:           "artifact",
					Dependencies:       artDeps,
					ContextEnvironment: nil,
					Cmd:                &cobra.Command{Use: "salsa"},
				}

				opts := CreateProvenanceOptions(scanCfg)
				slsaPredicate := GenerateSlsaPredicate(opts)

				// Predicate
				assert.Equal(t, test.buildType, slsaPredicate.BuildType)
				assert.Equal(t, test.builderId, slsaPredicate.Builder.ID)
				assert.Equal(t, test.buildConfig, slsaPredicate.BuildConfig)
				assert.Equal(t, test.configSource, slsaPredicate.Invocation.ConfigSource)
				assert.Empty(t, slsaPredicate.Invocation.Parameters)
				assert.Empty(t, slsaPredicate.Invocation.Environment)

				// metadata
				assert.Equal(t, test.buildInvocationId, slsaPredicate.Metadata.BuildInvocationID)
				assert.Equal(t, test.buildTimerIsSet, time.Now().UTC().After(*slsaPredicate.Metadata.BuildStartedOn))
				assert.Equal(t, test.buildTimerFinishedIsSet, time.Now().UTC().After(*slsaPredicate.Metadata.BuildFinishedOn))
				assert.Equal(t, false, slsaPredicate.Metadata.Reproducible)

				// completeness
				assert.Equal(t, false, slsaPredicate.Metadata.Completeness.Environment)
				assert.Equal(t, false, slsaPredicate.Metadata.Completeness.Materials)
				assert.Equal(t, false, slsaPredicate.Metadata.Completeness.Parameters)

				// materials
				assert.Equal(t, 1, len(slsaPredicate.Materials))
				assert.Equal(t, ExpectedDependenciesMaterial(), slsaPredicate.Materials)
			}
		})

	}
}

func expectedMetadata() *vcs.Metadata {
	return &vcs.Metadata{
		Arch: "",
		Env:  map[string]string{},
		Context: vcs.Context{
			Github: vcs.Github{
				RunId: "1234",
			},
			Runner: vcs.Runner{
				Os:   "Linux",
				Temp: "/home/runner/work/_temp"},
		},
	}
}

func ExpectedDependenciesMaterial() []slsa.ProvenanceMaterial {
	return []slsa.ProvenanceMaterial{
		{
			URI: "pkg:groupId:artifactId:v01",
			Digest: slsa.DigestSet{
				"todo": "todo",
			},
		},
	}
}

func ToExpectedMaterials() []slsa.ProvenanceMaterial {
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
				build.AlgorithmSHA1: "4321",
			},
		},
	}
}
