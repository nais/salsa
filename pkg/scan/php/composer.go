package php

import (
	"encoding/json"
	"fmt"
	"github.com/nais/salsa/pkg/scan"
)

type Dist struct {
	Shasum string `json:"shasum"`
}

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Dist    Dist   `json:"dist"`
}

type ComposerLock struct {
	Dependencies []Dependency `json:"packages"`
}

func ComposerDeps(composerLockJsonContents string) (*scan.BuildToolMetadata, error) {
	var lock ComposerLock
	err := json.Unmarshal([]byte(composerLockJsonContents), &lock)
	if err != nil {
		return nil, fmt.Errorf("unable to parse composer.lock: %v", err)
	}

	return transform(lock.Dependencies), nil
}

func transform(dependencies []Dependency) *scan.BuildToolMetadata {
	metadata := scan.CreateMetadata()

	for _, dep := range dependencies {
		metadata.Deps[dep.Name] = dep.Version
		metadata.Checksums[dep.Name] = scan.CheckSum{
			Algorithm: "sha1",
			Digest:    dep.Dist.Shasum,
		}
	}
	return metadata
}
