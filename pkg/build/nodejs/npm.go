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
	path := fmt.Sprintf("%s/%s", workDir, npmBuildFileName)
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps, err := NpmDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing deps: %v\n", err)
	}
	return &build.ArtifactDependencies{
		Cmd: build.Cmd{
			Path:     path,
			CmdFlags: npmBuildFileName,
		},
		RuntimeDeps: deps,
	}, nil
}

func NewNpm() build.Tool {
	return &Npm{
		BuildFilePatterns: []string{npmBuildFileName},
	}
}

func (n Npm) BuildFiles() []string {
	return n.BuildFilePatterns
}

func NpmDeps(packageLockJsonContents string) (map[string]build.Dependency, error) {
	var f interface{}
	err := json.Unmarshal([]byte(packageLockJsonContents), &f)
	if err != nil {
		return nil, fmt.Errorf("unable to parse package-lock.json: %v", err)
	}
	raw := f.(map[string]interface{})
	return transform(raw["dependencies"].(map[string]interface{})), nil
}

func transform(input map[string]interface{}) map[string]build.Dependency {
	deps := make(map[string]build.Dependency, 0)
	for key, value := range input {
		dependency := value.(map[string]interface{})
		integrity := fmt.Sprintf("%s", dependency["integrity"])
		shaDig := strings.Split(integrity, "-")
		checksum := build.CreateChecksum(fmt.Sprintf("%s", shaDig[0]), fmt.Sprintf("%s", shaDig[1]))
		deps[key] = build.CreateDependency(key, fmt.Sprintf("%s", dependency["version"]), checksum)
	}
	return deps
}
