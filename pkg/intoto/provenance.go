package intoto

import (
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
)

func GenerateSlsaPredicate(opts *ProvenanceOptions) *slsa.ProvenancePredicate {
	predicate := &slsa.ProvenancePredicate{
		Builder: slsa.ProvenanceBuilder{
			ID: opts.BuilderId,
		},
		BuildType:   opts.BuildType,
		BuildConfig: opts.BuildConfig,
		Metadata:    withMetadata(opts),
		Materials:   withMaterials(opts),
	}

	if opts.Invocation != nil {
		predicate.Invocation = *opts.Invocation
		return predicate
	}

	return predicate
}

func withMetadata(opts *ProvenanceOptions) *slsa.ProvenanceMetadata {
	timeFinished := opts.GetBuildFinishedOn()
	return &slsa.ProvenanceMetadata{
		BuildInvocationID: opts.BuildInvocationId,
		BuildStartedOn:    &opts.BuildStartedOn,
		BuildFinishedOn:   &timeFinished,
		Completeness:      withCompleteness(opts),
		Reproducible:      opts.Reproducible(),
	}
}

func withCompleteness(opts *ProvenanceOptions) slsa.ProvenanceComplete {
	return slsa.ProvenanceComplete{
		Environment: opts.Environment(),
		Materials:   opts.Materials(),
		Parameters:  opts.Parameters(),
	}
}

func withMaterials(opts *ProvenanceOptions) []slsa.ProvenanceMaterial {
	materials := make([]slsa.ProvenanceMaterial, 0)
	AppendRuntimeDependencies(opts, &materials)
	AppendBuildContext(opts, &materials)
	return materials
}

func AppendRuntimeDependencies(opts *ProvenanceOptions, materials *[]slsa.ProvenanceMaterial) {
	if opts.Dependencies == nil {
		return
	}

	for _, dep := range opts.Dependencies.RuntimeDeps {
		m := slsa.ProvenanceMaterial{
			URI:    dep.ToUri(),
			Digest: dep.ToDigestSet(),
		}
		*materials = append(*materials, m)
	}
}

func AppendBuildContext(opts *ProvenanceOptions, materials *[]slsa.ProvenanceMaterial) {
	if opts.BuilderRepoDigest != nil {
		*materials = append(*materials, *opts.BuilderRepoDigest)
	}
}
