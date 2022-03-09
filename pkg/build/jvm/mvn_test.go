package jvm

import (
	"github.com/nais/salsa/pkg/build/test"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/nais/salsa/pkg/build"
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
