package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile             string
	clusterKey          bool
	clusterKeyName      string
	clusterKeyNamespace string
	keyFile             string
	inPlace             bool

	envPrefix = "YAMLE"

	rootCmd = &cobra.Command{
		Use:   "yamle",
		Short: "A simple YAML encrypter",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.yamle.yaml)")
	rootCmd.PersistentFlags().BoolVar(&inPlace, "in-place", false, "encrypt or decrypt file in place")

	rootCmd.PersistentFlags().BoolVar(&clusterKey, "cluster-key", false, "use keys stored in Kubernetes")
	rootCmd.PersistentFlags().StringVar(&clusterKeyName, "cluster-key-name", "", "the name of the secret in Kubernetes")
	rootCmd.PersistentFlags().StringVar(&clusterKeyNamespace, "cluster-key-namespace", "", "the namespace of the secret in Kubernetes")

	rootCmd.PersistentFlags().StringVar(&keyFile, "key-file", "", "path to a PEM formatted key")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(genkeysCmd)
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find current working directory.
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in current working directory with name ".yamle" (without extension).
		v.AddConfigPath(cwd)
		v.SetConfigType("yaml")
		v.SetConfigName(".yamle")
	}

	if err := v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", v.ConfigFileUsed())
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
