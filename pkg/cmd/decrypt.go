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
	decryptCmd = &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt a YAML file",
		Run:   decrypt,
	}
)

func decrypt(cmd *cobra.Command, args []string) {
	var privKey *rsa.PrivateKey

	if clusterKey {
		key, _ := keys.GetKeyFromCluster(clusterKeyName, clusterKeyNamespace)
		privKey = key.PrivateKey
	} else {
		keyPEM, err := keys.ReadKeyFromFile(keyFile)
		if err != nil {
			fmt.Printf("error reading key '%s': %s\n", keyFile, err)
			return
		}
		privKey = keys.ExportPEMStrToPrivKey(keyPEM)
	}

	if privKey == nil {
		fmt.Printf("No private key available for decryption")
		return
	}

	for _, fileName := range args {
		ctx := context.Background()
		ctx = context.WithValue(ctx, yamle.PrivateKey, privKey)
		ctx = context.WithValue(ctx, yamle.PublicKey, nil)

		yamle.DoIt(ctx, "!encrypted", fileName)
	}
}
