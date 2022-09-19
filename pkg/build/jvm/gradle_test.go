package jvm

import (
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/build/test"
	"io/ioutil"
	"testing"

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

func TestBuildGradle(t *testing.T) {
	tests := []test.IntegrationTest{
		{
			Name:      "find build file and parse output",
			BuildType: BuildGradle(),
			WorkDir:   "testdata/jvm/gradle",
			BuildPath: "/usr/local/bin/gradle",
			Cmd:       "gradle -q dependencies --configuration runtimeClasspath -M sha256",
			Want: test.Want{
				Key:     "org.jetbrains.kotlin:kotlin-reflect",
				Version: "1.6.10",
				Algo:    "sha256",
				Digest:  "3277ac102ae17aad10a55abec75ff5696c8d109790396434b496e75087854203",
			},
		},
		{
			Name:         "cant find Gradle build file",
			BuildType:    BuildGradle(),
			WorkDir:      "testdata/whatever",
			Error:        true,
			ErrorMessage: "could not find match, reading dir open testdata/whatever: no such file or directory",
		},
		{
			Name:         "cant find supported build type",
			BuildType:    BuildMaven(""),
			WorkDir:      "testdata/jvm/gradle",
			Error:        true,
			ErrorMessage: "no supported build files found: testdata/jvm/gradle",
		},
	}

	test.Run(t, tests)
}
