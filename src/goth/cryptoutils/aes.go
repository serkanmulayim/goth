package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

func DecryptAES256GCM(base64Cipher *string, key []byte) ([]byte, error) {
	if base64Cipher == nil {
		return nil, errors.New("cipherText cannot be nil")
	}

	c, err := base64.StdEncoding.DecodeString(*base64Cipher)
	if err != nil {
		return nil, err
	}

	if len(key) != AES_KEY_SIZE {
		s := fmt.Sprintf("AES key size is %d. It should be %d", len(key), AES_KEY_SIZE)
		err = errors.New(s)
		return nil, err
	}
	ciphertext_len := len(c)

	if ciphertext_len < 48 { //minimum ciphertext size
		s := fmt.Sprintf("Ciphertext format is smaller than minimum (48), length:%d", ciphertext_len)
		err = errors.New(s)
		return nil, err
	}

	iv := c[0:AES_IV_SIZE]
	aad := c[AES_IV_SIZE:(AES_IV_SIZE + AES_AAD_SIZE)]
	tag := c[(ciphertext_len - AES_TAG_SIZE):ciphertext_len]
	ciphertext := c[(AES_IV_SIZE + AES_AAD_SIZE):(ciphertext_len - AES_TAG_SIZE)]

	return decryptAES(ciphertext, key, iv, aad, tag)
}

func decryptAES(ciphertext []byte, key []byte, iv []byte, aad []byte, tag []byte) ([]byte, error) {
	c := append(ciphertext, (tag)...)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, AES_IV_SIZE)
	if err != nil {
		return nil, err
	}

	return gcm.Open(c[:0], iv, c, aad)
}

func EncryptAES256GCM(plaintext *string, key []byte, aad []byte) (string, error) {

	if plaintext == nil {
		return "", errors.New("plaintext cannot be nil")
	}

	if len(key) != AES_KEY_SIZE {
		s := fmt.Sprintf("AES key size is %d. It should be %d", len(key), AES_KEY_SIZE)
		err := errors.New(s)
		return "", err
	}

	plaintext_bytes := []byte(*plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv, err := GenerateRandomBytes(AES_IV_SIZE)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, AES_IV_SIZE)
	if err != nil {
		return "", err
	}

	ct := gcm.Seal(plaintext_bytes[:0], iv, plaintext_bytes, aad)
	tag := ct[(len(ct) - AES_TAG_SIZE):len(ct)]
	ciphertext := ct[:(len(ct) - AES_TAG_SIZE)]

	out := append(iv, aad...)
	out = append(out, ciphertext...)
	out = append(out, tag...)

	outStr := base64.StdEncoding.EncodeToString(out)
	return outStr, nil

}
