package build_tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nais/salsa/pkg/intoto"
	"github.com/nais/salsa/pkg/scan"
	"github.com/nais/salsa/pkg/utils"
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
)

type BuildTool interface {
	Build(workDir, project string, context *vcs.AnyContext) error
	BuildTool(pattern string) bool
	BuildFiles() []string
}

func Scan(workingDir, project string, inputContext *string) error {
	context, err := vcs.CreateCIContext(inputContext)
	if err != nil {
		return err
	}

	gradle := NewGradle()
	mvn := NewMaven()
	golang := NewGolang()
	npm := NewNpm()
	yarn := NewYarn()
	composer := NewComposer()

	supportedBuildFiles := sumSupported(
		gradle.BuildFiles(),
		mvn.BuildFiles(),
		golang.BuildFiles(),
		npm.BuildFiles(),
		yarn.BuildFiles(),
		composer.BuildFiles(),
	)

	for index, pattern := range supportedBuildFiles {
		log.Printf("search for build type '%s'\n", pattern)
		buildFile := findBuildFile(workingDir, pattern)
		if index < len(supportedBuildFiles) {
			if buildFile != "" {
				log.Printf("found build type %s", buildFile)
				s := utils.StartSpinner(fmt.Sprintf("provenace generated for %s in %s", project, workingDir))
				switch true {
				case gradle.BuildTool(buildFile):
					err := gradle.Build(workingDir, project, context)
					if err != nil {
						return err
					}
				case mvn.BuildTool(buildFile):
					err := mvn.Build(workingDir, project, context)
					if err != nil {
						return err
					}
				case golang.BuildTool(buildFile):
					err := golang.Build(workingDir, project, context)
					if err != nil {
						return err
					}
				case npm.BuildTool(buildFile):
					err := npm.Build(workingDir, project, context)
					if err != nil {
						return err
					}
				case yarn.BuildTool(buildFile):
					err := yarn.Build(workingDir, project, context)
					if err != nil {
						return err
					}
				case composer.BuildTool(buildFile):
					err := composer.Build(workingDir, project, context)
					if err != nil {
						return err
					}

					// add more cases
				}
				// found break out!
				s.Stop()
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

func GenerateProvenance(workDir, project string, buildToolMetadata *scan.BuildToolMetadata, context *vcs.AnyContext) error {
	app := intoto.CreateApp(project, buildToolMetadata).With(context)
	predicate := intoto.GenerateSlsaPredicate(app)
	statement, err := json.Marshal(predicate)
	if err != nil {
		return fmt.Errorf("marshal: %v\n", err)
	}

	log.Println(string(statement))

	provenanceFileName := utils.ProvenanceFile(project)
	err = os.WriteFile(workDir+"/"+provenanceFileName, statement, 0644)
	if err != nil {
		return fmt.Errorf("write to file: %v\n", err)
	}
	return nil
}
