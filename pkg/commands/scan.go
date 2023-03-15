package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/build/golang"
	"github.com/nais/salsa/pkg/build/jvm"
	"github.com/nais/salsa/pkg/build/nodejs"
	"github.com/nais/salsa/pkg/build/php"
	"github.com/nais/salsa/pkg/config"
	"github.com/nais/salsa/pkg/intoto"
	"github.com/nais/salsa/pkg/utils"
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	buildContext  string
	runnerContext string
	mvnOpts       string
	envContext    string
	Config        *ProvenanceConfig
)

type ProvenanceConfig struct {
	WithDependencies bool
	BuildStartedOn   string
}

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

		workDir := PathFlags.WorkDir()
		log.Infof("prepare to scan path %s ...", workDir)

		deps := &build.ArtifactDependencies{}
		if Config.WithDependencies {
			generatedDeps, err := InitBuildTools(mvnOpts).DetectDeps(workDir)
			if err != nil {
				return fmt.Errorf("detecting dependecies: %v", err)
			}

			if generatedDeps != nil {
				deps = generatedDeps
			} else {
				log.Infof("no supported build files found for directory: %s, proceeding", workDir)
			}
		}

		contextEnvironment, err := vcs.ResolveBuildContext(&buildContext, &runnerContext, &envContext)
		if err != nil {
			return err
		}

		scanConfiguration := &config.ScanConfiguration{
			BuildStartedOn:     Config.BuildStartedOn,
			WorkDir:            workDir,
			RepoName:           PathFlags.Repo,
			Dependencies:       deps,
			ContextEnvironment: contextEnvironment,
			Cmd:                cmd,
		}

		err = GenerateProvenance(scanConfiguration)
		if err != nil {
			return err
		}
		return nil
	},
}

func GenerateProvenance(scanCfg *config.ScanConfiguration) error {
	opts := intoto.CreateProvenanceOptions(scanCfg)
	predicate := intoto.GenerateSlsaPredicate(opts)
	statement, err := json.Marshal(predicate)
	if err != nil {
		return fmt.Errorf("marshal: %v\n", err)
	}

	provenanceFileName := utils.ProvenanceFile(scanCfg.RepoName)
	path := fmt.Sprintf("%s/%s", scanCfg.WorkDir, provenanceFileName)
	err = os.WriteFile(path, statement, 0644)
	if err != nil {
		return fmt.Errorf("write to file: %v\n", err)
	}

	log.Infof("generated provenance in file: %s", path)
	return nil
}

func InitBuildTools(mavenOpts string) build.Tools {
	return build.Tools{
		Tools: []build.Tool{
			jvm.BuildGradle(),
			jvm.BuildMaven(mavenOpts),
			golang.BuildGo(),
			nodejs.BuildNpm(),
			nodejs.BuildYarn(),
			php.BuildComposer(),
		},
	}
}

func init() {
	Config = &ProvenanceConfig{}
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVar(&buildContext, "build-context", "", "context of build tool")
	scanCmd.Flags().StringVar(&runnerContext, "runner-context", "", "context of runner")
	scanCmd.Flags().StringVar(&envContext, "env-context", "", "environmental variables of current context")
	scanCmd.Flags().BoolVar(&Config.WithDependencies, "with-deps", true, "specify if the cli should generate dependencies for a provenance")
	scanCmd.Flags().StringVar(&mvnOpts, "mvn-opts", "", "pass additional Comma-delimited list of options to the maven build tool")
	scanCmd.Flags().StringVar(&Config.BuildStartedOn, "build-started-on", "", "set start time for the build")
}
