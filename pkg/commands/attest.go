package commands

import (
	"github.com/nais/salsa/pkg/exec"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const predicateType = "slsaprovenance"

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

		c := exec.CmdCfg{
			WorkDir: repoPath,
			Cmd:     "cosign",
			Args: []string{
				"attest",
				"--type",
				predicateType,
				"--predicate",
				predicateFile,
				"--key",
				kmsUrl,
				"--rekor-url",
				logUrl,
				"--no-upload=true",
				image,
			},
		}
		command, err := c.Exec()
		if err != nil {
			return err
		}
		os.Mkdir("attestations", os.FileMode(0755))
		os.WriteFile("./attestations/"+predicateFile+".json", []byte(command.Output), os.FileMode(0755))
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
