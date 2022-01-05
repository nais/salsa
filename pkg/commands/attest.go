package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var kmsUrl string
var predicateFile string
var logUrl string
var image string

var attestCmd = &cobra.Command{
	Use:   "attest",
	Short: "sign and upload in-toto attestation",
	RunE: func(cmd *cobra.Command, args []string) error {

		log.Infof("prepare to sign predicate %s ...", predicateFile)

		err := RequireCommand("cosign")
		if err != nil {
			return err
		}

		c := CmdCfg{
			workDir: repoPath,
			cmd:     "cosign",
			args: []string{
				"attest",
				"--predicate",
				predicateFile,
				"--key",
				kmsUrl,
				"--rekor-url",
				logUrl,
				image,
			},
		}
		command, err := c.Exec()
		if err != nil {
			return err
		}
		log.Infof("finished signing and uploading attestation. %s", command.Output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(attestCmd)
	attestCmd.Flags().StringVar(&kmsUrl, "kmsUrl", "default", "the kmsUrl used to sign attestation")
	attestCmd.Flags().StringVar(&predicateFile, "predicateFile", "default", "the attestation file")
	attestCmd.Flags().StringVar(&logUrl, "logUrl", "default", "the logUrl for attestation storage")
	attestCmd.Flags().StringVar(&image, "image", "default", "the docker image")
	viper.BindPFlags(attestCmd.Flags())
}
