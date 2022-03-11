package php

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nais/salsa/pkg/build"
)

const composerLockFileName = "composer.lock"

type Composer struct {
	BuildFilePatterns []string
}

func (c Composer) BuildFiles() []string {
	return c.BuildFilePatterns
}

func NewComposer() build.Tool {
	return &Composer{
		BuildFilePatterns: []string{composerLockFileName},
	}
}

func (c Composer) ResolveDeps(workDir string) (*build.ArtifactDependencies, error) {
	path := fmt.Sprintf("%s/%s", workDir, composerLockFileName)
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps, err := ComposerDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("scan: %v\n", err)
	}
	return build.ArtifactDependency(deps, path, composerLockFileName), nil
}

type dist struct {
	Shasum string `json:"shasum"`
}

type dep struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Dist    dist   `json:"dist"`
}

type composerLock struct {
	Dependencies []dep `json:"packages"`
}

func ComposerDeps(composerLockJsonContents string) (map[string]build.Dependency, error) {
	var lock composerLock
	err := json.Unmarshal([]byte(composerLockJsonContents), &lock)
	if err != nil {
		return nil, fmt.Errorf("unable to parse composer.lock: %v", err)
	}

	return transform(lock.Dependencies), nil
}

func transform(dependencies []dep) map[string]build.Dependency {
	deps := make(map[string]build.Dependency, 0)
	for _, d := range dependencies {
		checksum := build.Verification(build.AlgorithmSHA1, d.Dist.Shasum)
		deps[d.Name] = build.Dependence(d.Name, d.Version, checksum)
	}
	return deps
}
