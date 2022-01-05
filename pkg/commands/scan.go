package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan files and dependencies for a given project",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: check if path exists, if not create
		log.Infof("prepare to scan path %s ...", repoPath)
		// TODO: generalize into other build tools
		GradleScan(repoPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
