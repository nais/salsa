package commands

import (
	"errors"

	"github.com/nais/salsa/pkg/build-tool"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var project string
var inputContext string

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

		log.Infof("prepare to scan path %s for project %s...", PathFlags.WorkDir(), project)
		// TODO: generalize into other build tools
		err := build_tool.Scan(PathFlags.WorkDir(), PathFlags.Repo, &inputContext)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVar(&inputContext, "context", "", "context of build environment")
}
