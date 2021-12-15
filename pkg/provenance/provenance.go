package provenance

type Predicate struct {
	Builder   `json:"builder"`
	Metadata  `json:"metadata"`
	Recipe    `json:"recipe"`
	Materials []Material `json:"materials"`
}

type Material struct {
	URI    string    `json:"uri"`
	Digest DigestSet `json:"digest"`
}

type Builder struct {
	Id string `json:"id"`
}
type Metadata struct {
	BuildInvocationId string `json:"buildInvocationId"`
	Completeness      `json:"completeness"`
	Reproducible      bool `json:"reproducible"`
	// BuildStartedOn not defined as it's not available from a GitHub Action.
	BuildFinishedOn string `json:"buildFinishedOn"`
}

type Completeness struct {
	Arguments   bool `json:"arguments"`
	Environment bool `json:"environment"`
	Materials   bool `json:"materials"`
}

type DigestSet map[string]string

type Subject struct {
	Name   string    `json:"name"`
	Digest DigestSet `json:"digest"`
}

type Provenance struct {
	Type          string    `json:"_type"`
	Subject       []Subject `json:"subject"`
	PredicateType string    `json:"predicateType"`
	Predicate     Predicate `json:"predicate"`
}

type Arguments struct {
	CFLAGS string `json:"CFLAGS"`
}

type Recipe struct {
	Type              string    `json:"type"`
	DefinedInMaterial int       `json:"definedInMaterial"`
	EntryPoint        string    `json:"entryPoint"`
	Arguments         Arguments `json:"arguments"`
}

func createPredicate(materials []Material) Predicate {
	return Predicate{
		Builder:   Builder{},
		Metadata:  Metadata{},
		Recipe:    createRecipe(),
		Materials: materials,
	}
}

func createRecipe() Recipe {
	return Recipe{
		Type:              "type",
		DefinedInMaterial: 0,
		EntryPoint:        "point",
		Arguments:         Arguments{CFLAGS: "balls"},
	}
}

func createMaterial() {

}

func create() Provenance {
	material := Material{
		URI:    "",
		Digest: nil,
	}
	return Provenance{
		Type:          "type",
		Subject:       nil,
		PredicateType: "pre",
		Predicate:     createPredicate(),
	}
}
