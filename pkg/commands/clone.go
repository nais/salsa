package commands

import (
	"errors"
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var owner string

type Principal struct {
	Username, Password string
}

var auth Principal

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
		err := vcs.CloneRepo(owner, PathFlags.Repo, path, auth.Username, auth.Password)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVar(&owner, "owner", "", "owner of the repo")
	cloneCmd.Flags().StringVar(&auth.Username, "username", "", "username for private/internal repos")
	cloneCmd.Flags().StringVar(&auth.Password, "password", "", "Github PAT token")
	viper.BindPFlags(cloneCmd.Flags())
}
