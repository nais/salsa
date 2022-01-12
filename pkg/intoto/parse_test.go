package intoto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFindsAllMaterials(t *testing.T) {
	attPath := "testdata/cosign-dsse-attestation.json"
	fileContents, err := os.ReadFile(attPath)

	assert.NoError(t, err)
	statement, err := ParseEnvelope(fileContents)

	assert.NoError(t, err)
	assert.NotEmpty(t, statement.Predicate.Materials)
	fmt.Println(statement.Predicate.Materials)
}

func TestFindMaterial(t *testing.T) {
	valueToFind := "com.google.guava:guava"
	attPath := "testdata/cosign-dsse-attestation.json"
	fileContents, err := os.ReadFile(attPath)
	assert.NoError(t, err)

	assert.NoError(t, err)
	statement, err := ParseEnvelope(fileContents)

	assert.NoError(t, err)
	assert.NotEmpty(t, statement.Predicate.Materials)

	foundMaterial := FindMaterials(statement.Predicate.Materials, valueToFind)
	assert.NotEmpty(t, foundMaterial)
	assert.Contains(t, foundMaterial[0].URI, valueToFind)
	fmt.Println(len(foundMaterial))
	assert.Equal(t, 1, len(foundMaterial))
}
