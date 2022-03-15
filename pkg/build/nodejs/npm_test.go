package nodejs

import (
	"github.com/nais/salsa/pkg/build"
	"testing"
)

func TestPackageLockJsonParsing(t *testing.T) {
	got, _ := NpmDeps(packageLockContents)
	want := map[string]build.Dependency{}
	want["js-tokens"] = build.TestDependency("js-tokens", "4.0.0", "sha512", "45d2547e5704ddc5332a232a420b02bb4e853eef5474824ed1b7986cf84737893a6a9809b627dca02b53f5b7313a9601b690f690233a49bce0e026aeb16fcf29")
	want["loose-envify"] = build.TestDependency("loose-envify", "1.4.0", "sha512", "972bb13c6aff59f86b95e9b608bfd472751cd7372a280226043cee918ed8e45ff242235d928ebe7d12debe5c351e03324b0edfeb5d54218e34f04b71452a0add")
	want["object-assign"] = build.TestDependency("object-assign", "4.1.1", "sha1", "2109adc7965887cfc05cbbd442cac8bfbb360863")
	want["react"] = build.TestDependency("react", "17.0.2", "sha512", "82784fb7be62fddabfcf7ffaabfd1ab0fefc0f4bb9f760f92f5a5deccf0ff9d724e85bbf8c978bea25552b6ddfa6d494663f158dffbeef05c0f1435c94641c6c")

	build.AssertEqual(t, got, want)
}

const packageLockContents = `
{
  "name": "myproject",
  "version": "1.0.0",
  "lockfileVersion": 2,
  "requires": true,
  "packages": {
    "": {
      "name": "dillings",
      "version": "1.0.0",
      "license": "ISC",
      "dependencies": {
        "react": "^17.0.2"
      }
    },
    "node_modules/js-tokens": {
      "version": "4.0.0",
      "resolved": "https://registry.npmjs.org/js-tokens/-/js-tokens-4.0.0.tgz",
      "integrity": "sha512-RdJUflcE3cUzKiMqQgsCu06FPu9UdIJO0beYbPhHN4k6apgJtifcoCtT9bcxOpYBtpD2kCM6Sbzg4CausW/PKQ=="
    },
    "node_modules/loose-envify": {
      "version": "1.4.0",
      "resolved": "https://registry.npmjs.org/loose-envify/-/loose-envify-1.4.0.tgz",
      "integrity": "sha512-lyuxPGr/Wfhrlem2CL/UcnUc1zcqKAImBDzukY7Y5F/yQiNdko6+fRLevlw1HgMySw7f611UIY408EtxRSoK3Q==",
      "dependencies": {
        "js-tokens": "^3.0.0 || ^4.0.0"
      },
      "bin": {
        "loose-envify": "cli.js"
      }
    },
    "node_modules/object-assign": {
      "version": "4.1.1",
      "resolved": "https://registry.npmjs.org/object-assign/-/object-assign-4.1.1.tgz",
      "integrity": "sha1-IQmtx5ZYh8/AXLvUQsrIv7s2CGM=",
      "engines": {
        "node": ">=0.10.0"
      }
    },
    "node_modules/react": {
      "version": "17.0.2",
      "resolved": "https://registry.npmjs.org/react/-/react-17.0.2.tgz",
      "integrity": "sha512-gnhPt75i/dq/z3/6q/0asP78D0u592D5L1pd7M8P+dck6Fu/jJeL6iVVK23fptSUZj8Vjf++7wXA8UNclGQcbA==",
      "dependencies": {
        "loose-envify": "^1.1.0",
        "object-assign": "^4.1.1"
      },
      "engines": {
        "node": ">=0.10.0"
      }
    }
  },
  "dependencies": {
    "js-tokens": {
      "version": "4.0.0",
      "resolved": "https://registry.npmjs.org/js-tokens/-/js-tokens-4.0.0.tgz",
      "integrity": "sha512-RdJUflcE3cUzKiMqQgsCu06FPu9UdIJO0beYbPhHN4k6apgJtifcoCtT9bcxOpYBtpD2kCM6Sbzg4CausW/PKQ=="
    },
    "loose-envify": {
      "version": "1.4.0",
      "resolved": "https://registry.npmjs.org/loose-envify/-/loose-envify-1.4.0.tgz",
      "integrity": "sha512-lyuxPGr/Wfhrlem2CL/UcnUc1zcqKAImBDzukY7Y5F/yQiNdko6+fRLevlw1HgMySw7f611UIY408EtxRSoK3Q==",
      "requires": {
        "js-tokens": "^3.0.0 || ^4.0.0"
      }
    },
    "object-assign": {
      "version": "4.1.1",
      "resolved": "https://registry.npmjs.org/object-assign/-/object-assign-4.1.1.tgz",
      "integrity": "sha1-IQmtx5ZYh8/AXLvUQsrIv7s2CGM="
    },
    "react": {
      "version": "17.0.2",
      "resolved": "https://registry.npmjs.org/react/-/react-17.0.2.tgz",
      "integrity": "sha512-gnhPt75i/dq/z3/6q/0asP78D0u592D5L1pd7M8P+dck6Fu/jJeL6iVVK23fptSUZj8Vjf++7wXA8UNclGQcbA==",
      "requires": {
        "loose-envify": "^1.1.0",
        "object-assign": "^4.1.1"
      }
    }
  }
}`

func TestBuildNpm(t *testing.T) {
	tests := []build.IntegrationTest{
		{
			Name:      "find build file and parse output",
			BuildType: BuildNpm(),
			WorkDir:   "testdata/nodejs/npm",
			BuildPath: "testdata/nodejs/npm/package-lock.json",
			Cmd:       "package-lock.json",
			Want: build.Want{
				Key:     "@ampproject/remapping",
				Version: "2.1.0",
				Algo:    "sha512",
				Digest:  "779472b13949ee19b0e53c38531831718de590c7bdda7f2c5c272e2cf0322001caea3f0379f0f0b469d554380e9eff919c4a2cba50c9f4d3ca40bdbb6c321dd2",
			},
		},
		{
			Name:         "cant find build file",
			BuildType:    BuildNpm(),
			WorkDir:      "testdata/whatever",
			Error:        true,
			ErrorMessage: "could not find match, reading dir open testdata/whatever: no such file or directory",
		},
	}

	build.RunTests(t, tests)
}
