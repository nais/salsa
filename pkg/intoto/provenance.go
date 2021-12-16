package intoto

const (
	SlsaPredicateType = "https://slsa.dev/provenance/v0.2"
	StatementType     = "https://in-toto.io/Statement/v0.1"
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

func GenerateStatement(dependencies map[string]string, builderId string) *Provenance {
	statement := &Provenance{
		PredicateType: SlsaPredicateType,
		Subject:       nil,
		Type:          StatementType,
	}
	return statement.withPredicate(dependencies, builderId)
}

func (p *Provenance) withPredicate(dependencies map[string]string, builderId string) *Provenance {
	p.Predicate = Predicate{
		Builder: Builder{
			Id: builderId,
		},
		BuildType: "uri",
	}.withMetadata("", "", false).
		withRecipe("", "", "").
		withMaterials(dependencies)
	return p
}
