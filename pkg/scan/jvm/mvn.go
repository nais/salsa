package jvm

import (
	"errors"
	"github.com/nais/salsa/pkg/scan"
	"regexp"
	"strings"
)

func MavenCompileAndRuntimeTimeDeps(mvnOutput string) (*scan.BuildToolMetadata, error) {
	mavenMetadata := scan.CreateMetadata()
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
		name := elements[0] + ":" + elements[1]
		mavenMetadata.Deps[name] = elements[3]
	}
	return mavenMetadata, nil
}
