package build_tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nais/salsa/pkg/intoto"
	log "github.com/sirupsen/logrus"
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

	for index, pattern := range supportedBuildFiles {
		log.Printf("search for build type '%s'", pattern)
		buildFile := findBuildFile(workingDir, pattern)

		if index < len(supportedBuildFiles) {
			log.Printf("searching..")
			if buildFile != "" {
				log.Printf("found build type %s", buildFile)

				switch true {
				case gradle.BuildTool(buildFile):
					err := gradle.Build(workingDir, project)
					if err != nil {
						return err
					}
				case mvn.BuildTool(buildFile):
					err := mvn.Build(workingDir, project)
					if err != nil {
						return err
					}
				case golang.BuildTool(buildFile):
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

func findBuildFile(root, pattern string) (result string) {
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

func GenerateProvenance(workDir, project string, deps map[string]string) error {
	app := createApp(project, deps)
	s := intoto.GenerateSlsaPredicate(app)

	statement, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal: %v\n", err)
	}

	log.Println(string(statement))
	provenanceName := fmt.Sprintf("%s.provenance", project)

	err = os.WriteFile(fmt.Sprintf("%s/%s", workDir, provenanceName), statement, 0644)
	if err != nil {
		return fmt.Errorf("write to file: %v\n", err)
	}
	return nil
}

func createApp(name string, deps map[string]string) intoto.App {
	return intoto.App{
		Name:         name,
		BuilderId:    "todoId",
		BuildType:    "todoType",
		Dependencies: deps,
	}
}
