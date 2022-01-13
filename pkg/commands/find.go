package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nais/salsa/pkg/intoto"
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
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("could not read dir %w", err)
		}

		for _, file := range files {
			var attFilePath = "./" + path + "/" + file.Name()

			if ext := filepath.Ext(file.Name()); ext != ".att" {
				continue
			}

			fileContents, _ := os.ReadFile(attFilePath)
			provenance, err := intoto.ParseEnvelope(fileContents)
			if err != nil {
				return fmt.Errorf("could not read file %s, %w", attFilePath, err)
			}
			result := intoto.FindMaterials(provenance.Predicate.Materials, artifact)
			if len(result) > 0 {
				app := strings.Split(file.Name(), ".")[0]
				log.Infof("found dependency %s in app %s", result, app)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.Flags().StringVar(&artifact, "artifact", "", "artifact to search after")
	viper.BindPFlags(findCmd.Flags())
}
