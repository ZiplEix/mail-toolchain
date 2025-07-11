package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

var privateKey *rsa.PrivateKey

func LoadPrivateKey(path string) error {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	privateKey = privKey

	return nil
}

func EncodePrivateKeyToPEM(key *rsa.PrivateKey) ([]byte, error) {
	privDER := x509.MarshalPKCS1PrivateKey(key)

	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}
	return pem.EncodeToMemory(privBlock), nil
}
