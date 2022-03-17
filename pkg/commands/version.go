package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var info bool

var versionCmd = &cobra.Command{
	Use:   "version [flags]",
	Short: "Show 'salsa' client version",
	RunE: func(command *cobra.Command, args []string) error {
		if info {
			fmt.Printf("%s: %s commit: %s date: %s builtBy: %s",
				command.CommandPath(),
				Client.Version,
				Client.Commit,
				Client.Date, Client.BuiltBy,
			)
		}
		fmt.Printf("%s: %s", command.CommandPath(), Client.Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	scanCmd.Flags().BoolVar(&info, "info", false, "Detailed commit information")
}
