package commands

import (
	"errors"

	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var url string

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clones the given project into user defined path",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		if PathFlags.Repo == "" || url == "" {
			return errors.New("repo and url must be specified")
		}

		path, err := PathFlags.WorkDir()
		if err != nil {
			return err
		}
		log.Infof("prepare to checkout %s into path %s ...", url, path)
		err = vcs.CloneRepo(url, path)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVar(&url, "url", "", "repo to clone")
	viper.BindPFlags(cloneCmd.Flags())
}
