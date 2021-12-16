package intoto

type Predicate struct {
	Builder   `json:"builder"`
	BuildType string `json:"buildType"`
	Metadata  `json:"metadata"`
	Recipe    `json:"recipe"`
	Materials []Material `json:"materials"`
}

type Builder struct {
	Id string `json:"id"`
}

type Recipe struct {
	Type              string    `json:"type"`
	DefinedInMaterial int       `json:"definedInMaterial"`
	EntryPoint        string    `json:"entryPoint"`
	Arguments         Arguments `json:"arguments"`
}

type Arguments struct {
	CFLAGS string `json:"CFLAGS"`
}

type Material struct {
	URI    string    `json:"uri"`
	Digest DigestSet `json:"digest"`
}

func (p Predicate) withMetadata(buildInvitationId, buildFinished string, rp bool) *Predicate {
	p.Metadata = Metadata{
		BuildInvocationId: buildInvitationId,
		Reproducible:      rp,
		BuildFinishedOn:   buildFinished,
	}.withCompleteness(false, false, false)
	return &p
}

func (p Predicate) withRecipe(recipeType, entryPoint, cflags string) *Predicate {
	p.Recipe = Recipe{
		Type:              recipeType,
		DefinedInMaterial: 0,
		EntryPoint:        entryPoint,
		Arguments:         Arguments{CFLAGS: cflags},
	}
	return &p
}

func (p Predicate) withMaterials(deps map[string]string) Predicate {
	materials := make([]Material, 0)
	for k, v := range deps {
		m := Material{
			URI:    k + ":" + v,
			Digest: nil,
		}
		materials = append(materials, m)
	}
	p.Materials = materials
	return p
}
