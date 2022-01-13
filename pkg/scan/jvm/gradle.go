package jvm

import (
	"errors"
	"github.com/nais/salsa/pkg/scan"
	"regexp"
	"strings"
)

func GradleDeps(gradleOutput string) (*scan.BuildToolMetadata, error) {
	gradleMetatdata := scan.CreateMetadata()
	pattern := regexp.MustCompile(`(?m)---\s[a-zA-Z0-9.]+:.*$`)
	matches := pattern.FindAllString(gradleOutput, -1)
	if matches == nil {
		return nil, errors.New("unable to find any dependencies")
	}

	for _, match := range matches {
		replacedDownGrades := strings.Replace(match, " -> ", ":", -1)
		elements := strings.Split(strings.Replace(replacedDownGrades, "--- ", "", -1), ":")
		e := strings.Split(elements[2], " ")
		name := elements[0] + ":" + elements[1]
		gradleMetatdata.Deps[name] = e[0]
	}

	return gradleMetatdata, nil
}
