package jvm

import (
	"github.com/nais/salsa/pkg/digest"
	"io/ioutil"
	"testing"

	"github.com/nais/salsa/pkg/build"
	"github.com/stretchr/testify/assert"
)

func TestGradleDeps(t *testing.T) {
	gradleOutput, _ := ioutil.ReadFile("testdata/gradle_output.txt")
	checksumXml, _ := ioutil.ReadFile("testdata/verification-metadata.xml")
	got, err := GradleDeps(string(gradleOutput), checksumXml)
	assert.NoError(t, err)
	want := []build.Dependency{
		dep("ch.qos.logback:logback-classic", "1.2.10", "3160ae988af82c8bf3024ddbe034a82da98bb186fd09e76c50543c5b9da5cc5e"),
	}

	count := 0
	for _, wantDep := range want {
		for _, gotDep := range got {
			if wantDep == gotDep {
				count++
			}
		}
	}

	assert.Equal(t, len(want), count)
}

func dep(coordinates, version, checksum string) build.Dependency {
	return build.Dependency{
		Coordinates: coordinates,
		Version:     version,
		CheckSum: build.CheckSum{
			Algorithm: digest.AlgorithmSHA256,
			Digest:    checksum,
		},
	}
}
