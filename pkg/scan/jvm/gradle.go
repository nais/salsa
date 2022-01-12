package jvm

import (
	"errors"
	"regexp"
	"strings"
)

func GradleDeps(gradleOutput string) (map[string]string, error) {
	pattern := regexp.MustCompile(`(?m)---\s[a-zA-Z0-9.]+:.*$`)
	matches := pattern.FindAllString(gradleOutput, -1)
	if matches == nil {
		return nil, errors.New("unable to find any dependencies")
	}

	dependencies := make(map[string]string)
	for _, match := range matches {
		elements := strings.Split(strings.Replace(match, "--- ", "", -1), ":")
		e := strings.Split(elements[2], " ")
		name := elements[0] + ":" + elements[1]
		dependencies[name] = e[0]
	}

	return dependencies, nil
}
