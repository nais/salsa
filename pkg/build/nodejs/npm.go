package nodejs

import (
	"encoding/json"
	"fmt"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/utils"
	"os"
	"strings"
)

const npmBuildFileName = "package-lock.json"

type Npm struct {
	BuildFilePatterns []string
}

func BuildNpm() build.Tool {
	return &Npm{
		BuildFilePatterns: []string{npmBuildFileName},
	}
}

func (n Npm) BuildFiles() []string {
	return n.BuildFilePatterns
}

func (n Npm) ResolveDeps(workDir string) (*build.ArtifactDependencies, error) {
	path := fmt.Sprintf("%s/%s", workDir, npmBuildFileName)
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps, err := NpmDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing deps: %v\n", err)
	}
	return build.ArtifactDependency(deps, path, npmBuildFileName), nil
}

func NpmDeps(packageLockJsonContents string) (map[string]build.Dependency, error) {
	var f interface{}
	err := json.Unmarshal([]byte(packageLockJsonContents), &f)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s: %v", packageLockJsonContents, err)
	}

	raw := f.(map[string]interface{})

	trans, err := transform(raw["dependencies"].(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("try to derive dependecies from %s %w", packageLockJsonContents, err)
	}
	return trans, nil
}

func transform(input map[string]interface{}) (map[string]build.Dependency, error) {
	deps := make(map[string]build.Dependency, 0)
	for key, value := range input {
		dependency := value.(map[string]interface{})
		integrity := fmt.Sprintf("%s", dependency["integrity"])
		shaDig := strings.Split(integrity, "-")
		decodedDigest, err := utils.DecodeDigest(shaDig[1])
		if err != nil {
			return nil, err
		}
		checksum := build.Verification(shaDig[0], decodedDigest)
		deps[key] = build.Dependence(key, fmt.Sprintf("%s", dependency["version"]), checksum)
	}
	return deps, nil
}
