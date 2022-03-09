package jvm

import (
	"github.com/nais/salsa/pkg/build/test"
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
	want := map[string]build.Dependency{}
	want["ch.qos.logback:logback-classic:"] = test.Dependency(
		"ch.qos.logback:logback-classic",
		"1.2.10",
		"sha256",
		"3160ae988af82c8bf3024ddbe034a82da98bb186fd09e76c50543c5b9da5cc5e",
	)
	want["org.jetbrains.kotlinx:kotlinx-coroutines-core"] = test.Dependency(
		"org.jetbrains.kotlinx:kotlinx-coroutines-core",
		"1.5.2-native-mt",
		"sha256",
		"78492527a0d09e0c53c81aacc2e073a83ee0fc3105e701496819ec67c98df16f",
	)
	want["com.fasterxml.jackson.core:jackson-annotations"] = test.Dependency(
		"com.fasterxml.jackson.core:jackson-annotations",
		"2.13.0",
		"sha256",
		"81f9724d8843e8b08f8f6c0609e7a2b030d00c34861c4ac7e2099a7235047d6f",
	)
	want["com.fasterxml.jackson.core:jackson-databind"] = test.Dependency(
		"com.fasterxml.jackson.core:jackson-databind",
		"2.13.0",
		"sha256",
		"9c826d27176268777adcf97e1c6e2051c7e33a7aaa2c370c2e8c6077fd9da3f4",
	)

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
