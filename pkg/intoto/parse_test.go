package intoto

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFindsAllMaterials(t *testing.T) {
	attPath := "testdata/cosign-dsse-attestation.json"
	fileContents, _ := os.ReadFile(attPath)
	env, _ := ParseEnvelope(fileContents)
	assert.NotEmpty(t, env.Predicate.Materials)
}
