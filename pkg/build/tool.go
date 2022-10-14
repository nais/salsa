package build

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	AlgorithmSHA256 = "sha256"
	AlgorithmSHA1   = "sha1"
)

type Tool interface {
	BuildFiles() []string
	ResolveDeps(workDir string) (*ArtifactDependencies, error)
}

type Tools struct {
	Tools []Tool
}

func (t Tools) DetectDeps(workDir string) (*ArtifactDependencies, error) {
	log.Info("search for build files\n")
	for _, tool := range t.Tools {
		foundMatch, buildFile, err := match(tool, workDir)
		if err != nil {
			return nil, fmt.Errorf("could not find match, %v", err)
		}

		if foundMatch {
			log.Infof("found build type '%s'\n", buildFile)
			deps, err := tool.ResolveDeps(workDir)
			if err != nil {
				return nil, fmt.Errorf("could not resolve deps: %v", err)
			}

			return deps, nil
		}
	}
	return nil, fmt.Errorf(fmt.Sprintf("no supported build files found: %s", workDir))
}

func match(t Tool, workDir string) (bool, string, error) {
	for _, file := range t.BuildFiles() {
		buildFile, err := findBuildFile(workDir, file)

		if err != nil {
			return false, "", err
		}

		if file == buildFile && len(buildFile) != 0 {
			return true, buildFile, nil
		}
	}
	return false, "", nil
}

func findBuildFile(root, pattern string) (string, error) {
	var result = ""
	files, err := os.ReadDir(root)
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
