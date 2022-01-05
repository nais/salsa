package intoto

import (
	"fmt"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFindsAllMaterials(t *testing.T) {
	attPath := "testdata/cosign-dsse-attestation.json"
	fileContents, err := os.ReadFile(attPath)
	assert.NoError(t, err)

	statement, err := ParseEnvelope(fileContents)

	fmt.Printf("yo %v", statement.Predicate)

	assert.NoError(t, err)

	p := &in_toto.ProvenancePredicate{}
	err = mapstructure.Decode(statement.Predicate, p)

	assert.NotEmpty(t, p.Materials)
}

func TestFindMaterial(t *testing.T) {
	valueToFind := "com.google.guava:guava"
	attPath := "testdata/cosign-dsse-attestation.json"
	fileContents, err := os.ReadFile(attPath)
	assert.NoError(t, err)

	statement, err := ParseEnvelope(fileContents)
	assert.NoError(t, err)

	p := &in_toto.ProvenancePredicate{}
	err = mapstructure.Decode(statement.Predicate, p)

	foundMaterial := FindMaterials(*p, valueToFind)
	assert.NotEmpty(t, foundMaterial)
}
