package RSA

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GenerateRSAKeys(privateKeyPath, publicKeyPath string) error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048);
	if err != nil {
		return fmt.Errorf("[ERROR]: Error Generating Private Key !! %v", err);
	}

	// Save Private Key to PEM file
	privateKeyFile, err := os.Create(privateKeyPath);
	if err != nil {
		return fmt.Errorf("[ERROR]: Error in Cereating File for Private Key !! %v", err);
	}
	defer privateKeyFile.Close();

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return fmt.Errorf("[ERROR]: Error in Writing Private Key !! %v", err)
	}

	// Generate Public Key
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("[ERROR]: Error in Marshalling Public Key !! %v", err)
	}

	// Save Public Key to PEM file
	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		return fmt.Errorf("[ERROR]: Error in Creating File for Public Key !! %v", err)
	}

	defer publicKeyFile.Close()

	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		return fmt.Errorf("[ERROR]: Error in Writing Public Key !! %v", err)
	}

	fmt.Println("RSA Key Pair Generated Successfully..")
	return nil
}
