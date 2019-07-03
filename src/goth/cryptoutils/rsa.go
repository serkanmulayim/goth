package cryptoutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

const (
	RSA_KEY_SIZE           = 2048
	RSA_PEM_KEY_TYPE_PKCS1 = "RSA PRIVATE KEY"
	RSA_PEM_KEY_TYPE_PKCS8 = "PRIVATE KEY"
)

func GenerateRSAKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE)
}

func EncryptRSA(plaintext []byte, publicKey *rsa.PublicKey) (string, error) {
	label := []byte("")
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, label[:])
	if err != nil {
		return "", err
	}
	outB64 := base64.StdEncoding.EncodeToString(ciphertext)
	return outB64, nil
}

func DecryptRSA(ciphertextB64 *string, privateKey *rsa.PrivateKey) ([]byte, error) {
	if ciphertextB64 == nil {
		return nil, errors.New("RSA ciphertext cannot be nil")
	}
	label := []byte("")
	ct_bytes, err := base64.StdEncoding.DecodeString(*ciphertextB64)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ct_bytes, label)
}

//Reads from PKCS8 and PKCS1 formatted key
func PemToRSAPrivateKey(pemContent *string) (*rsa.PrivateKey, error) {
	if pemContent == nil {
		return nil, errors.New("RSA PEM private key string cannot be nil")
	}

	pemBytes := []byte(*pemContent)

	privPem, _ := pem.Decode(pemBytes)
	if privPem == nil {
		return nil, errors.New("PEM content is not decodable")
	}

	var outI interface{}
	outI, err := x509.ParsePKCS8PrivateKey(privPem.Bytes)
	if err != nil {
		outI, err = x509.ParsePKCS1PrivateKey(privPem.Bytes)
		if err != nil {
			return nil, errors.New("Private key is not parsable")
		}
	}

	out, ok := outI.(*rsa.PrivateKey)
	if !ok { //will not happen
		return nil, errors.New("This is not a Private Key")
	}

	return out, nil

}
