package intoto

import (
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"strings"
)

func FindMaterials(materials []common.ProvenanceMaterial, value string) []common.ProvenanceMaterial {
	found := make([]common.ProvenanceMaterial, 0)
	for _, m := range materials {
		if find(m, value) {
			found = append(found, m)
		}
	}
	return found
}

func find(material common.ProvenanceMaterial, value string) bool {
	return strings.Contains(material.URI, value)
}
