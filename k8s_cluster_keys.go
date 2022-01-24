package main

import (
	"context"
	"crypto/rsa"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type RSAKeyPair struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func SaveKeyToCluster(pubPem, privPem, name, namespace string) {
	s := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		StringData: make(map[string]string),
		Type:       "yamle.everadaptive.tech/rsa",
	}

	if namespace != "" {
		s.Namespace = namespace
	}

	s.StringData["public.pem"] = pubPem
	s.StringData["private.pem"] = privPem

	c, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Printf("error %s\n", err)
	}

	err = c.Create(context.Background(), &s)
	if err != nil {
		fmt.Printf("error %s\n", err)
	}
}

func GetKeyFromCluster(name, namespace string) (*RSAKeyPair, error) {
	c, err := client.New(config.GetConfigOrDie(), client.Options{})
	s := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	if namespace != "" {
		s.Namespace = namespace
	}

	err = c.Get(context.Background(), client.ObjectKeyFromObject(&s), &s)
	if err != nil {
		return nil, err
	}

	kp := RSAKeyPair{}

	if len(s.Data["private.pem"]) > 0 {
		kp.PrivateKey = exportPEMStrToPrivKey(s.Data["private.pem"])
	}

	if len(s.Data["public.pem"]) > 0 {
		kp.PublicKey = exportPEMStrToPubKey(s.Data["public.pem"])
	}

	return &kp, nil
}
