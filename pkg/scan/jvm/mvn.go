package jvm

import (
	"errors"
	"regexp"
	"strings"
)

func MavenCompileAndRuntimeTimeDeps(mvnOutput string) (map[string]string, error) {
	pattern := regexp.MustCompile(`(?m)\s{4}[a-zA-Z0-9.]+:.*`)
	matches := pattern.FindAllString(mvnOutput, -1)
	if matches == nil {
		return nil, errors.New("unable to find any dependencies")
	}

	dependencies := make(map[string]string)
	for _, match := range matches {
		elements := strings.Split(strings.TrimSpace(match), ":")
		if elements[4] == "test" {
			continue
		}
		name := elements[0] + ":" + elements[1]
		dependencies[name] = elements[3]
	}
	return dependencies, nil
}
