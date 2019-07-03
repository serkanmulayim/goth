package cryptoutils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const (
	USER_ID_LITERAL = "UserID"
	MEMBER_LITERAL  = "Member"
	AES_IV_SIZE     = 16
	AES_AAD_SIZE    = 16
	AES_TAG_SIZE    = 16
	AES_KEY_SIZE    = 32
)

func GetUserIdHash(userId int64) string {
	s := fmt.Sprintf("%s%d", USER_ID_LITERAL, userId)
	//t := append(inputBytes, []byte(userId_str)...)
	x := sha256.Sum256([]byte(s))
	out := append(x[:], 1, 1)
	return base64.StdEncoding.EncodeToString(out)
}

//TODO: write the test and compare
func GetAESKeyObjectHash(securityContextHashB64 string, principalHashB64 string) (string, error) {
	scb, err := base64.StdEncoding.DecodeString(securityContextHashB64)
	if err != nil {
		return "", fmt.Errorf("Error in decoding security context hash: %s", err.Error())
	}
	phb, err := base64.StdEncoding.DecodeString(principalHashB64)
	if err != nil {
		return "", fmt.Errorf("Error in decoding pricipal hash: %s", err.Error())
	}
	memberByte := []byte(MEMBER_LITERAL)

	h := make([]byte, 0, len(memberByte)+len(scb)+len(phb))
	h = append(h, memberByte...)
	h = append(h, scb...)
	h = append(h, phb...)

	x := sha256.Sum256(h)
	out := append(x[:], 1, 1)

	return base64.StdEncoding.EncodeToString(out), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	rand.Read(b)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
