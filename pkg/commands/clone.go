package commands

import (
	"errors"
	"github.com/nais/salsa/pkg/clone"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var owner string

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clones the given project into user defined path",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		if PathFlags.Repo == "" || owner == "" {
			return errors.New("repo and owner must be specified")
		}

		path := PathFlags.WorkDir()
		log.Infof("prepare to checkout %s into path %s ...", PathFlags.Repo, path)
		err := clone.Repo(owner, PathFlags.Repo, path, Auth.GithubToken)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVar(&owner, "owner", "", "owner of the repo")
	err := viper.BindPFlags(cloneCmd.Flags())
	if err != nil {
		log.Errorf("setting viper flag: %v", err)
	}
}
