package jvm

import (
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/build/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMavenDeps(t *testing.T) {
	got, err := MavenCompileAndRuntimeTimeDeps("testdata/target/dependency")
	assert.NoError(t, err)
	want := map[string]build.Dependency{}
	want["org.springframework:spring-core"] = test.Dependency(
		"org.springframework:spring-core",
		"5.3.16", "sha256",
		"0903d17e58654a2c79f4e46df79dc73ccaa49b6edbc7c3278359db403b687f6e",
	)
	want["org.yaml:snakeyaml"] = test.Dependency(
		"org.yaml:snakeyaml",
		"1.26",
		"sha256",
		"d87d607e500885356c03c1cae61e8c2e05d697df8787d5aba13484c2eb76a844",
	)

	test.AssertEqual(t, got, want)
}

func TestBuildMaven(t *testing.T) {
	tests := []test.IntegrationTest{
		{
			Name:      "find build file and parse output",
			BuildType: BuildMaven(""),
			WorkDir:   "testdata/jvm/maven",
			BuildPath: "/usr/local/bin/mvn",
			Cmd:       "mvn dependency:copy-dependencies -DincludeScope=runtime -Dmdep.useRepositoryLayout=true",
			Want: test.Want{
				Key:     "com.google.code.gson:gson",
				Version: "2.8.6",
				Algo:    "sha256",
				Digest:  "c8fb4839054d280b3033f800d1f5a97de2f028eb8ba2eb458ad287e536f3f25f",
			},
		},
		{
			Name:         "cant find build file",
			BuildType:    BuildMaven(""),
			WorkDir:      "testdata/whatever",
			Error:        true,
			ErrorMessage: "could not find match, reading dir open testdata/whatever: no such file or directory",
		},
		{
			Name:      "Add additional command line arguments as a part of the mvn command",
			BuildType: BuildMaven("-s .m2/maven-settings.xml"),
			Cmd:       "mvn dependency:copy-dependencies -DincludeScope=runtime -Dmdep.useRepositoryLayout=true -s.m2/maven-settings.xml",
			WorkDir:   "testdata/jvm/maven",
			BuildPath: "/usr/local/bin/mvn",
			Want: test.Want{
				Key:     "com.google.code.gson:gson",
				Version: "2.8.6",
				Algo:    "sha256",
				Digest:  "c8fb4839054d280b3033f800d1f5a97de2f028eb8ba2eb458ad287e536f3f25f",
			},
		},
		{
			Name:      "Add additional commands line arguments as a part of the mvn command",
			BuildType: BuildMaven("--also-make, --threads=2, --batch-mode, --settings=.m2/maven-settings.xml"),
			Cmd:       "mvn dependency:copy-dependencies -DincludeScope=runtime -Dmdep.useRepositoryLayout=true --also-make --threads=2 --batch-mode --settings=.m2/maven-settings.xml",
			WorkDir:   "testdata/jvm/maven",
			BuildPath: "/usr/local/bin/mvn",
			Want: test.Want{
				Key:     "com.google.code.gson:gson",
				Version: "2.8.6",
				Algo:    "sha256",
				Digest:  "c8fb4839054d280b3033f800d1f5a97de2f028eb8ba2eb458ad287e536f3f25f",
			},
		},
	}

	test.Run(t, tests)
}
