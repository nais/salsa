package deps

import (
	"fmt"
	"io/ioutil"

	"github.com/nais/salsa/pkg/scan/common"
	"github.com/nais/salsa/pkg/scan/golang"
	"github.com/nais/salsa/pkg/scan/jvm"
	"github.com/nais/salsa/pkg/scan/nodejs"
	"github.com/nais/salsa/pkg/scan/php"
	log "github.com/sirupsen/logrus"
)

func Dependencies(workDir string) (*common.ArtifactDependencies, error) {
	supported := common.BuildTools{
		Tools: []common.BuildTool{
			jvm.NewGradle(),
			jvm.NewMaven(),
			golang.NewGolang(),
			nodejs.NewNpm(),
			nodejs.NewYarn(),
			php.NewComposer(),
		},
	}
	for _, tool := range supported.Tools {
		log.Infof("search for build files '%s'\n", tool.BuildFiles())
		if match(tool, workDir) {
			log.Infof("found build type")
			deps, err := tool.ResolveDeps(workDir)
			if err != nil {
				return nil, fmt.Errorf("could not resolve deps, %v", err)
			}
			return deps, nil
		}
	}
	return nil, nil
}

func match(t common.BuildTool, workDir string) bool {
	for _, file := range t.BuildFiles() {
		buildFile := findBuildFile(workDir, file)
		if buildFile != "" {
			return true
		}
	}
	return false
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

func concat(slices ...[]string) (result []string) {
	for _, slice := range slices {
		for _, s := range slice {
			result = append(result, s)
		}
	}
	return result
}
