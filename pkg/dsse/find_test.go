package dsse

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

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
	assert.Equal(t, 1, len(foundMaterial))
}
