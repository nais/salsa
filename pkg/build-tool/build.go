package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/intoto"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

var supportedBuildFiles = []string{
	"build.gradle.kts",
}

type BuildTool interface {
	Build(workingDir, project string) error
	BuildTool(pattern string) bool
}

func Scan(workingDir, project string) error {
	gradle := NewGradle()

	for _, pattern := range supportedBuildFiles {
		log.Printf("search for build type '%s'", pattern)
		foundBuildType := findBuildType(workingDir, pattern)

		if foundBuildType != "" {
			log.Printf("found build type %s", foundBuildType)
			switch true {
			case gradle.BuildTool(foundBuildType):
				err := gradle.Build(workingDir, project)
				if err != nil {
					return err
				}
				// add more cases
			}
			break

		} else {
			return fmt.Errorf("unknown build type")
		}
	}
	return nil
}

func findBuildType(root, pattern string) (result string) {
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		// filepath.Ext(d.Name())
		if d.Name() == pattern {
			trimmed := strings.Trim(s, fmt.Sprintf("%s/", root))
			result = trimmed
			return io.EOF
		}
		return nil
	})
	if err == io.EOF {
		err = nil
	}
	return result
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func createApp(name string, deps map[string]string) intoto.App {
	return intoto.App{
		Name:         name,
		BuilderId:    "todoId",
		BuildType:    "todoType",
		Dependencies: deps,
	}
}
