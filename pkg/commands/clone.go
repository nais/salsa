package commands

import (
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clones the given project into user defined path",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		url := viper.GetString("url")
		// TODO: check if path exists, if not create
		log.Infof("prepare to checkout %s into path %s ...", url, repoPath)
		vcs.CloneRepo(url, repoPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// onboardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// onboardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cloneCmd.Flags().String("url", "https://github.com/someorg/somerepo", "repo to clone")
	viper.BindPFlags(cloneCmd.Flags())
}
