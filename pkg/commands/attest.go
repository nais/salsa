package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
		var attest AttestOptions
		err := viper.Unmarshal(&attest)
		if err != nil {
			return err
		}

		if PathFlags.Repo == "" {
			return errors.New("repo name must be specified")
		}

		if attest.PredicateFile == "" {
			file := utils.ProvenanceFile(PathFlags.Repo)
			log.Infof("no predicate specified, using default pattern %s", file)
			attest.PredicateFile = file
		}

		s := utils.StartSpinner(fmt.Sprintf("finished attestation for %s", PathFlags.Repo))
		filePath := PathFlags.RepoDir + "/" + PathFlags.Repo + ".att"
		if verify {
			raw, err := attest.Verify(args)
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
				err = os.WriteFile(PathFlags.RepoDir+"/"+PathFlags.Repo+".raw.txt", []byte(raw), os.FileMode(0755))
			} else {
				log.Infof("no attestations found from cosign verify-attest command")
			}
		} else {
			out, err := attest.Exec(args)
			if err != nil {
				return err
			}
			if attest.NoUpload {
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

func (o AttestOptions) Verify(a []string) (string, error) {
	err := utils.RequireCommand("cosign")
	if err != nil {
		return "", err
	}
	args := []string{
		"verify-attestation",
		"--key", o.Key,
	}
	args = append(args, a...)

	cmd := exec.Command(
		"cosign",
		args...,
	)

	cmd.Dir = PathFlags.WorkDir()
	return utils.Exec(cmd)
}

func (o AttestOptions) Exec(a []string) (string, error) {
	err := utils.RequireCommand("cosign")
	if err != nil {
		return "", err
	}
	args := []string{
		"attest",
		"--type", o.PredicateType,
		"--predicate", o.PredicateFile,
		"--key", o.Key,
		"--rekor-url", o.RekorURL,
	}
	if o.NoUpload {
		args = append(args, "--no-upload")
	}
	args = append(args, a...)

	cmd := exec.Command(
		"cosign",
		args...,
	)

	cmd.Dir = PathFlags.WorkDir()
	return utils.Exec(cmd)
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
	viper.BindPFlags(attestCmd.Flags())
}
