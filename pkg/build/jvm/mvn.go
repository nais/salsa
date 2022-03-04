package jvm

import (
	"crypto/sha256"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/utils"
)

const mavenBuildFileName = "pom.xml"

type Maven struct {
	BuildFilePatterns []string
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
	return &build.ArtifactDependencies{
		Cmd:         fmt.Sprintf("%s %v", cmd.Path, cmd.Args),
		RuntimeDeps: deps,
	}, nil
}

func NewMaven() build.BuildTool {
	return &Maven{
		BuildFilePatterns: []string{mavenBuildFileName},
	}
}

func (m Maven) BuildFiles() []string {
	return m.BuildFilePatterns
}

func MavenCompileAndRuntimeTimeDeps(rootPath string) ([]build.Dependency, error) {
	files, err := findJarFiles(rootPath)
	if err != nil {
		return nil, err
	}

	deps := make([]build.Dependency, 0)

	for _, file := range files {
		f := strings.Split(file, rootPath)[1]

		path := strings.Split(f, "/")
		version := path[len(path)-2]
		artifactId := path[len(path)-3]
		groupId := strings.Join(path[1:(len(path)-3)], ".")

		fmt.Printf("yolo %s:%s:%s\n", groupId, artifactId, version)
		digest, err := hashFile(file)
		if err != nil {
			return nil, err
		}
		deps = append(deps, build.Dependency{
			Coordinates: fmt.Sprintf("%s:%s", groupId, artifactId),
			Version:     version,
			CheckSum: build.CheckSum{
				Algorithm: "sha256",
				Digest:    digest,
			},
		})
	}
	return deps, nil
}

func hashFile(file string) (string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(content))

	return hash, nil
}

func findJarFiles(rootPath string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".jar" {
			files = append(files, path)
			//	fmt.Printf("File Name: %s path:%s\n", info.Name(), path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v", files)
	return files, nil
}
