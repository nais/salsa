package jvm

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/nais/salsa/pkg/digest"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/utils"
)

const gradleBuildFileName = "build.gradle.kts"

type Gradle struct {
	BuildFilePatterns []string
}

func (g Gradle) ResolveDeps(workDir string) (*build.ArtifactDependencies, error) {
	cmd := exec.Command(
		"gradle",
		"-q", "dependencies", "--configuration", "runtimeClasspath", "-M", "sha256",
	)
	cmd.Dir = workDir

	err := utils.RequireCommand("gradle")
	if err != nil {
		return nil, fmt.Errorf("exec: %v\n", err)
	}

	depsOutput, err := utils.Exec(cmd)
	if err != nil {
		return nil, fmt.Errorf("exec: %v\n", err)
	}

	xml, err := ioutil.ReadFile(workDir + "/gradle/verification-metadata.xml")
	if err != nil {
		return nil, fmt.Errorf("readfile: %v\n", err)
	}

	deps, err := GradleDeps(depsOutput, xml)
	if err != nil {
		return nil, fmt.Errorf("could not get gradle deps: %w", err)
	}

	return &build.ArtifactDependencies{
		Cmd:         fmt.Sprintf("%s %v", cmd.Path, cmd.Args),
		RuntimeDeps: deps,
	}, nil
}

func NewGradle() build.BuildTool {
	return &Gradle{
		BuildFilePatterns: []string{gradleBuildFileName},
	}
}

func (g Gradle) BuildFiles() []string {
	return g.BuildFilePatterns
}

func GradleDeps(depsOutput string, checksumXml []byte) ([]build.Dependency, error) {
	pattern := regexp.MustCompile(`(?m)---\s[a-zA-Z0-9.]+:.*$`)
	matches := pattern.FindAllString(depsOutput, -1)
	if matches == nil {
		return nil, errors.New("unable to find any dependencies")
	}

	deps := make([]build.Dependency, 0)

	sum := GradleChecksum{}
	err := xml.Unmarshal(checksumXml, &sum)
	if err != nil {
		return nil, fmt.Errorf("xml parsing: %v", err)
	}

	for _, match := range matches {
		replacedDownGrades := strings.Replace(match, " -> ", ":", -1)
		elements := strings.Split(strings.Replace(replacedDownGrades, "--- ", "", -1), ":")
		groupId := elements[0]
		artifactId := elements[1]
		version := strings.Split(elements[2], " ")[0]
		deps = append(deps, build.Dependency{
			Coordinates: fmt.Sprintf("%s:%s", groupId, artifactId),
			Version:     version,
			CheckSum: build.CheckSum{
				Algorithm: digest.AlgorithmSHA256,
				Digest:    sum.checksum(groupId, artifactId, version),
			},
		})
	}

	return deps, nil
}

func (g GradleChecksum) checksum(groupId, artifactId, version string) string {
	for _, c := range g.Components.Components {
		if c.Group == groupId && c.Name == artifactId && c.Version == version {
			for _, a := range c.Artifacts {
				if strings.HasSuffix(a.Name, ".jar") {
					return a.Sha256.Value
				}
			}
		}
	}

	return ""
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
