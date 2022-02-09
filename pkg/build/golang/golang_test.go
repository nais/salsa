package golang

import (
	"reflect"
	"testing"

	"github.com/nais/salsa/pkg/build"
)

func TestGoDeps(t *testing.T) {
	got := GoDeps(goSumContents)
	wantDeps := []build.Dependency{
		dep("github.com/google/uuid", "1.0.0", "b4Gk+7WdP/d3HZH8EJsZpvV7EtDOgaZLtnaNGIu1adA="),
		dep("github.com/pborman/uuid", "1.2.1", "+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw="),
	}

	if !reflect.DeepEqual(got, wantDeps) {
		t.Errorf("got %q, wanted %q", got, wantDeps)
	}
}

func dep(coordinates, version, digest string) build.Dependency {
	return build.Dependency{
		Coordinates: coordinates,
		Version:     version,
		CheckSum: build.CheckSum{
			Algorithm: "sha256",
			Digest:    digest,
		},
	}
}

const goSumContents = `
github.com/google/uuid v1.0.0 h1:b4Gk+7WdP/d3HZH8EJsZpvV7EtDOgaZLtnaNGIu1adA=
github.com/google/uuid v1.0.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/pborman/uuid v1.2.1 h1:+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw=
github.com/pborman/uuid v1.2.1/go.mod h1:X/NO0urCmaxf9VXbdlT7C2Yzkj2IKimNn4k+gtPdI/k=
`
