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
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true

		log.Infof("prepare to scan path %s ...", repoPath)

        c := CmdConfig{
            workDir: repoPath,
            cmd:     "cosign",
            args:    []string{"attest", "--predicate", predicateFile, "--key", kmsUrl, "--rekor-url", logUrl, image},
        }
        command, err := c.ExecuteCommand()
        if err != nil {
            log.Errorf("got error %w", err)
            return err
        }
        log.Printf("cosign output %s", command.Output)
        return nil
	},
}

func init() {
	rootCmd.AddCommand(attestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// onboardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// onboardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// attestCmd.Flags().String("repoPath", defaultPath, "project to scan")
    attestCmd.Flags().StringVar(&kmsUrl, "kmsUrl", "default", "the kmsUrl used to sign attestation")
    attestCmd.Flags().StringVar(&predicateFile, "predicateFile", "default", "the attestation file")
    attestCmd.Flags().StringVar(&logUrl, "logUrl", "default", "the logUrl for attestation storage")
    attestCmd.Flags().StringVar(&image, "image", "default", "the docker image")

	viper.BindPFlags(attestCmd.Flags())
}
