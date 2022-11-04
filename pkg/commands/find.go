package commands

import (
	"errors"
	"fmt"
	"github.com/nais/salsa/pkg/dsse"
	"github.com/nais/salsa/pkg/intoto"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var artifact string
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "find artifact from attestations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			artifact = args[0]
		}

		if artifact == "" {
			return errors.New("missing artifact")
		}

		path := PathFlags.RepoDir
		dirs, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("could not read dir %w", err)
		}

		for _, dir := range dirs {
			files, err := os.ReadDir(fmt.Sprintf("./%s/%s", path, dir.Name()))
			if err != nil {
				return fmt.Errorf("could not read dir %w", err)
			}

			for _, file := range files {
				var attFilePath = fmt.Sprintf("./%s/%s/%s", path, dir.Name(), file.Name())

				if ext := filepath.Ext(file.Name()); ext != ".att" {
					continue
				}

				fileContents, err := os.ReadFile(attFilePath)
				if err != nil {
					return fmt.Errorf("read .att file content %s, %w", attFilePath, err)
				}

				provenance, err := dsse.ParseEnvelope(fileContents)
				if err != nil {
					return fmt.Errorf("could not read file %s, %w", attFilePath, err)
				}
				result := intoto.FindMaterials(provenance.Predicate.Materials, artifact)
				app := strings.Split(file.Name(), ".")[0]

				if len(result) == 0 {
					log.Infof("no dependcies where found in app %s", app)
				} else {
					log.Infof("found dependency(ies) in app %s:", app)
					for _, f := range result {
						log.Infof("-uri: %s", f.URI)
						for k, d := range f.Digest {
							log.Infof("--digest: %s:%s", k, d)
						}
					}
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.Flags().StringVar(&artifact, "artifact", "", "artifact to search after")
	err := viper.BindPFlags(findCmd.Flags())
	if err != nil {
		log.Errorf("setting viper flag: %v", err)
	}
}
