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
	encryptCmd = &cobra.Command{
		Use:   "encrypt <file>",
		Short: "Encrypt a YAML file",
		Args:  cobra.MinimumNArgs(1),
		Run:   encrypt,
	}
)

func encrypt(cmd *cobra.Command, args []string) {
	var pubKey *rsa.PublicKey

	if clusterKey {
		key, _ := keys.GetKeyFromCluster(clusterKeyName, clusterKeyNamespace)
		pubKey = key.PublicKey
		if pubKey == nil {
			pubKey = &key.PrivateKey.PublicKey
		}
	} else {
		keyPEM, err := keys.ReadKeyFromFile(keyFile)
		if err != nil {
			cobra.CheckErr(fmt.Errorf("error reading key file %s", err))
		}
		pubKey = keys.ExportPEMStrToPubKey(keyPEM)
		if pubKey == nil {
			privKey := keys.ExportPEMStrToPrivKey(keyPEM)
			pubKey = &privKey.PublicKey
		}
	}

	if pubKey == nil {
		cobra.CheckErr(fmt.Errorf("no public key available in '%s'", keyFile))
	}

	for _, fileName := range args {
		ctx := context.Background()
		ctx = context.WithValue(ctx, yamle.PrivateKey, nil)
		ctx = context.WithValue(ctx, yamle.PublicKey, pubKey)

		out, err := yamle.DoIt(ctx, "!encrypt", fileName)
		cobra.CheckErr(err)

		if inPlace {
			ioutil.WriteFile(fileName, out, 0644)
			return
		}

		fmt.Printf("%s\n", string(out))
	}
}
