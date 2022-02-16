package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nais/salsa/pkg/vcs"
	"os"

	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/build/golang"
	"github.com/nais/salsa/pkg/build/jvm"
	"github.com/nais/salsa/pkg/build/nodejs"
	"github.com/nais/salsa/pkg/build/php"
	"github.com/nais/salsa/pkg/intoto"
	"github.com/nais/salsa/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var githubContext string
var runnerContext string
var envContext string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan files and dependencies for a given project",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		if len(args) == 1 {
			artifact = args[0]
		}

		if PathFlags.Repo == "" {
			return errors.New("repo name must be specified")
		}

		log.Infof("prepare to scan path %s for project %s...", PathFlags.WorkDir(), PathFlags.Repo)
		workDir := PathFlags.WorkDir()

		tools := build.SupportedBuildTools{
			Tools: []build.BuildTool{
				jvm.NewGradle(),
				jvm.NewMaven(),
				golang.NewGolang(),
				nodejs.NewNpm(),
				nodejs.NewYarn(),
				php.NewComposer(),
			},
		}

		deps, err := tools.DetectDeps(workDir)
		if deps == nil {
			return errors.New("could not find any supported build tools in " + workDir)
		} else if err != nil {
			return fmt.Errorf("detecting dependecies")
		}

		ciEnv, err := vcs.CreateCIEnvironment(&githubContext, &runnerContext, &envContext)
		if err != nil {
			return err
		}

		err = GenerateProvenance(workDir, PathFlags.Repo, deps, ciEnv)
		if err != nil {
			return err
		}
		return nil
	},
}

func GenerateProvenance(workDir, project string, dependencies *build.ArtifactDependencies, ciEnv *vcs.Environment) error {
	opts := intoto.CreateProvenanceOptions(project, dependencies, ciEnv)
	predicate := intoto.GenerateSlsaPredicate(opts)
	statement, err := json.Marshal(predicate)
	if err != nil {
		return fmt.Errorf("marshal: %v\n", err)
	}

	provenanceFileName := utils.ProvenanceFile(project)
	path := fmt.Sprintf("%s/%s", workDir, provenanceFileName)
	err = os.WriteFile(path, statement, 0644)
	if err != nil {
		return fmt.Errorf("write to file: %v\n", err)
	}

	log.Infof("generated provenance in file: %s", path)
	return nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVar(&githubContext, "github_context", "", "context of github environment")
	scanCmd.Flags().StringVar(&runnerContext, "runner_context", "", "context of runner environment")
	scanCmd.Flags().StringVar(&envContext, "env_context", "", "environmental variables of current context")
}
