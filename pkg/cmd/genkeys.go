package cmd

import (
	"fmt"

	"github.com/everadaptive/yamle/pkg/keys"
	"github.com/spf13/cobra"
)

var genkeysCmd = &cobra.Command{
	Use:   "gen-keys",
	Short: "Generate RSA keys to use for encryption and decryption",
	Run:   genkeys,
}

func genkeys(cmd *cobra.Command, args []string) {
	priv, pub := keys.GenerateKeyPair(keySize)

	if clusterKey {
		err := keys.SaveKeyToCluster(keys.ExportPubKeyAsPEMStr(pub), keys.ExportPrivKeyAsPEMStr(priv), clusterKeyName, clusterKeyNamespace)
		cobra.CheckErr(err)
	} else {
		err := keys.SaveKeyToFile(keys.ExportPubKeyAsPEMStr(pub), fmt.Sprintf("%s.%s", keyFile, "pub"))
		cobra.CheckErr(err)

		err = keys.SaveKeyToFile(keys.ExportPrivKeyAsPEMStr(priv), keyFile)
		cobra.CheckErr(err)
	}
}
