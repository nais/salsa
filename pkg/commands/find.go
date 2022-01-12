package commands

import (
	"errors"
	"io/ioutil"
	"os"
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

		files, err := ioutil.ReadDir("./attestations/")
		if err != nil {
			return err
		}

		for _, file := range files {
			var attFilePath = "./attestations/" + file.Name()
			fileContents, _ := os.ReadFile(attFilePath)
			provenance, err := intoto.ParseEnvelope(fileContents)
			if err != nil {
				return err
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
