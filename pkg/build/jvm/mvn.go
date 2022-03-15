package jvm

import (
	"crypto/sha256"
	"fmt"
	"github.com/nais/salsa/pkg/build"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nais/salsa/pkg/utils"
)

const mavenBuildFileName = "pom.xml"

type Maven struct {
	BuildFilePatterns []string
}

func BuildMaven() build.Tool {
	return &Maven{
		BuildFilePatterns: []string{mavenBuildFileName},
	}
}

func (m Maven) BuildFiles() []string {
	return m.BuildFilePatterns
}

func (m Maven) ResolveDeps(workDir string) (*build.ArtifactDependencies, error) {
	cmd := exec.Command(
		"mvn",
		"dependency:copy-dependencies",
		"-DincludeScope=runtime",
		"-Dmdep.useRepositoryLayout=true",
	)
	cmd.Dir = workDir
	rootPath := workDir + "/target/dependency"

	_, err := utils.Exec(cmd)
	if err != nil {
		return nil, fmt.Errorf("exec: %v\n", err)
	}

	deps, err := MavenCompileAndRuntimeTimeDeps(rootPath)
	if err != nil {
		return nil, fmt.Errorf("scan: %v\n", err)
	}
	return build.ArtifactDependency(deps, cmd.Path, strings.Join(cmd.Args, " ")), nil
}

func MavenCompileAndRuntimeTimeDeps(rootPath string) (map[string]build.Dependency, error) {
	files, err := findJarFiles(rootPath)
	if err != nil {
		return nil, err
	}

	deps := make(map[string]build.Dependency, 0)

	for _, file := range files {
		f := strings.Split(file, rootPath)[1]

		path := strings.Split(f, "/")
		version := path[len(path)-2]
		artifactId := path[len(path)-3]
		groupId := strings.Join(path[1:(len(path)-3)], ".")

		checksum, err := buildChecksum(file)
		if err != nil {
			return nil, err
		}
		coordinates := fmt.Sprintf("%s:%s", groupId, artifactId)
		deps[coordinates] = build.Dependence(coordinates, version, checksum)
	}
	return deps, nil
}

func buildChecksum(file string) (build.CheckSum, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return build.CheckSum{}, err
	}
	checksum := fmt.Sprintf("%x", sha256.Sum256(content))
	return build.Verification(build.AlgorithmSHA256, checksum), nil
}

func findJarFiles(rootPath string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".jar" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
