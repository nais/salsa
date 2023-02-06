package intoto

import (
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"strings"
)

func FindMaterials(materials []slsa.ProvenanceMaterial, value string) []slsa.ProvenanceMaterial {
	found := make([]slsa.ProvenanceMaterial, 0)
	for _, m := range materials {
		if find(m, value) {
			found = append(found, m)
		}
	}
	return found
}

func find(material slsa.ProvenanceMaterial, value string) bool {
	return strings.Contains(material.URI, value)
}
