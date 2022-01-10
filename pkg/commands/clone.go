package commands

import (
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/url"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clones the given project into user defined path",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		repoUrl, err := url.Parse(viper.GetString("url"))
		if err != nil {
			return err
		}

		// TODO: check if path exists, if not create
		log.Infof("prepare to checkout %s into path %s ...", repoUrl.Path, repoPath)
		err = vcs.CloneRepo(repoUrl.String(), repoPath)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().String("url", "https://github.com/someorg/somerepo", "repo to clone")
	viper.BindPFlags(cloneCmd.Flags())
}
