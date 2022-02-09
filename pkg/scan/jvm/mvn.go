package jvm

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nais/salsa/pkg/scan/common"
	"github.com/nais/salsa/pkg/utils"
)

const mavenBuildFileName = "pom.xml"

type Maven struct {
	BuildFilePatterns []string
}

func (m Maven) ResolveDeps(workDir string) (*common.ArtifactDependencies, error) {
	cmd := exec.Command(
		"mvn",
		"dependency:list",
	)
	cmd.Dir = workDir

	output, err := utils.Exec(cmd)
	if err != nil {
		return nil, fmt.Errorf("exec: %v\n", err)
	}

	deps, err := MavenCompileAndRuntimeTimeDeps(output)
	if err != nil {
		return nil, fmt.Errorf("scan: %v\n", err)
	}
	return &common.ArtifactDependencies{
		Cmd:         fmt.Sprintf("%s %v", cmd.Path, cmd.Args),
		RuntimeDeps: deps,
	}, nil
}

func NewMaven() common.BuildTool {
	return &Maven{
		BuildFilePatterns: []string{mavenBuildFileName},
	}
}

func (m Maven) BuildFiles() []string {
	return m.BuildFilePatterns
}

func MavenCompileAndRuntimeTimeDeps(mvnOutput string) ([]common.Dependency, error) {
	deps := make([]common.Dependency, 0)
	pattern := regexp.MustCompile(`(?m)\s{4}[a-zA-Z0-9.]+:.*`)
	matches := pattern.FindAllString(mvnOutput, -1)
	if matches == nil {
		return nil, errors.New("unable to find any dependencies")
	}
	for _, match := range matches {
		elements := strings.Split(strings.TrimSpace(match), ":")
		if elements[4] == "test" {
			continue
		}
		deps = append(deps, common.Dependency{
			Coordinates: fmt.Sprintf("%s:%s", elements[0], elements[1]),
			Version:     elements[3],
			CheckSum: common.CheckSum{
				Algorithm: "todo",
				Digest:    "todo",
			},
		})
	}
	return deps, nil
}
