package intoto

import (
	"github.com/in-toto/in-toto-golang/in_toto"
	"strings"
	"time"
)

const (
	SlsaPredicateType = "https://slsa.dev/provenance/v0.2"
	StatementType     = "https://in-toto.io/Statement/v0.1"
)

func GenerateStatement(app App) in_toto.Statement {
	return in_toto.Statement{
		StatementHeader: in_toto.StatementHeader{
			Type:          StatementType,
			PredicateType: SlsaPredicateType,
			Subject:       nil,
		},
		Predicate: withPredicate(app),
	}
}

func withPredicate(app App) in_toto.ProvenancePredicate {
	return in_toto.ProvenancePredicate{
		Builder: in_toto.ProvenanceBuilder{
			ID: "",
		},
		Recipe:    withRecipe(),
		Metadata:  withMetadata(false, time.Now(), time.Now()),
		Materials: withMaterials(app),
	}
}

func FindMaterials(predicate in_toto.ProvenancePredicate, value string) []in_toto.ProvenanceMaterial {
	var found []in_toto.ProvenanceMaterial
	for _, m := range predicate.Materials {
		if find(m, value) {
			found = append(found, m)
		}
	}
	return found
}

func find(material in_toto.ProvenanceMaterial, value string) bool {
	return strings.Contains(material.URI, value)
}

func withMetadata(rp bool, buildStarted, buildFinished time.Time) *in_toto.ProvenanceMetadata {
	return &in_toto.ProvenanceMetadata{
		BuildStartedOn:  &buildStarted,
		BuildFinishedOn: &buildFinished,
		Completeness:    withCompleteness(false, false, false),
		Reproducible:    rp,
	}
}

func withCompleteness(arguments, environment, materials bool) in_toto.ProvenanceComplete {
	return in_toto.ProvenanceComplete{
		Arguments:   arguments,
		Environment: environment,
		Materials:   materials,
	}
}

func withRecipe() in_toto.ProvenanceRecipe {
	return in_toto.ProvenanceRecipe{
		Type:              "",
		DefinedInMaterial: nil,
		EntryPoint:        "",
		Arguments:         nil,
		Environment:       nil,
	}
}

// TODO: use other type of materials aswell, e.g. github actions run in the build
func withMaterials(app App) []in_toto.ProvenanceMaterial {
	materials := make([]in_toto.ProvenanceMaterial, 0)
	for k, v := range app.Dependencies {
		m := in_toto.ProvenanceMaterial{
			URI:    k + ":" + v,
			Digest: nil,
		}
		materials = append(materials, m)
	}
	return materials
}
