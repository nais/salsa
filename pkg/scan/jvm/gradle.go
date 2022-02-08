package jvm

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/nais/salsa/pkg/digest"
	"github.com/nais/salsa/pkg/scan"
)

// TODO: get gradle checksums
func GradleDeps(depsOutput string) ([]scan.Dependency, error) {
	pattern := regexp.MustCompile(`(?m)---\s[a-zA-Z0-9.]+:.*$`)
	matches := pattern.FindAllString(depsOutput, -1)
	if matches == nil {
		return nil, errors.New("unable to find any dependencies")
	}

	deps := make([]scan.Dependency, 0)

	for _, match := range matches {
		replacedDownGrades := strings.Replace(match, " -> ", ":", -1)
		elements := strings.Split(strings.Replace(replacedDownGrades, "--- ", "", -1), ":")
		groupId := elements[0]
		artifactId := elements[1]
		version := strings.Split(elements[2], " ")[0]
		deps = append(deps, scan.Dependency{
			Coordinates: fmt.Sprintf("%s:%s", groupId, artifactId),
			Version:     version,
			CheckSum: scan.CheckSum{
				Algorithm: "todo",
				Digest:    "todo",
			},
		})
	}

	return deps, nil
}

// TODO: use some of this for checksums above
func GradleDepsAndSums(metadata *scan.BuildArtifactMetadata, outputSums []byte) error {
	sum := GradleChecksum{}
	err := xml.Unmarshal(outputSums, &sum)
	if err != nil {
		return err
	}
	for _, c := range sum.Components.Components {
		depName := fmt.Sprintf("%s:%s", c.Group, c.Name)
		metadata.Deps[depName] = c.Version
		for _, a := range c.Artifacts {
			encodedChecksum := base64.StdEncoding.EncodeToString([]byte(a.Sha256.Value))
			if a.preferredArtifactType() {
				metadata.Checksums[depName] = scan.CheckSum{
					Algorithm: digest.SHA256, Digest: fmt.Sprintf("%s", encodedChecksum),
				}
			}
		}
	}
	return nil
}

type GradleChecksum struct {
	XMLName       xml.Name      `xml:"verification-metadata"`
	Configuration Configuration `xml:"configuration"`
	Components    Components    `xml:"components"`
}

type Configuration struct {
	XMLName          xml.Name `xml:"configuration"`
	VerifyMetadata   bool     `xml:"verify-metadata"`
	VerifySignatures bool     `xml:"verify-signatures"`
}

type Components struct {
	XMLName    xml.Name    `xml:"components"`
	Components []Component `xml:"component"`
}

type Component struct {
	XMLName   xml.Name   `xml:"component"`
	Group     string     `xml:"group,attr"`
	Name      string     `xml:"name,attr"`
	Version   string     `xml:"version,attr"`
	Artifacts []Artifact `xml:"artifact"`
}

type Artifact struct {
	XMLName xml.Name `xml:"artifact"`
	Name    string   `xml:"name,attr"`
	Sha256  Sha256   `xml:"sha256"`
}

type Sha256 struct {
	XMLName xml.Name `xml:"sha256"`
	Value   string   `xml:"value,attr"`
}

// for now, only jar is added, pom is also an alternativ
func (a Artifact) preferredArtifactType() bool {
	if strings.Contains(a.Name, ".pom") {
		return false
	}
	return true
}
