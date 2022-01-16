package golang

import (
	"github.com/nais/salsa/pkg/digest"
	"github.com/nais/salsa/pkg/scan"
	"reflect"
	"testing"
)

func TestGoDeps(t *testing.T) {
	got := GoDeps(goSumContents)
	wantDeps := map[string]string{
		"github.com/google/uuid":  "1.0.0",
		"github.com/pborman/uuid": "1.2.1",
	}

	wantChecksum := map[string]scan.CheckSum{
		"github.com/google/uuid":  {Algorithm: digest.SHA256, Digest: "b4Gk+7WdP/d3HZH8EJsZpvV7EtDOgaZLtnaNGIu1adA="},
		"github.com/pborman/uuid": {Algorithm: digest.SHA256, Digest: "+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw="},
	}

	if !reflect.DeepEqual(got.Deps, wantDeps) {
		t.Errorf("got %q, wanted %q", got.Deps, wantDeps)
	}

	if !reflect.DeepEqual(got.Checksums, wantChecksum) {
		t.Errorf("got %v, wanted %v", got.Checksums, wantChecksum)
	}
}

const goSumContents = `
github.com/google/uuid v1.0.0 h1:b4Gk+7WdP/d3HZH8EJsZpvV7EtDOgaZLtnaNGIu1adA=
github.com/google/uuid v1.0.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/pborman/uuid v1.2.1 h1:+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw=
github.com/pborman/uuid v1.2.1/go.mod h1:X/NO0urCmaxf9VXbdlT7C2Yzkj2IKimNn4k+gtPdI/k=
`
