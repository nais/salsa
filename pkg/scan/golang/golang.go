package golang

import (
	"github.com/nais/salsa/pkg/digest"
	"github.com/nais/salsa/pkg/scan"
	"strings"
)

func GoDeps(goSumContents string) *scan.BuildToolMetadata {
	goMetatdata := scan.CreateMetadata()
	lines := strings.Split(goSumContents, "\n")
	for _, line := range lines {
		if isNotInteresting(line) {
			continue
		}
		parts := strings.Split(line, " ")
		goMetatdata.Deps[parts[0]] = parts[1][1:]
		strings.Split(parts[2], ":")
		stringDecodedDigest := digest.Digest(strings.Split(parts[2], ":")[1])
		goMetatdata.Checksums[parts[0]] = scan.CheckSum{Algorithm: digest.SHA256, Digest: string(stringDecodedDigest)}

	}
	return goMetatdata
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
