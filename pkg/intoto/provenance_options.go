package intoto

import (
	"github.com/nais/salsa/pkg/config"
	"github.com/nais/salsa/pkg/digest"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/vcs"
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
		opts.BuildType = vcs.BuildType
		opts.BuildInvocationId = scanCfg.CiEnvironment.BuildInvocationId()
		opts.BuilderId = scanCfg.CiEnvironment.BuilderId()
		opts.withBuilderRepoDigest(scanCfg.CiEnvironment).withBuilderInvocation(scanCfg.CiEnvironment)
		return opts
	}

	opts.BuildConfig = GenerateBuildConfig(scanCfg)
	opts.BuilderId = vcs.DefaultBuildId
	opts.BuildType = vcs.AdHocBuildType
	opts.Invocation = nil
	return opts
}

func (in *ProvenanceOptions) withBuilderRepoDigest(env *vcs.Environment) *ProvenanceOptions {
	in.BuilderRepoDigest = &slsa.ProvenanceMaterial{
		URI: "git+" + env.RepoUri(),
		Digest: slsa.DigestSet{
			digest.AlgorithmSHA1: env.GithubSha(),
		},
	}
	return in
}

func (in *ProvenanceOptions) withBuilderInvocation(env *vcs.Environment) *ProvenanceOptions {
	in.Invocation = &slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI: "git+" + env.RepoUri(),
			Digest: slsa.DigestSet{
				digest.AlgorithmSHA1: env.GithubSha(),
			},
			EntryPoint: env.GitHubContext.Workflow,
		},
		Parameters:  env.AddUserDefinedParameters(),
		Environment: NonReproducibleMetadata(env),
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

	if in.Invocation.Parameters.(*vcs.Event) == nil {
		return false
	}

	return in.Invocation.Parameters.(*vcs.Event).Inputs != nil
}

func (in *ProvenanceOptions) HasEnvironment() bool {
	if in.Invocation == nil {
		return false
	}

	return in.Invocation.Environment != nil
}

func NonReproducibleMetadata(env *vcs.Environment) *Metadata {
	// Other variables that are required to reproduce the build and that cannot be
	// recomputed using existing information.
	//(Documentation would explain how to recompute the rest of the fields.)
	return &Metadata{
		Arch: env.RunnerContext.Arch,
		Env:  env.GetCurrentEnvironment(),
		Context: Context{
			Github: Github{
				RunId: env.GitHubContext.RunId,
			},
			Runner: Runner{
				Os:   env.RunnerContext.OS,
				Temp: env.RunnerContext.Temp,
			},
		},
	}
}
