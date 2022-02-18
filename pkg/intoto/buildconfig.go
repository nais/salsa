package intoto

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type BuildConfig struct {
	Commands []string `json:"commands"`
	// Indicates how to parse the strings in commands.
	Shell string `json:"shell"`
}

func GenerateBuildConfig(cmd *cobra.Command) *BuildConfig {

	flagsUsed := make([]*pflag.Flag, 0)
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			flagsUsed = append(flagsUsed, f)
		}
	})

	cmdFlags := ""
	for _, c := range flagsUsed {
		cmdFlags += fmt.Sprintf("--%s %s", c.Name, c.Value.String())
	}

	return &BuildConfig{
		Commands: []string{
			fmt.Sprintf("%s %s", cmd.CommandPath(), cmdFlags),
		},
		Shell: "bash",
	}
}
