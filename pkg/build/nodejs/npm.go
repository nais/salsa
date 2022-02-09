package nodejs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nais/salsa/pkg/build"
)

const npmBuildFileName = "package-lock.json"

type Npm struct {
	BuildFilePatterns []string
}

func (n Npm) ResolveDeps(workDir string) (*build.ArtifactDependencies, error) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, npmBuildFileName))
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps, err := NpmDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing deps: %v\n", err)
	}
	return &build.ArtifactDependencies{
		Cmd:         npmBuildFileName,
		RuntimeDeps: deps,
	}, nil
}

func NewNpm() build.BuildTool {
	return &Npm{
		BuildFilePatterns: []string{npmBuildFileName},
	}
}

func (n Npm) BuildFiles() []string {
	return n.BuildFilePatterns
}

func NpmDeps(packageLockJsonContents string) ([]build.Dependency, error) {
	var f interface{}
	err := json.Unmarshal([]byte(packageLockJsonContents), &f)
	if err != nil {
		return nil, fmt.Errorf("unable to parse package-lock.json: %v", err)
	}
	raw := f.(map[string]interface{})
	return transform(raw["dependencies"].(map[string]interface{})), nil
}

func transform(input map[string]interface{}) []build.Dependency {
	deps := make([]build.Dependency, 0)
	for key, value := range input {
		dependency := value.(map[string]interface{})
		integrity := fmt.Sprintf("%s", dependency["integrity"])
		shaDig := strings.Split(integrity, "-")
		deps = append(deps, build.Dependency{
			Coordinates: key,
			Version:     fmt.Sprintf("%s", dependency["version"]),
			CheckSum: build.CheckSum{
				Algorithm: fmt.Sprintf("%s", shaDig[0]),
				Digest:    fmt.Sprintf("%s", shaDig[1]),
			},
		})
	}
	return deps
}
