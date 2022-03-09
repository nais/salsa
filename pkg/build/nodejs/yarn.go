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

func NewYarn() build.Tool {
	return &Yarn{
		BuildFilePatterns: []string{yarnBuildFileName},
	}
}

func (y Yarn) BuildFiles() []string {
	return y.BuildFilePatterns
}

func (y Yarn) ResolveDeps(workDir string) (*build.ArtifactDependencies, error) {
	path := fmt.Sprintf("%s/%s", workDir, yarnBuildFileName)
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps := YarnDeps(string(fileContent))
	return build.ArtifactDependency(deps, path, yarnBuildFileName), nil
}

func YarnDeps(yarnLockContents string) map[string]build.Dependency {
	deps := make(map[string]build.Dependency, 0)
	lines := strings.Split(yarnLockContents, "\n")
	blockLines := blockLineNumbers(lines)
	for _, startLine := range blockLines {
		depName := parseDependency(lines[startLine])
		depVersion := parseVersion(lines[startLine+1])
		integrityLine := lines[startLine+3]
		checksum := buildChecksum(integrityLine)
		deps[depName] = build.Dependence(depName, depVersion, checksum)
	}
	return deps
}

func blockLineNumbers(yarnLockLines []string) []int {
	var startsOfEntries []int
	for index, line := range yarnLockLines {
		if isNewEntry(line) {
			startsOfEntries = append(startsOfEntries, index)
		}
	}
	return startsOfEntries
}

func isNewEntry(str string) bool {
	return !strings.HasPrefix(str, " ") && strings.HasSuffix(str, ":")
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

func buildChecksum(line string) build.CheckSum {
	trimPrefixIntegrity := strings.TrimPrefix(line, "  integrity ")
	fields := strings.Split(trimPrefixIntegrity, "-")
	return build.Verification(fields[0], fields[1])
}
