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
	return deps
}

// TODO: Find answer to why dependencies has either '"' or starts without in header of dependencies, what's the difference?
// strings.HasPrefix(line, "\"") || strings.Contains(line, "@^") || strings.Contains(line, "@~")
// full files all, ut to what cost?
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

func lastElementInSlice(slice []string) string {
	return fmt.Sprintf("%v", slice[len(slice)-1])
}

func parseName(line string) string {
	regex := regexp.MustCompile(`^"?(?P<pkgname>.*)@[^~]?.*$`)
	matches := regex.FindStringSubmatch(line)
	pkgnameIndex := regex.SubexpIndex("pkgname")
	return matches[pkgnameIndex]
}

func parseVersion(line string) string {
	//fmt.Println(line)
	regex := regexp.MustCompile(`.*"(?P<pkgversion>.*)"$`)
	matches := regex.FindStringSubmatch(line)
	pkgversionIndex := regex.SubexpIndex("pkgversion")
	return matches[pkgversionIndex]
}

func yarnShaDigest(line string) build.CheckSum {
	trimPrefixIntegrity := strings.TrimPrefix(line, "  integrity ")
	fields := strings.Split(trimPrefixIntegrity, "-")
	// Better to keep the digest base64 encoded when signing envelope
	// decodedDigest, _ := base64.StdEncoding.DecodeString(algoDigest[1])
	// s.Digest = fmt.Sprintf("%x", decodedDigest)
	return build.CheckSum{
		Algorithm: fields[0],
		Digest:    fields[1],
	}
}
