package build

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type Tool interface {
	BuildFiles() []string
	ResolveDeps(workDir string) (*ArtifactDependencies, error)
}

type SupportedBuildTools struct {
	Tools []Tool
}

func (t SupportedBuildTools) DetectDeps(workDir string) (*ArtifactDependencies, error) {
	log.Info("search for build files\n")
	for _, tool := range t.Tools {
		foundMatch, err := match(tool, workDir)
		if err != nil {
			return nil, fmt.Errorf("could not find match, %v", err)
		}

		if foundMatch {
			log.Infof("found build type '%s'\n", tool.BuildFiles())
			deps, err := tool.ResolveDeps(workDir)
			if err != nil {
				return nil, fmt.Errorf("could not resolve deps, %v", err)
			}

			return deps, nil
		}
	}
	return nil, nil
}

func match(t Tool, workDir string) (bool, error) {
	for _, file := range t.BuildFiles() {
		buildFile, err := findBuildFile(workDir, file)

		if err != nil {
			return false, err
		}

		if buildFile != "" {
			return true, nil
		}
	}
	return false, nil
}

func findBuildFile(root, pattern string) (string, error) {
	var result = ""
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return "", fmt.Errorf("reading dir %v", err)
	}

	for _, file := range files {
		if file.Name() == pattern {
			result = file.Name()
			break
		}
	}
	return result, nil
}
