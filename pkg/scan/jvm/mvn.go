package jvm

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/nais/salsa/pkg/scan"
)

func MavenCompileAndRuntimeTimeDeps(mvnOutput string) ([]scan.Dependency, error) {
	deps := make([]scan.Dependency, 0)
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
		deps = append(deps, scan.Dependency{
			Coordinates: fmt.Sprintf("%s:%s", elements[0], elements[1]),
			Version:     elements[3],
			CheckSum: scan.CheckSum{
				Algorithm: "todo",
				Digest:    "todo",
			},
		})
	}
	return deps, nil
}
