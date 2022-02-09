package intoto

import (
	"time"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
)

func GenerateSlsaPredicate(pa *ProvenanceArtifact) slsa.ProvenancePredicate {
	return withPredicate(pa)
}

func withPredicate(pa *ProvenanceArtifact) slsa.ProvenancePredicate {
	return slsa.ProvenancePredicate{
		Builder: slsa.ProvenanceBuilder{
			ID: pa.BuilderId,
		},
		BuildType:   pa.BuildType,
		Invocation:  pa.Invocation,
		BuildConfig: nil,
		Metadata:    withMetadata(pa, false, time.Now().UTC()),
		Materials:   pa.withMaterials(pa),
	}
}

// TODO: use other type of materials aswell, e.g. github actions run in the build
func (in *ProvenanceArtifact) withMaterials(pa *ProvenanceArtifact) []slsa.ProvenanceMaterial {
	materials := make([]slsa.ProvenanceMaterial, 0)
	AppendRuntimeDependencies(pa, &materials)
	AppendBuildContext(pa, &materials)
	return materials
}

func AppendRuntimeDependencies(pa *ProvenanceArtifact, materials *[]slsa.ProvenanceMaterial) {
	for _, dep := range pa.Dependencies.RuntimeDeps {
		m := slsa.ProvenanceMaterial{
			URI:    dep.ToUri(),
			Digest: dep.ToDigestSet(),
		}
		*materials = append(*materials, m)
	}
}

func AppendBuildContext(pa *ProvenanceArtifact, materials *[]slsa.ProvenanceMaterial) {
	if pa.BuilderRepoDigest != nil {
		*materials = append(*materials, *pa.BuilderRepoDigest)
	}
}

func withMetadata(pa *ProvenanceArtifact, rp bool, buildFinished time.Time) *slsa.ProvenanceMetadata {
	return &slsa.ProvenanceMetadata{
		BuildInvocationID: pa.BuildInvocationId,
		BuildStartedOn:    &pa.BuildStartedOn,
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
