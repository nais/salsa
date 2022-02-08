package php

import (
	"encoding/json"
	"fmt"

	"github.com/nais/salsa/pkg/scan"
)

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

func ComposerDeps(composerLockJsonContents string) ([]scan.Dependency, error) {
	var lock composerLock
	err := json.Unmarshal([]byte(composerLockJsonContents), &lock)
	if err != nil {
		return nil, fmt.Errorf("unable to parse composer.lock: %v", err)
	}

	return transform(lock.Dependencies), nil
}

func transform(dependencies []dep) []scan.Dependency {
	deps := make([]scan.Dependency, 0)
	for _, d := range dependencies {
		deps = append(deps, scan.Dependency{
			Coordinates: d.Name,
			Version:     d.Version,
			CheckSum: scan.CheckSum{
				Algorithm: "sha1",
				Digest:    d.Dist.Shasum,
			},
		})
	}
	return deps
}
