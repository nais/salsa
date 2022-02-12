package nodejs

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nais/salsa/pkg/build"
)

const yarnBuildFileName = "yarn.lock"

type Yarn struct {
	BuildFilePatterns []string
}

func NewYarn() build.BuildTool {
	return &Yarn{
		BuildFilePatterns: []string{yarnBuildFileName},
	}
}

func (y Yarn) ResolveDeps(workDir string) (*build.ArtifactDependencies, error) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, yarnBuildFileName))
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps := YarnDeps(string(fileContent))

	return &build.ArtifactDependencies{
		Cmd:         yarnBuildFileName,
		RuntimeDeps: deps,
	}, nil
}

func (y Yarn) BuildFiles() []string {
	return y.BuildFilePatterns
}

func YarnDeps(yarnLockContents string) []build.Dependency {
	deps := make([]build.Dependency, 0)
	lines := strings.Split(yarnLockContents, "\n")
	blockLines := blockLineNumbers(lines)
	for _, startLine := range blockLines {
		depName := parseDependency(lines[startLine])
		depVersion := parseVersion(lines[startLine+1])
		integrityLine := lines[startLine+3]

		deps = append(deps, build.Dependency{
			Coordinates: depName,
			Version:     depVersion,
			CheckSum:    yarnShaDigest(integrityLine),
		})
	}
	return deduplicate(deps)
}

// TODO: Find answer to why header of dependencies in yarn.lock starts with either '"' or without,
// strings.HasPrefix(line, "\"") || strings.Contains(line, "@^") || strings.Contains(line, "@~")
// fulfills all cases, but to what cost?
// several of the dependencies come in duplicates with several versions.
// see: func deduplicate(deps []build.Dependency) []build.Dependency
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
		allPossibilities := strings.Split(depLine, ", ")
		return lastElementInSlice(allPossibilities)
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

func lastElementInSlice(slice []string) string {
	return trim(fmt.Sprintf("%v", slice[len(slice)-1]))
}

func trim(line string) string {
	return strings.TrimPrefix(line, "\"")
}

func parseVersion(line string) string {
	regex := regexp.MustCompile(`.*"(?P<pkgversion>.*)"$`)
	matches := regex.FindStringSubmatch(line)
	pkgversionIndex := regex.SubexpIndex("pkgversion")
	return matches[pkgversionIndex]
}

func yarnShaDigest(line string) build.CheckSum {
	trimPrefixIntegrity := strings.TrimPrefix(line, "  integrity ")
	fields := strings.Split(trimPrefixIntegrity, "-")
	return build.CheckSum{
		Algorithm: fields[0],
		Digest:    fields[1],
	}
}

func deduplicate(deps []build.Dependency) []build.Dependency {
	type DependencyKey struct{ coordinate string }
	var unique []build.Dependency
	transients := make(map[DependencyKey]int)

	for _, d := range deps {
		k := DependencyKey{d.Coordinates}
		// Overwrite with last
		if i, ok := transients[k]; ok {
			unique[i] = d
		} else {
			// recalculate size
			transients[k] = len(unique)
			unique = append(unique, d)
		}
	}
	return unique
}
