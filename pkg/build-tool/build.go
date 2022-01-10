package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/intoto"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type BuildTool interface {
	Build(workDir, project string) error
	BuildTool(pattern string) bool
	BuildFiles() []string
}

func Scan(workingDir, project string) error {
	gradle := NewGradle()
	mvn := NewMaven()
	golang := NewGolang()

	supportedBuildFiles := sumSupported(
		gradle.BuildFiles(),
		mvn.BuildFiles(),
		golang.BuildFiles(),
	)

	totalSupported := len(supportedBuildFiles)
	for index, pattern := range supportedBuildFiles {
		log.Printf("search for build type '%s'", pattern)
		foundBuildType := findBuildType(workingDir, pattern)

		if index+1 <= totalSupported {
			log.Printf("searching..")
			if foundBuildType != "" {
				log.Printf("found build type %s", foundBuildType)
				switch true {
				case gradle.BuildTool(foundBuildType):
					err := gradle.Build(workingDir, project)
					if err != nil {
						return err
					}
				case mvn.BuildTool(foundBuildType):
					err := mvn.Build(workingDir, project)
					if err != nil {
						return err
					}
				case golang.BuildTool(foundBuildType):
					err := golang.Build(workingDir, project)
					if err != nil {
						return err
					}

					// add more cases
				}
				// found break out!
				break
			}

		} else {
			return fmt.Errorf("unknown build type")
		}
	}
	return nil
}

func findBuildType(root, pattern string) (result string) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == pattern {
			result = file.Name()
			break
		}
	}
	return result
}

func createApp(name string, deps map[string]string) intoto.App {
	return intoto.App{
		Name:         name,
		BuilderId:    "todoId",
		BuildType:    "todoType",
		Dependencies: deps,
	}
}
