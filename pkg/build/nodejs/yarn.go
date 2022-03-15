package nodejs

import (
	"errors"
	"fmt"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/utils"
	"os"
	"regexp"
	"strings"
)

const yarnBuildFileName = "yarn.lock"

type Yarn struct {
	BuildFilePatterns []string
}

func BuildYarn() build.Tool {
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
		return nil, fmt.Errorf("read file: %w", err)
	}

	deps, err := YarnDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("try to derive depencies from: %s %w", path, err)
	}

	return build.ArtifactDependency(deps, path, yarnBuildFileName), nil
}

func YarnDeps(yarnLockContents string) (map[string]build.Dependency, error) {
	deps := make(map[string]build.Dependency, 0)
	lines := strings.Split(yarnLockContents, "\n")
	blockLines := blockLineNumbers(lines)
	for _, startLine := range blockLines {
		depName := parseDependency(lines[startLine])
		depVersion := parseVersion(lines[startLine+1])
		if !strings.Contains(lines[startLine+3], "integrity") {
			return nil, errors.New("integrity is missing")
		}
		integrityLine := lines[startLine+3]
		checksum, err := buildChecksum(integrityLine)
		if err != nil {
			return nil, fmt.Errorf("building checksum %w", err)
		}
		deps[depName] = build.Dependence(depName, depVersion, checksum)
	}
	return deps, nil
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

func buildChecksum(line string) (build.CheckSum, error) {
	trimPrefixIntegrity := strings.TrimPrefix(line, "  integrity ")
	fields := strings.Split(trimPrefixIntegrity, "-")
	decodedDigest, err := utils.DecodeDigest(fields[1])
	if err != nil {
		return build.CheckSum{}, err
	}
	return build.Verification(fields[0], decodedDigest), nil
}
