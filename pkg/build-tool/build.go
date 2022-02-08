package build_tool

import (
	"fmt"
	"io/ioutil"

	"github.com/nais/salsa/pkg/scan"
	log "github.com/sirupsen/logrus"
)

// TODO: rename?
type BuildTool interface {
	Supported(buildFile string) bool
	BuildFiles() []string
	ResolveDeps(workDir string) (*scan.ArtifactDependencies, error)
}

type BuildTools struct {
	Tools []BuildTool
}

func (b BuildTools) supportedBuildFiles() []string {
	files := make([]string, 0)
	for _, tool := range b.Tools {
		files = concat(files, tool.BuildFiles())
	}
	return files
}

func Match(t BuildTool, workingDir string) bool {
	for _, file := range t.BuildFiles() {
		buildFile := findBuildFile(workingDir, file)
		if buildFile != "" {
			return true
		}
	}
	return false
}

func Dependencies(workingDir string) (*scan.ArtifactDependencies, error) {
	supported := BuildTools{
		Tools: []BuildTool{
			NewGradle(),
			NewMaven(),
			NewGolang(),
			NewNpm(),
			NewYarn(),
			NewComposer(),
		},
	}
	for _, tool := range supported.Tools {
		log.Infof("search for build files '%s'\n", tool.BuildFiles())
		if Match(tool, workingDir) {
			log.Infof("found build type")
			deps, err := tool.ResolveDeps(workingDir)
			if err != nil {
				return nil, fmt.Errorf("could not resolve deps, %v", err)
			}
			return deps, nil
		}
	}
	return nil, nil
}

func findBuildFile(root, pattern string) (result string) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == pattern {
			result = file.Name()
			break
		}
	}
	return result
}
