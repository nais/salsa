package intoto

import (
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/config"
	"github.com/nais/salsa/pkg/github"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
)

type ProvenanceOptions struct {
	BuildConfig       *BuildConfig
	BuilderId         string
	BuilderRepoDigest *slsa.ProvenanceMaterial
	BuildInvocationId string
	BuildStartedOn    time.Time
	BuildType         string
	Dependencies      *build.ArtifactDependencies
	Invocation        *slsa.ProvenanceInvocation
	Name              string
}

func CreateProvenanceOptions(scanCfg *config.ScanConfiguration) *ProvenanceOptions {
	opts := &ProvenanceOptions{
		BuildStartedOn: time.Now().UTC(),
		Dependencies:   scanCfg.Dependencies,
		Name:           scanCfg.RepoName,
	}

	if scanCfg.CiEnvironment != nil {
		opts.BuildType = github.BuildType
		opts.BuildInvocationId = scanCfg.CiEnvironment.BuildInvocationId()
		opts.BuilderId = scanCfg.CiEnvironment.BuilderId()
		opts.withBuilderRepoDigest(scanCfg.CiEnvironment).withBuilderInvocation(scanCfg.CiEnvironment)
		return opts
	}

	opts.BuildConfig = GenerateBuildConfig(scanCfg)
	opts.BuilderId = github.DefaultBuildId
	opts.BuildType = github.AdHocBuildType
	opts.Invocation = nil
	return opts
}

func (in *ProvenanceOptions) withBuilderRepoDigest(env *github.Environment) *ProvenanceOptions {
	in.BuilderRepoDigest = &slsa.ProvenanceMaterial{
		URI: "git+" + env.RepoUri(),
		Digest: slsa.DigestSet{
			build.AlgorithmSHA1: env.GithubSha(),
		},
	}
	return in
}

func (in *ProvenanceOptions) withBuilderInvocation(env *github.Environment) *ProvenanceOptions {
	in.Invocation = &slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI: "git+" + env.RepoUri(),
			Digest: slsa.DigestSet{
				build.AlgorithmSHA1: env.GithubSha(),
			},
			EntryPoint: env.GitHubContext.Workflow,
		},
		Parameters:  env.AddUserDefinedParameters(),
		Environment: env.NonReproducibleMetadata(),
	}
	return in
}

func (in *ProvenanceOptions) HasBuilderRepoDigest() bool {
	if in.BuilderRepoDigest == nil {
		return false
	}

	return in.BuilderRepoDigest.Digest != nil && in.BuilderRepoDigest.URI != ""

}

func (in *ProvenanceOptions) HasDependencies() bool {
	if in.Dependencies == nil {
		return false
	}

	return len(in.Dependencies.RuntimeDeps) > 0
}

func (in *ProvenanceOptions) HasParameters() bool {
	if in.Invocation == nil {
		return false
	}

	if in.Invocation.Parameters == nil {
		return false
	}

	if in.Invocation.Parameters.(*github.Event) == nil {
		return false
	}

	return in.Invocation.Parameters.(*github.Event).Inputs != nil
}

func (in *ProvenanceOptions) HasEnvironment() bool {
	if in.Invocation == nil {
		return false
	}

	return in.Invocation.Environment != nil
}
