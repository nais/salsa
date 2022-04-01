package intoto

import (
	"fmt"
	"github.com/nais/salsa/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

type BuildConfig struct {
	Commands []string `json:"commands"`
	// Indicates how to parse the strings in commands.
	Shell string `json:"shell"`
}

func GenerateBuildConfig(scanConfig *config.ScanConfiguration) *BuildConfig {
	buildConfig := &BuildConfig{
		Commands: []string{
			fmt.Sprintf("%s %s",
				scanConfig.Cmd.CommandPath(),
				salsaCmdFlags(scanConfig.Cmd),
			),
		},
		Shell: "bash",
	}

	if len(scanConfig.Dependencies.RuntimeDeps) > 0 {
		buildConfig.Commands = append(buildConfig.Commands, fmt.Sprintf("%s %s",
			scanConfig.Dependencies.CmdPath(),
			scanConfig.Dependencies.CmdFlags(),
		))
	}

	return buildConfig
}

func salsaCmdFlags(cmd *cobra.Command) string {
	flagsUsed := make([]*pflag.Flag, 0)
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			flagsUsed = append(flagsUsed, f)
		}
	})

	cmdFlags := ""
	for _, c := range flagsUsed {
		if strings.Contains(c.Name, "token") {
			cmdFlags += fmt.Sprintf(" --%s %s", c.Name, "***")
		} else {
			cmdFlags += fmt.Sprintf(" --%s %s", c.Name, c.Value.String())
		}
	}

	return cmdFlags
}
