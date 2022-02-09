package nodejs

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nais/salsa/pkg/scan/common"
)

const yarnBuildFileName = "yarn.lock"

type Yarn struct {
	BuildFilePatterns []string
}

func NewYarn() common.BuildTool {
	return &Yarn{
		BuildFilePatterns: []string{yarnBuildFileName},
	}
}

func (y Yarn) ResolveDeps(workDir string) (*common.ArtifactDependencies, error) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, yarnBuildFileName))
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps := YarnDeps(string(fileContent))

	return &common.ArtifactDependencies{
		Cmd:         yarnBuildFileName,
		RuntimeDeps: deps,
	}, nil
}

func (y Yarn) BuildFiles() []string {
	return y.BuildFilePatterns
}

func YarnDeps(yarnLockContents string) []common.Dependency {
	deps := make([]common.Dependency, 0)
	lines := strings.Split(yarnLockContents, "\n")
	blockLines := blockLineNumbers(lines)
	for _, startLine := range blockLines {
		depName := parseDependency(lines[startLine])
		depVersion := parseVersion(lines[startLine+1])
		integrityLine := lines[startLine+3]

		deps = append(deps, common.Dependency{
			Coordinates: depName,
			Version:     depVersion,
			CheckSum:    yarnShaDigest(integrityLine),
		})
	}
	return deps
}

func blockLineNumbers(yarnLockLines []string) []int {
	var startsOfEntries []int
	for index, line := range yarnLockLines {
		if strings.HasPrefix(line, "\"") {
			startsOfEntries = append(startsOfEntries, index)
		}
	}
	return startsOfEntries
}

func parseDependency(depLine string) string {
	if len(strings.Split(depLine, ", ")) > 1 {
		depLine = parseName(depLine)
		return strings.Split(depLine, ", ")[1]
	} else {
		return parseName(depLine)
	}
}

func parseName(line string) string {
	regex := regexp.MustCompile(`^"?(?P<pkgname>.*)@[^~]?.*$`)
	matches := regex.FindStringSubmatch(line)
	pkgnameIndex := regex.SubexpIndex("pkgname")
	return matches[pkgnameIndex]
}

func parseVersion(line string) string {
	regex := regexp.MustCompile(`.*"(?P<pkgversion>.*)"$`)
	matches := regex.FindStringSubmatch(line)
	pkgversionIndex := regex.SubexpIndex("pkgversion")
	return matches[pkgversionIndex]
}

func yarnShaDigest(line string) common.CheckSum {
	trimPrefixIntegrity := strings.TrimPrefix(line, "  integrity ")
	fields := strings.Split(trimPrefixIntegrity, "-")
	// Better to keep the digest base64 encoded when signing envelope
	// decodedDigest, _ := base64.StdEncoding.DecodeString(algoDigest[1])
	// s.Digest = fmt.Sprintf("%x", decodedDigest)
	return common.CheckSum{
		Algorithm: fields[0],
		Digest:    fields[1],
	}
}
