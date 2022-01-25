package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile             string
	clusterKey          bool
	clusterKeyName      string
	clusterKeyNamespace string
	keyFile             string

	rootCmd = &cobra.Command{
		Use:   "yamle",
		Short: "A simple YAML encrypter",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")

	rootCmd.PersistentFlags().BoolVar(&clusterKey, "cluster-key", false, "use keys stored in Kubernetes")
	rootCmd.PersistentFlags().StringVar(&clusterKeyName, "cluster-key-name", "", "the name of the secret in Kubernetes")
	rootCmd.PersistentFlags().StringVar(&clusterKeyNamespace, "cluster-key-namespace", "", "the namespace of the secret in Kubernetes")

	rootCmd.PersistentFlags().StringVar(&keyFile, "key-file", "", "path to a PEM formatted key")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(genkeysCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
