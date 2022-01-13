package golang

import (
	"reflect"
	"testing"
)

func TestGoDeps(t *testing.T) {
	got := GoDeps(goSumContents)
	wantDeps := map[string]string{
		"github.com/google/uuid":  "1.0.0",
		"github.com/pborman/uuid": "1.2.1",
	}

	wantChecksum := map[string]string{
		"github.com/google/uuid":  "6f81a4fbb59d3ff7771d91fc109b19a6f57b12d0ce81a64bb6768d188bb569d0",
		"github.com/pborman/uuid": "f99648c39f2dfe8cdd8d169787fddac077e65916f363126801d349c5eff7a6fc",
	}

	if !reflect.DeepEqual(got.Deps, wantDeps) {
		t.Errorf("got %q, wanted %q", got.Deps, wantDeps)
	}

	if !reflect.DeepEqual(got.Checksums, wantChecksum) {
		t.Errorf("got %q, wanted %q", got.Checksums, wantChecksum)
	}
}

const goSumContents = `
github.com/google/uuid v1.0.0 h1:b4Gk+7WdP/d3HZH8EJsZpvV7EtDOgaZLtnaNGIu1adA=
github.com/google/uuid v1.0.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/pborman/uuid v1.2.1 h1:+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw=
github.com/pborman/uuid v1.2.1/go.mod h1:X/NO0urCmaxf9VXbdlT7C2Yzkj2IKimNn4k+gtPdI/k=
`
