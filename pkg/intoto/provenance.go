package intoto

const (
	// SlsaPredicateType the predicate type for SLSA intoto statements
	SlsaPredicateType = "https://slsa.dev/provenance/v0.2"
	// StatementType the type of the intoto statement
	StatementType = "https://in-toto.io/Statement/v0.1"
)

type DigestSet map[string]string

type Provenance struct {
	Type          string    `json:"_type"`
	Subject       []Subject `json:"subject"`
	PredicateType string    `json:"predicateType"`
	Predicate     Predicate `json:"predicate"`
}

type Subject struct {
	Name   string    `json:"name"`
	Digest DigestSet `json:"digest"`
}

func (p *Provenance) withPredicate(dependencies map[string]string, builderId string) *Provenance {
	p.Predicate = Predicate{}.
		withBuilder(builderId).
		withBuildType("https://").
		withMetadata("", "", false).
		withRecipe("", "", "").
		withMaterials(dependencies)
	return p
}

func GenerateStatement(dependencies map[string]string, builderId string) *Provenance {
	statement := &Provenance{
		PredicateType: SlsaPredicateType,
		Subject:       nil,
		Type:          StatementType,
	}
	return statement.withPredicate(dependencies, builderId)
}
