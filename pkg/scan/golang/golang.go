package golang

import (
	"strings"

	"github.com/nais/salsa/pkg/digest"
	"github.com/nais/salsa/pkg/scan"
)

func GoDeps(goSumContents string) []scan.Dependency {
	deps := make([]scan.Dependency, 0)
	lines := strings.Split(goSumContents, "\n")
	for _, line := range lines {
		if isNotInteresting(line) {
			continue
		}
		parts := strings.Split(line, " ")
		version := parts[1][1:]
		stringDecodedDigest := digest.Digest(strings.Split(parts[2], ":")[1])
		deps = append(deps, scan.Dependency{
			Coordinates: parts[0],
			Version:     version,
			CheckSum:    scan.CheckSum{Algorithm: digest.SHA256, Digest: string(stringDecodedDigest)},
		})
	}
	return deps
}

func isNotInteresting(line string) bool {
	return isEmpty(line) || isMod(line)
}

func isEmpty(line string) bool {
	return strings.TrimSpace(line) == ""
}

func isMod(line string) bool {
	idx := strings.Index(line, "go.mod")
	return idx > -1
}
