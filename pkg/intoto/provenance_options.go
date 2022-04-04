package intoto

import (
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/config"
	"github.com/nais/salsa/pkg/vcs"
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
)

const (
	// AdHocBuildType no entry point, and the commands were run in an ad-hoc fashion
	AdHocBuildType = "https://github.com/nais/salsa/ManuallyRunCommands@v1"
	DefaultBuildId = "https://github.com/nais/salsa"
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
		BuilderId:    DefaultBuildId,
		BuildType:    AdHocBuildType,
		Dependencies: scanCfg.Dependencies,
		Name:         scanCfg.RepoName,
	}

	context := scanCfg.ContextEnvironment
	if context != nil {
		opts.BuildStartedOn = context.GetBuildStartedOn()
		opts.BuildType = context.BuildType()
		opts.BuildInvocationId = context.BuildInvocationId()
		opts.BuilderId = context.BuilderId()
		opts.withBuilderRepoDigest(context).
			withBuilderInvocation(context)
		return opts
	}

	opts.BuildConfig = GenerateBuildConfig(scanCfg)
	opts.BuildStartedOn = time.Now().UTC()
	opts.Invocation = nil
	return opts
}

func (in *ProvenanceOptions) withBuilderRepoDigest(env vcs.ContextEnvironment) *ProvenanceOptions {
	in.BuilderRepoDigest = &slsa.ProvenanceMaterial{
		URI: "git+" + env.RepoUri(),
		Digest: slsa.DigestSet{
			build.AlgorithmSHA1: env.Sha(),
		},
	}
	return in
}

func (in *ProvenanceOptions) withBuilderInvocation(env vcs.ContextEnvironment) *ProvenanceOptions {
	in.Invocation = &slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI: "git+" + env.RepoUri(),
			Digest: slsa.DigestSet{
				build.AlgorithmSHA1: env.Sha(),
			},
			EntryPoint: env.Context(),
		},
		Parameters:  env.UserDefinedParameters(),
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

func (in *ProvenanceOptions) Parameters() bool {
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

func (in *ProvenanceOptions) Environment() bool {
	if in.Invocation == nil {
		return false
	}

	return in.Invocation.Environment != nil
}

func (in *ProvenanceOptions) Materials() bool {
	return in.HasDependencies() && in.HasBuilderRepoDigest()
}

func (in *ProvenanceOptions) Reproducible() bool {
	return in.Environment() && in.Materials() && in.Parameters()
}
