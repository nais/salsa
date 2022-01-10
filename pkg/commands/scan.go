package commands

import (
	"github.com/nais/salsa/pkg/build-tool"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var project string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan files and dependencies for a given project",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: check if path exists, if not create
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		log.Infof("prepare to scan path %s for project %s...", repoPath, project)
		// TODO: generalize into other build tools
		// Check repo for build file
		err := build_tool.Scan(repoPath, project)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVar(&project, "project", "tokendings", "project name")
}
