package cmd

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/everadaptive/yamle/pkg/keys"
	"github.com/everadaptive/yamle/pkg/yamle"
	"github.com/spf13/cobra"
)

var (
	encryptCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt a YAML file",
		Run:   encrypt,
	}
)

func encrypt(cmd *cobra.Command, args []string) {
	var pubKey *rsa.PublicKey

	if clusterKey {
		key, _ := keys.GetKeyFromCluster(clusterKeyName, clusterKeyNamespace)
		pubKey = key.PublicKey
	} else {
		keyPEM, err := keys.ReadKeyFromFile(keyFile)
		if err != nil {
			fmt.Printf("error reading key '%s': %s\n", keyFile, err)
			return
		}
		pubKey = keys.ExportPEMStrToPubKey(keyPEM)
	}

	if pubKey == nil {
		fmt.Printf("No private key available for encryption")
		return
	}

	for _, fileName := range args {
		ctx := context.Background()
		ctx = context.WithValue(ctx, yamle.PrivateKey, nil)
		ctx = context.WithValue(ctx, yamle.PublicKey, pubKey)

		yamle.DoIt(ctx, "!encrypt", fileName)
	}
}
