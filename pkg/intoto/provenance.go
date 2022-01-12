package intoto

import (
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	"github.com/nais/salsa/pkg/digest"
	"strings"
	"time"
)

func GenerateSlsaPredicate(app App) slsa.ProvenancePredicate {
	return withPredicate(app)
}

func withPredicate(app App) slsa.ProvenancePredicate {
	return slsa.ProvenancePredicate{
		Builder:   slsa.ProvenanceBuilder{},
		BuildType: "yolo",
		Invocation: slsa.ProvenanceInvocation{
			ConfigSource: slsa.ConfigSource{},
			Parameters:   nil,
			Environment:  nil,
		},
		BuildConfig: nil,
		Metadata:    withMetadata(false, time.Now(), time.Now()),
		Materials:   withMaterials(app),
	}
}

func FindMaterials(materials []slsa.ProvenanceMaterial, value string) []slsa.ProvenanceMaterial {
	f := make([]slsa.ProvenanceMaterial, 0)
	for _, m := range materials {
		if find(m, value) {
			f = append(f, m)
		}
	}
	return f
}

func find(material slsa.ProvenanceMaterial, value string) bool {
	return strings.Contains(material.URI, value)
}

func withMetadata(rp bool, buildStarted, buildFinished time.Time) *slsa.ProvenanceMetadata {
	return &slsa.ProvenanceMetadata{
		BuildStartedOn:  &buildStarted,
		BuildFinishedOn: &buildFinished,
		Completeness:    withCompleteness(false, false),
		Reproducible:    rp,
	}
}

func withCompleteness(environment, materials bool) slsa.ProvenanceComplete {
	return slsa.ProvenanceComplete{
		Environment: environment,
		Materials:   materials,
	}
}

// TODO: use other type of materials aswell, e.g. github actions run in the build
func withMaterials(app App) []slsa.ProvenanceMaterial {
	materials := make([]slsa.ProvenanceMaterial, 0)
	for k, v := range app.Dependencies {
		var uri = k + ":" + v
		var hashedDigest = digest.Hash(uri)
		m := slsa.ProvenanceMaterial{
			URI:    uri,
			Digest: slsa.DigestSet{hashedDigest.Algorithm().String(): hashedDigest.Encoded()},
		}
		materials = append(materials, m)
	}
	return materials
}
