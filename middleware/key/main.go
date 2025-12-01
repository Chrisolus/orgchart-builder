package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Failed to generate private key:", err)
	}

	privatePEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := os.WriteFile("key/private.pem", pem.EncodeToMemory(privatePEM), 0600); err != nil {
		log.Fatal("Failed to write private key:", err)
	}
	log.Println("Private Key generated: private.pem (KEEP SECRET)")

	publicASN1 := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	publicPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicASN1,
	}
	if err := os.WriteFile("key/public.pem", pem.EncodeToMemory(publicPEM), 0644); err != nil {
		log.Fatal("Failed to write public key:", err)
	}
	log.Println("Public Key generated: public.pem (SHARE FREELY)")
}
