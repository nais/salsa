package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var info bool

var versionCmd = &cobra.Command{
	Use:   "version [flags]",
	Short: "Show 'salsa' client version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if info {
			log.Infof("%s: %s commit: %s date: %s builtBy: %s",
				cmd.CommandPath(),
				Client.Version,
				Client.Commit,
				Client.Date,
				Client.BuiltBy,
			)
			return nil
		}
		log.Infof("%s: %s", cmd.CommandPath(), Client.Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&info, "info", false, "Detailed commit information")
}
