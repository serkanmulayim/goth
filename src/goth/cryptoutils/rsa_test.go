package cryptoutils

import (
	"crypto/rsa"
	"goth/utils/testutils"
	"testing"
)

const (
	PRIVATE_RSA_KEY_PKCS1 = "-----BEGIN RSA PRIVATE KEY-----\n" +
		"MIIEpAIBAAKCAQEAvrddN/7oBwPoYCSc1DfRk6CejMt/5h3Vz3FqBds1bn+7gOWW\n" +
		"GJ+Z6XmQRxGzvkzbudM9sW9r0KwA7elRU4v5+GTggjRX3EwQUT29rxaZp5uloz60\n" +
		"cKVt8bgy814YQ//LrUSqTYN46AC0DykB9iTw2tgFbE+E5DGOfyUXh24iPElLv4eM\n" +
		"5nVqVNhkCbiE4uxkcTpchUPF9vbTzGHAJqAu28NJ/GIX9bfkuCU3IqBVDFg+87sm\n" +
		"WaWqFB5+qjrBW5omfLlW1bJmfyCItBKspIUNwu+UYj9VwcZhB7YXDl7/B8fPEcw0\n" +
		"9y3WcTfGS50Keq90MPJX1lBMbyMnSUEMlGUwVQIDAQABAoIBADps2AzodVJ42v5h\n" +
		"GP5WX0jXgtrlGLh5WV/kgbNrlTHVxa5WJyZB9pe02wM6pBLfXszNru+lj6TsnJhF\n" +
		"ytlKX3i+Bp08xdHCJ5mLzYlO5iGXqCWbdxGvTEApysoNgGeMfGwHRhja9vY0CPU/\n" +
		"/c7XQEX6uaaVscNqfxnqVgvyLGDFF1xvJ8tSkJt6cnCGGiap5k91hBQlDe+EE0H/\n" +
		"AfpdnESFENYWw0C85QO653kw0DsncUf11i/zi5myPhk9GyvGla2rQYXR7mzhWnk7\n" +
		"rFJ1fU6xEINqt6V8BxyBZ2WgcbbVp8sP03+KeEV9rJnSXkPfr8R7MVqwg3lF1P4v\n" +
		"P1iEVyECgYEA9eKC1QFNUKzsIZIqCGgMINR2E9MxaO+AL58rbQhbJJhMUoksmnA8\n" +
		"LmBHx6vhbDEBOax64VLwCSiM2PIqXVelfF6uCTkjjOS/EDBXYTl6tZVztOmntNrP\n" +
		"2SeDrOL31cAXk3ebKZ/DNQFHLPEllhoIgZRUPpvQZ1dzEHIbtxfhyFkCgYEAxo/b\n" +
		"KQPedvDojNfrye6XAq4feLuLQPUGhoRfwLHd4I28CKLjw/A0289n4daaeZE5GffJ\n" +
		"L+AKFWWLz169xbbkcIfw+Fhix50py7gzIiEGFq/LHgIVGL7TyikzurQyXj5mqjHO\n" +
		"A3tzdOYaa7yv8KUS+ly5KHekUnj5L9G5QB5tqF0CgYEA7Uv9F7R5+THpcTyudA52\n" +
		"JfjlO5zGQo9hFpR2RHOcAmU4wy6/bQyEB/3DSGAI3XEtzYs3y91O6ofH5ldPq2W6\n" +
		"v5xgOWZ7eK4J9oWwzpO6aSQi1qXoxBGoJdqMqX0h0ZfhKOB1TmZw7ead2jGgJIxb\n" +
		"STLtWr7lEdbDpUt6k+jobYECgYBlmzeKxXARbTWS3OrLakvcEz4HifC4TKoKK7LT\n" +
		"6DCht+WAhdoHaeOil3+RET/69VVC9Ij+9qYyTMQ+WTzC79P1wZqNeq0ReFrq5FdG\n" +
		"H5/9+/b0ZBxnjT7TNEJER/F+SHBApzvEjGeIajca6nMdsVsEwm39RzvO+BiLX8Zs\n" +
		"CSYUAQKBgQCalbMIyKFOdmTJnM5Gq0lQ6MWqOpDq+WwSN0N2QQ2Mb/WVjrcF4M9W\n" +
		"K9eq06Xqfi3S1laOyrQLMKrK55x2gEbZQ5Da82wAFUcAwpKX2YzibVaDty9nyScB\n" +
		"DUewUd5eiJxJYTubhibiPUyqJ+/33FXlPjeWnn+m4Y5RCcOalcSPQw==\n" +
		"-----END RSA PRIVATE KEY-----"

	PRIVATE_RSA_KEY_PKCS8 = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDa4iGqOnyyu7An
B+hQolua6W8t40grXbzx2RHl5vtmPGvYJ7PyzBheIk3SJCX39iKAGNo7sO9gfFv+
xG/4Obbxfl1KIAhKxC4XGnASvTZqzTV8rvlswuipEflk7BNFRVzNw1pr6QCtTsO9
cBFEINOfKLt6Hrv3kmhA/a6IltIAr2oF9/vG908BdV+3CIlNNPmZLsc54Wcps30L
Ka73OlJ+Vguj6IYBcFeQhKMr8hXelHTTF0CSQUOg13Vt3OQpwA/jdzFjXhs9V4az
9RPRIqGqBb2v2GrXmZMfZh9OWItrap414y6G3gGljeHq7IMElBkcnv37Jf4hq5/Y
cnM12ehvAgMBAAECggEANQ9iV9DLz8wd3KhXEcz55ei2seoJ4bFcj0guxYuj39m/
zut4/N+q/I7WTJ3EsoJxRJpBtPykWsBq0Kyv8/2RsqMOXFy9hdpezezzxCcoHOXY
FHyaLFpnaizYRzXPShiRcbPspIIcTubsHm0fUmOeyVNndygiKtbL/Q1CFQxypgVt
JfB9LiLKegVUoM0ol4GqFHL3XzA4yoz82WDwM/sEqJJpC1YWdxRzE8eYGmTMr55q
1dl460w5Fyhy9NaCZ6dUiyHiL6OH8sA7I7iU1BgQZGIP+XMDF4ujXvPqYqOO8MK0
D7Co8lFK967R1H3g1vZxK7/biimYnpnaDxaXgpDeWQKBgQDxDtzu7Hy0jT5kW3oS
Wm0u/6fR2seCTgJ5nxSLVpCBVcfvKT2S+1fRVgNMPrxZ9rLf5OWyGYc4oP2R94YZ
df9BGtycYZcrOuypNtjQQWuSQNLTdE/ISSi896OtBCAmB5uTA+gA4OWxvvBdjACf
hMMimtSkn5H3vSMseSEz3tcYUwKBgQDoc2W0GsohllzvfunjVc/Of1ds9zS1KjvT
zM0RoTlItp0Br52YG2ZNmdb8fwJYhnSYxv/2Q6lkIoTzaQmij19ZKVgfpEX0BGmk
bjWqoWjmgcPedhNNOY6+v6CmfV3Nj3ugHNgnG0IsVouS+yIjdSDWPnuMQssCP49z
FU/y3Hm79QKBgQDeTn2zncaX39ZVSQN334rnmAAViXUKl5SywuF4atmTXR+oUNkn
LsJbHL6n1wdu1BM8ZeTq7Z4FvHp83c/+tRI04WfolBuMU6gjmaAz1tE0rLGBLrfR
Fp8KPjrk+XQIfmWcHDMedEmANX2IV+/PLOmkhTNrqnk8BmJkxkS3iF/HXwKBgQCb
imQagPaSRPgI9cZxbUExLvqEGmJ1ez4vOlJaIqSfKDqlHyr31hW9hVxa9m3OaKHq
fPZXhez56TNHYRimYwNtOQIToiuA3dcGxQw6EemMnHZBDIdb3FDNCJLp9OdonkMd
308v08iSvJKGlm7AhSak1Yh8UVFgPsGxQyiNHMSEZQKBgGFXzsNQ6W83uNT8oV4+
XEvgOORjgM9NULHjeStNbV9BX0+Ez9r1T0Oo1yvQv/JMeIYKHCGBGxHLtL7VaBeA
ve85DpaY4bDQiW3zOghvC/oRQnP7s1rsWvDvqxkvSSUYffWID883B3DFniIHddc5
gOv+1f5L5wA/IGcVjNb7bnxa
-----END PRIVATE KEY-----`
)

func TestFullRSA_WithKeyGeneration(t *testing.T) {

	expected := "help meee, they are encrypting me"
	pt_bytes := []byte(expected)

	privateKey1, err := GenerateRSAKey()

	if err != nil {
		t.Errorf("error in generating RSA key")
	}

	//encrypt
	publicKey1 := privateKey1.Public().(*rsa.PublicKey)
	ciphertext1, err := EncryptRSA(pt_bytes, publicKey1)
	if err != nil {
		t.Errorf("Error in RSA encryptiong: %s", err)
	}

	//decrypt
	plaintext_out_bytes1, err := DecryptRSA(&ciphertext1, privateKey1)
	if err != nil {
		t.Errorf("Error in RSA decryption: %s", err)
	}
	plaintext1 := string(plaintext_out_bytes1)
	testutils.AssertEqualsString(t, plaintext1, expected)
}

func Test_PrivateKeyFromPEM_WRONGFORMAT(t *testing.T) {
	pkstr := "Do you think this string really looks like a private key?"

	_, err := PemToRSAPrivateKey(&pkstr)
	if err == nil {
		t.Errorf("Method should have failed since the key is not a valid key in PKCS1")
	}

}

func TestFullRSA_WithKeyFromPEM_PKCSx(t *testing.T) {

	expected := "help meee, they are encrypting me"
	pt_bytes := []byte(expected)

	pkstr1 := (PRIVATE_RSA_KEY_PKCS1)
	pkstr8 := (PRIVATE_RSA_KEY_PKCS8)

	privateKey1, err := PemToRSAPrivateKey(&pkstr1)
	if err != nil {
		t.Errorf("error in loading RSA key PKCS1, %s", err)
	}

	privateKey8, err := PemToRSAPrivateKey(&pkstr8)
	if err != nil {
		t.Errorf("error in loading RSA key PKCS8, %s", err)
	}

	//encrypt
	publicKey1 := privateKey1.Public().(*rsa.PublicKey)
	publicKey8 := privateKey8.Public().(*rsa.PublicKey)

	ciphertext1, err := EncryptRSA(pt_bytes, publicKey1)
	if err != nil {
		t.Errorf("Error in RSA encryption: %s", err)
	}

	ciphertext8, err := EncryptRSA(pt_bytes, publicKey8)
	if err != nil {
		t.Errorf("Error in RSA encryption: %s", err)
	}

	// //decrypt
	plaintext_out_bytes1, err := DecryptRSA(&ciphertext1, privateKey1)
	if err != nil {
		t.Errorf("Error in RSA decryption: %s", err)
	}

	plaintext_out_bytes8, err := DecryptRSA(&ciphertext8, privateKey8)
	if err != nil {
		t.Errorf("Error in RSA decryption: %s", err)
	}

	plaintext1 := string(plaintext_out_bytes1)
	plaintext8 := string(plaintext_out_bytes8)
	testutils.AssertEqualsString(t, plaintext1, expected)
	testutils.AssertEqualsString(t, plaintext8, expected)
}
