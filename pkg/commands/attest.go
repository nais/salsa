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

type AttestCmd interface {
	Run(args []string, runner utils.CmdRunner) error
}

var attestCmd = &cobra.Command{
	Use:   "attest",
	Short: "sign and upload in-toto attestation",
	RunE: func(cmd *cobra.Command, args []string) error {
		var options AttestOptions
		err := viper.Unmarshal(&options)
		if err != nil {
			return err
		}

		_, err = options.Run(args, utils.ExecCmd{})
		if err != nil {
			return err
		}
		return nil
	},
}

func (o AttestOptions) Run(args []string, runner utils.CmdRunner) (string, error) {
	if PathFlags.Repo == "" {
		return "", errors.New("repo name must be specified")
	}

	if o.PredicateFile == "" {
		file := utils.ProvenanceFile(PathFlags.Repo)
		log.Infof("no predicate specified, using default pattern %s", file)
		o.PredicateFile = file
	}

	workDir := PathFlags.WorkDir()
	s := utils.StartSpinner(fmt.Sprintf("finished attestation for %s\n", PathFlags.Repo))
	defer s.Stop()
	filePath := fmt.Sprintf("%s/%s.%s", workDir, PathFlags.Repo, "att")
	// TODO: could be a subcommand e.g bin/salsa attest verify
	if verify {
		cmd := o.verifyCmd(args, runner)
		out, err := cmd.Run()
		if err != nil {
			return "", err
		}
		docs := strings.Split(out, "\n")
		if len(docs) > 0 {
			// TODO: fix so that we dont have to make this assumption
			//remove last line which is a newline
			docs := docs[:len(docs)-1]
			if len(docs) == 0 {
				return "", errors.New(fmt.Sprintf("unexpected output from cosign command: %s", out))
			}

			doc := docs[len(docs)-1]

			err = os.WriteFile(filePath, []byte(doc), os.FileMode(0755))
			if err != nil {
				return "", fmt.Errorf("could not write file %s %w", filePath, err)
			}

			err := os.WriteFile(fmt.Sprintf("%s/%s.%s", workDir, PathFlags.Repo, "raw.txt"), []byte(out), os.FileMode(0755))
			if err != nil {
				return "", fmt.Errorf("could not write file %s %w", workDir, err)
			}
		} else {
			log.Infof("no attestations found from cosign verify-attest command")
		}
		return out, nil
	} else {
		cmd := o.attestCmd(args, runner)
		out, err := cmd.Run()

		if err != nil {
			return "", err
		}
		if o.NoUpload {
			err = os.WriteFile(filePath, []byte(out), os.FileMode(0755))
			if err != nil {
				return "", fmt.Errorf("could not write file %s %w", filePath, err)
			}
		}
		return out, nil
	}
}

func (o AttestOptions) verifyCmd(a []string, runner utils.CmdRunner) utils.Cmd {
	return utils.Cmd{
		Name:    "cosign",
		SubCmd:  "verify-attestation",
		Flags:   []string{"--key", o.Key},
		Args:    a,
		WorkDir: PathFlags.WorkDir(),
		Runner:  runner,
	}
}

func (o AttestOptions) attestCmd(a []string, runner utils.CmdRunner) utils.Cmd {
	return utils.Cmd{
		Name:   "cosign",
		SubCmd: "attest",
		Flags: []string{
			"--key", o.Key,
			"--predicate", o.PredicateFile,
			"--type", o.PredicateType,
			"--rekor-url", o.RekorURL,
			fmt.Sprintf("--no-upload=%s", strconv.FormatBool(o.NoUpload)),
		},
		Args:    a,
		WorkDir: PathFlags.WorkDir(),
		Runner:  runner,
	}
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
