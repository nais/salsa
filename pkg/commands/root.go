package commands

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/spf13/viper"
)

const (
	envVarPrefix = "SALSA"
	cmdName      = "salsa"
	defaultPath  = "tmp"
)

var (
	RepoPath string
	cfgFile  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   cmdName,
	Short: "Command-line tool for SLSA (SALSA)",
	Long:  `Scan files and dependencies, sign them and upload them`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&RepoPath, "repoPath", defaultPath, "project to scan")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/."+cmdName+".yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	v := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".salsa" (without extension).
		v.AddConfigPath(home)
		v.SetConfigType("yaml")
		v.SetConfigName("." + cmdName)
	}

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", v.ConfigFileUsed())
	}

	v.SetEnvPrefix(envVarPrefix + "_" + cmd.Name())

	v.AutomaticEnv() // read in environment variables that match
	bindFlags(cmd, v)
	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		bindEnv(v, cmd, f)
		if !f.Changed {
			setFlagIfPresent(v, cmd, f)
		}
	})
}

func bindEnv(v *viper.Viper, cmd *cobra.Command, f *pflag.Flag) {
	if strings.Contains(f.Name, "-") {
		suffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
		prefix := envVarPrefix + "_" + strings.ToUpper(cmd.Name())
		env := fmt.Sprintf("%s_%s", prefix, suffix)
		err := v.BindEnv(f.Name, env)
		if err != nil {
			fmt.Printf("could not bind to env: %v", err)
		}
	}
}

func setFlagIfPresent(v *viper.Viper, cmd *cobra.Command, f *pflag.Flag) {
	val := v.Get(f.Name)
	if val == nil {
		val = v.Get(cmd.Name() + "." + f.Name)
	}
	if val != nil {
		err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		if err != nil {
			fmt.Printf("could not set flag: %v", err)
			return
		}
	}
}
