package cryptoutils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"goth/utils/testutils"
	"testing"
)

const (
	USERID_9999_HASH = "Vxw/ycsW+m+H18EasOszVhPUNZksOwb1lV8MuiZULZcBAQ=="

	AES_DECRYPT_KEY    = "fcbc7eb62716dc7f792b6194d26d6d569eaee07a9d3c37ca42854090661e1845"
	AES_DECRYPT_IV     = "4c8c4624279b23b495c788844c76d225"
	AES_DECRYPT_AAD    = "3c182af19c46ff4acbdacecf70b42fb5"
	AES_DECRYPT_TAG    = "9a50233fb7e4af0b7d82ebddee6fc0ca"
	AES_DECRYPT_INPUT  = "40d682b5796e642069f101c22199c0b1"
	AES_DECRYPT_OUTPUT = "22144fc12f7bc5522b88b76c8ded1c76"
)

func TestGetUserIdHash(t *testing.T) {
	x := GetUserIdHash(int64(9999))
	testutils.AssertEqualsString(t, x, USERID_9999_HASH)
}

func TestDecrypt(t *testing.T) {

	iv, _ := hex.DecodeString(AES_DECRYPT_IV)
	key, _ := hex.DecodeString(AES_DECRYPT_KEY)
	aad, _ := hex.DecodeString(AES_DECRYPT_AAD)
	tag, _ := hex.DecodeString(AES_DECRYPT_TAG)
	ct, _ := hex.DecodeString(AES_DECRYPT_INPUT)
	expected, _ := hex.DecodeString(AES_DECRYPT_OUTPUT)

	out, _ := decryptAES(ct, key, iv, aad, tag)

	outStr := hex.EncodeToString(out)
	expectedStr := hex.EncodeToString(expected)
	testutils.AssertEqualsString(t, outStr, expectedStr)
}

func TestDecryptAES256_GCM(t *testing.T) {
	iv, _ := hex.DecodeString(AES_DECRYPT_IV)
	key, _ := hex.DecodeString(AES_DECRYPT_KEY)
	aad, _ := hex.DecodeString(AES_DECRYPT_AAD)
	tag, _ := hex.DecodeString(AES_DECRYPT_TAG)
	ct, _ := hex.DecodeString(AES_DECRYPT_INPUT)
	expected, _ := hex.DecodeString(AES_DECRYPT_OUTPUT)

	ciphertext := append(iv, aad...)
	ciphertext = append(ciphertext, ct...)
	ciphertext = append(ciphertext, tag...)

	ciphertextStr := base64.StdEncoding.EncodeToString(ciphertext)

	out, _ := DecryptAES256GCM(&ciphertextStr, key)

	outStr := hex.EncodeToString(out)
	expectedStr := hex.EncodeToString(expected)
	testutils.AssertEqualsString(t, outStr, expectedStr)
}

func TestFull(t *testing.T) {
	var plaintext string
	aad := make([]byte, AES_AAD_SIZE) //all zeros
	key := make([]byte, AES_KEY_SIZE) //all zeros

	for i := 0; i < 66; i++ {
		if i > 0 {
			plaintext = fmt.Sprintf("%s%d", plaintext, i)
		}
		out, err := EncryptAES256GCM(&plaintext, key, aad)

		if err != nil {
			t.Errorf("%s", err.Error())
		}

		decrypted, err := DecryptAES256GCM(&out, key)
		testutils.AssertEqualsString(t, string(decrypted), plaintext)
	}
}

func TestDecryptionError(t *testing.T) {
	plaintext := "Encrypt me man..."
	key := make([]byte, AES_KEY_SIZE)
	aad := make([]byte, AES_AAD_SIZE)

	_, err := EncryptAES256GCM(nil, key, aad)
	if err == nil {
		t.Error("should have been a \"plaintext cannot be nil error\"")
	}
	cipher, _ := EncryptAES256GCM(&plaintext, key, aad)
	key2 := make([]byte, AES_KEY_SIZE)
	key2[0] = 1

	_, err = DecryptAES256GCM(&cipher, key2)
	if err == nil {
		t.Error("There should have been a decryption error")
	}
}
