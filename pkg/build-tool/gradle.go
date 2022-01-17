package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/scan"
	"github.com/nais/salsa/pkg/scan/jvm"
	"github.com/nais/salsa/pkg/utils"
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

const gradleBuildFileName = "build.gradle.kts"

type Gradle struct {
	BuildFilePatterns []string
}

func NewGradle() BuildTool {
	return &Gradle{
		BuildFilePatterns: []string{gradleBuildFileName},
	}
}

func (g Gradle) Build(workDir string, project string, context *vcs.AnyContext) error {
	//cmd := exec.Command(
	//	"./gradlew",
	//	"-q", "dependencies", "--configuration", "runtimeClasspath",
	//)
	//cmd.Dir = workDir
	//
	//depsOutput, err := utils.Exec(cmd)
	//if err != nil {
	//	return fmt.Errorf("exec: %v\n", err)
	//}

	cmd := exec.Command(
		"./gradlew",
		"-M", "sha256",
	)
	cmd.Dir = workDir
	sumsOutput, err := utils.Exec(cmd)
	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}
	log.Info(sumsOutput)

	// deps, err := jvm.GradleDeps(depsOutput)
	depsSums := scan.CreateMetadata()
	log.Info(workDir)

	depsSumsFile, err := os.ReadFile(fmt.Sprintf("%s/gradle/verification-metadata.xml", workDir))
	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}

	err = jvm.GradleDepsAndSums(depsSums, depsSumsFile)
	if err != nil {
		return fmt.Errorf("scan: %v\n", err)
	}

	err = GenerateProvenance(workDir, project, depsSums, context)
	if err != nil {
		return fmt.Errorf("generating provencance %v", err)
	}
	return nil
}

func (g Gradle) BuildTool(pattern string) bool {
	return Contains(g.BuildFilePatterns, pattern)
}

func (g Gradle) BuildFiles() []string {
	return g.BuildFilePatterns
}
