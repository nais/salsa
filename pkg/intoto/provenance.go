package intoto

const (
	SlsaPredicateType = "https://slsa.dev/provenance/v0.2"
	StatementType     = "https://in-toto.io/Statement/v0.1"
)

type DigestSet map[string]string

type Statement struct {
	Type          string    `json:"_type"`
	Subject       []Subject `json:"subject"`
	PredicateType string    `json:"predicateType"`
	Predicate     Predicate `json:"predicate"`
}

type Subject struct {
	Name   string    `json:"name"`
	Digest DigestSet `json:"digest"`
}

func GenerateStatement(app App) *Statement {
	statement := &Statement{
		PredicateType: SlsaPredicateType,
		Subject:       nil,
		Type:          StatementType,
	}
	return statement.withPredicate(app)
}

func (p *Statement) withPredicate(app App) *Statement {
	p.Predicate = Predicate{
		Builder: Builder{
			Id: app.BuilderId,
		},
		BuildType: app.BuildType,
	}.withMetadata("", "", false).
		withRecipe("", "", "").
		withMaterials(app)
	return p
}
