package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var tagResolvers = make(map[string]func(*yaml.Node, *CustomTagProcessor) (*yaml.Node, error))

type Fragment struct {
	content *yaml.Node
}

func (f *Fragment) UnmarshalYAML(value *yaml.Node, ctp *CustomTagProcessor) error {
	var err error
	// process includes in fragments
	f.content, err = resolveTags(value, ctp)
	return err
}

type CustomTagProcessor struct {
	target  interface{}
	pubKey  *rsa.PublicKey
	privKey *rsa.PrivateKey
}

func (i *CustomTagProcessor) UnmarshalYAML(value *yaml.Node) error {
	resolved, err := resolveTags(value, i)
	if err != nil {
		return err
	}
	return resolved.Decode(i.target)
}

func resolveTags(node *yaml.Node, ctp *CustomTagProcessor) (*yaml.Node, error) {
	for tag, fn := range tagResolvers {
		if node.Tag == tag {
			return fn(node, ctp)
		}
	}
	if node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode {
		var err error
		for i := range node.Content {
			node.Content[i], err = resolveTags(node.Content[i], ctp)
			if err != nil {
				return nil, err
			}
		}
	}
	return node, nil
}

func resolveEncrypt(node *yaml.Node, ctp *CustomTagProcessor) (*yaml.Node, error) {
	if node.Kind != yaml.ScalarNode {
		return nil, errors.New("!encrypt on a non-scalar node")
	}

	f := Fragment{
		content: node,
	}

	cipherText, _ := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		ctp.pubKey,
		[]byte(node.Value),
		nil,
	)

	f.content.Value = fmt.Sprintf("%s", base64.StdEncoding.EncodeToString(cipherText))
	f.content.Tag = "!encrypted"
	f.content.Style = yaml.TaggedStyle
	return f.content, nil
}

func resolveEncrypted(node *yaml.Node, ctp *CustomTagProcessor) (*yaml.Node, error) {
	if node.Kind != yaml.ScalarNode {
		return nil, errors.New("!encrypted on a non-scalar node")
	}

	f := Fragment{
		content: node,
	}

	v, err := base64.StdEncoding.DecodeString(string(node.Value))
	if err != nil {

	}

	cipherText, _ := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		ctp.privKey,
		v,
		nil,
	)

	f.content.Value = fmt.Sprintf("%s", cipherText)
	f.content.Tag = "!encrypt"
	f.content.Style = yaml.TaggedStyle
	return f.content, nil
}

func AddResolvers(tag string, fn func(*yaml.Node, *CustomTagProcessor) (*yaml.Node, error)) {
	tagResolvers[tag] = fn
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "f", "", "YAML file to parse.")

	var encrypt bool
	flag.BoolVar(&encrypt, "encrypt", false, "YAML file to parse.")

	var decrypt bool
	flag.BoolVar(&decrypt, "decrypt", false, "YAML file to parse.")

	var genKeys bool
	flag.BoolVar(&genKeys, "gen-keys", false, "YAML file to parse.")

	var pubKeyFile string
	flag.StringVar(&pubKeyFile, "pub-key-file", "./example/keys/pub.pem", "YAML file to parse.")

	var privKeyFile string
	flag.StringVar(&privKeyFile, "priv-key-file", "./example/keys/priv.pem", "YAML file to parse.")

	var clusterKey bool
	flag.BoolVar(&clusterKey, "cluster-key", false, "YAML file to parse.")

	var clusterKeySecret string
	flag.StringVar(&clusterKeySecret, "key-secret", "yamle-key", "YAML file to parse.")

	var clusterKeysNamespace string
	flag.StringVar(&clusterKeysNamespace, "key-namespace", "default", "YAML file to parse.")

	flag.Parse()

	if genKeys {
		fmt.Printf("Writing keys...\n")

		priv, pub := generateKeyPair(2048)
		if clusterKey {
			SaveKeyToCluster(exportPubKeyAsPEMStr(pub), exportPrivKeyAsPEMStr(priv), clusterKeySecret, clusterKeysNamespace)
		} else {
			if _, err := os.Stat(privKeyFile); err == nil {
				fmt.Printf("priv key already exists %s\n", privKeyFile)
				return
			}

			if _, err := os.Stat(pubKeyFile); err == nil {
				fmt.Printf("pub key already exists %s\n", pubKeyFile)
				return
			}

			saveKeyToFile(exportPrivKeyAsPEMStr(priv), privKeyFile)
			saveKeyToFile(exportPubKeyAsPEMStr(pub), pubKeyFile)
		}

		return
	}

	if encrypt == decrypt {
		fmt.Printf("Please choose exactly one action --decrypt or --encrypt\n")
		return
	}

	var pubKey *rsa.PublicKey
	var privKey *rsa.PrivateKey

	// Register custom tag resolvers
	if encrypt {
		if clusterKey {
			key, _ := GetKeyFromCluster(clusterKeySecret, clusterKeysNamespace)
			pubKey = key.PublicKey
		} else {
			pubKeyPEM := readKeyFromFile(pubKeyFile)
			pubKey = exportPEMStrToPubKey(pubKeyPEM)
		}

		AddResolvers("!encrypt", resolveEncrypt)
	}

	if decrypt {
		if clusterKey {
			key, _ := GetKeyFromCluster(clusterKeySecret, clusterKeysNamespace)
			privKey = key.PrivateKey
		} else {
			privKeyPEM := readKeyFromFile(privKeyFile)
			privKey = exportPEMStrToPrivKey(privKeyPEM)
		}

		if privKey == nil {
			fmt.Printf("No private key available for decryption")
			return
		}

		AddResolvers("!encrypted", resolveEncrypted)
	}

	if fileName == "" {
		fmt.Println("Please provide yaml file by using -f option")
		return
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}

	var s yaml.Node
	err = yaml.Unmarshal(yamlFile, &CustomTagProcessor{
		target:  &s,
		privKey: privKey,
		pubKey:  pubKey,
	})
	if err != nil {
		panic("Error encountered during unmarshalling")
	}

	d, err := yaml.Marshal(&s)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%s\n", string(d))
}
