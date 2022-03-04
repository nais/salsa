package jvm

import (
	"reflect"
	"testing"

	"github.com/nais/salsa/pkg/build"
)

func TestMavenDeps(t *testing.T) {
	got, _ := MavenCompileAndRuntimeTimeDeps("testdata/target/dependency")
	want := []build.Dependency{
		mvnDep("org.springframework:spring-core", "5.3.16", "0903d17e58654a2c79f4e46df79dc73ccaa49b6edbc7c3278359db403b687f6e"),
		mvnDep("org.yaml:snakeyaml", "1.26", "d87d607e500885356c03c1cae61e8c2e05d697df8787d5aba13484c2eb76a844"),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func mvnDep(coordinates, version string, digest string) build.Dependency {
	return build.Dependency{
		Coordinates: coordinates,
		Version:     version,
		CheckSum: build.CheckSum{
			Algorithm: "sha256",
			Digest:    digest,
		},
	}
}
