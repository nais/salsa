package jvm

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/nais/salsa/pkg/digest"
	"github.com/nais/salsa/pkg/scan"
	"regexp"
	"strings"
)

func GradleDeps(depsOutput string) (*scan.BuildToolMetadata, error) {
	gradleMetatdata := scan.CreateMetadata()
	pattern := regexp.MustCompile(`(?m)---\s[a-zA-Z0-9.]+:.*$`)
	matches := pattern.FindAllString(depsOutput, -1)
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

func GradleDepsAndSums(metadata *scan.BuildToolMetadata, outputSums []byte) error {
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
