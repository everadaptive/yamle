package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var genkeysCmd = &cobra.Command{
	Use:   "gen-keys",
	Short: "Generate RSA keys to use for encryption and decryption",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}
