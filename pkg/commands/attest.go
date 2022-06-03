package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nais/salsa/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AttestOptions struct {
	Key           string `mapstructure:"key"`
	NoUpload      bool   `mapstructure:"no-upload"`
	RekorURL      string `mapstructure:"rekor-url"`
	PredicateFile string `mapstructure:"predicate"`
	PredicateType string `mapstructure:"type"`
}

var verify bool

var attestCmd = &cobra.Command{
	Use:   "attest",
	Short: "sign and upload in-toto attestation",
	RunE: func(cmd *cobra.Command, args []string) error {
		var options AttestOptions
		err := viper.Unmarshal(&options)
		if err != nil {
			return err
		}

		if PathFlags.Repo == "" {
			return errors.New("repo name must be specified")
		}

		if options.PredicateFile == "" {
			file := utils.ProvenanceFile(PathFlags.Repo)
			log.Infof("no predicate specified, using default pattern %s", file)
			options.PredicateFile = file
		}

		workDir := PathFlags.WorkDir()
		s := utils.StartSpinner(fmt.Sprintf("finished attestation for %s\n", PathFlags.Repo))
		filePath := fmt.Sprintf("%s/%s.%s", workDir, PathFlags.Repo, "att")
		// TODO: could be a subcommand e.g bin/salsa attest verify
		if verify {
			cmd := options.verifyCmd(args)
			raw, err := cmd.Run()
			if err != nil {
				return err
			}
			docs := strings.Split(raw, "\n")
			if len(docs) > 0 {
				// TODO: fix so that we dont have to make this assumption
				//remove last line which is a newline
				docs := docs[:len(docs)-1]
				doc := docs[len(docs)-1]

				err = os.WriteFile(filePath, []byte(doc), os.FileMode(0755))
				if err != nil {
					return fmt.Errorf("could not write file %s %w", filePath, err)
				}

				err := os.WriteFile(fmt.Sprintf("%s/%s.%s", workDir, PathFlags.Repo, "raw.txt"), []byte(raw), os.FileMode(0755))
				if err != nil {
					return fmt.Errorf("could not write file %s %w", workDir, err)
				}
			} else {
				log.Infof("no attestations found from cosign verify-attest command")
			}
		} else {
			cmd := options.attestCmd(args)
			out, err := cmd.Run()
			if err != nil {
				return err
			}
			if options.NoUpload {
				err = os.WriteFile(filePath, []byte(out), os.FileMode(0755))
				if err != nil {
					return fmt.Errorf("could not write file %s %w", filePath, err)
				}
			}
		}
		s.Stop()
		return nil
	},
}

func (o AttestOptions) verifyCmd(a []string) utils.Cmd {
	return utils.NewCmd(
		"cosign",
		"verify-attestation",
		[]string{"--key", o.Key},
		a,
		PathFlags.WorkDir(),
	)
}

func (o AttestOptions) attestCmd(a []string) utils.Cmd {
	return utils.NewCmd(
		"cosign",
		"attest",
		[]string{
			"--key", o.Key,
			"--predicate", o.PredicateFile,
			"--type", o.PredicateType,
			"--rekor-url", o.RekorURL,
			"--no-upload", strconv.FormatBool(o.NoUpload),
		},
		a,
		PathFlags.WorkDir(),
	)
}

func init() {
	rootCmd.AddCommand(attestCmd)
	attestCmd.Flags().String("key", "",
		"path to the private key file, KMS URI or Kubernetes Secret")
	attestCmd.Flags().BoolVar(&verify, "verify", false, "if true, verifies attestations - default is false")
	attestCmd.Flags().Bool("no-upload", false,
		"do not upload the generated attestation")
	attestCmd.Flags().String("rekor-url", "https://rekor.sigstore.dev",
		"address of transparency log")
	attestCmd.Flags().String("predicate", "",
		"the predicate file used for attestation")
	attestCmd.Flags().String("type", "slsaprovenance",
		"specify a predicate type (slsaprovenance|link|spdx|custom) or an URI (default \"slsaprovenance\")\n")
	err := viper.BindPFlags(attestCmd.Flags())
	if err != nil {
		log.Errorf("setting viper flag: %v", err)
	}
}
