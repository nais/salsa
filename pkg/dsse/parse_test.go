package dsse

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	attPath := "testdata/cosign-dsse-attestation.json"
	fileContents, err := os.ReadFile(attPath)

	assert.NoError(t, err)
	statement, err := ParseEnvelope(fileContents)

	assert.NoError(t, err)
	assert.NotEmpty(t, statement.Predicate.Materials)
}
