package cmd

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	"github.com/everadaptive/yamle/pkg/keys"
	"github.com/everadaptive/yamle/pkg/yamle"
	"github.com/spf13/cobra"
)

var (
	decryptCmd = &cobra.Command{
		Use:   "decrypt <file>",
		Short: "Decrypt a YAML file",
		Args:  cobra.MinimumNArgs(1),
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
			cobra.CheckErr(fmt.Errorf("error reading key file %s", err))
		}
		privKey = keys.ExportPEMStrToPrivKey(keyPEM)
	}

	if privKey == nil {
		cobra.CheckErr(fmt.Errorf("no private key available in '%s'", keyFile))
	}

	for _, fileName := range args {
		ctx := context.Background()
		ctx = context.WithValue(ctx, yamle.PrivateKey, privKey)
		ctx = context.WithValue(ctx, yamle.PublicKey, nil)

		out, err := yamle.DoIt(ctx, "!encrypted", fileName)
		cobra.CheckErr(err)

		if inPlace {
			ioutil.WriteFile(fileName, out, 0644)
			return
		}

		fmt.Printf("%s\n", string(out))
	}
}
