package cryptoutils

import (
	"fmt"
	"goth/utils/testutils"
	"testing"
)

func TestValidationFullCycle(t *testing.T) {
	key := []byte("blah")
	privateKey, err := GenerateRSAKey()
	publicKey := privateKey.Public()
	tokenStr, err := JWTCreateToken(1, "qa4", "api-gateway", key, "HS256")
	if err != nil {
		t.Errorf("could not create the token: %s", err.Error())
		return
	}
	tokenStr2, err := JWTCreateToken(1, "qa4", "api-gateway", privateKey, "RS512")
	if err != nil {
		t.Errorf("could not create the token :%s", err.Error())
		return
	}
	claims, err := JWTValidateAndGetClaims(&tokenStr, key)
	if err != nil || claims == nil {
		t.Errorf("Could not validate the token: %s", err.Error())
		return
	}
	claims2, err := JWTValidateAndGetClaims(&tokenStr2, publicKey)
	if err != nil || claims2 == nil {
		t.Errorf("Could not validate the token: %s", err.Error())
		return
	}

	testutils.AssertEqualsInt64(t, 1, claims.UserId)
	testutils.AssertEqualsString(t, "qa4", claims.TenantId)
	testutils.AssertEqualsString(t, "api-gateway", claims.Originator)
	testutils.AssertEqualsInt64(t, 1, claims2.UserId)
	testutils.AssertEqualsString(t, "qa4", claims2.TenantId)
	testutils.AssertEqualsString(t, "api-gateway", claims2.Originator)
}

func TestJWTCreateErrors(t *testing.T) {
	key := []byte("blah")

	_, err := JWTCreateToken(1, "tenantId", "originator", key, "badSigningMethod")
	if err == nil {
		t.Error("should have seen create token error")
	}
	_, err = JWTCreateToken(1, "tenantId", "originator", nil, "RS512")
	if err == nil {
		t.Error("should have seen a key nil error")
	}
	_, err = JWTCreateToken(1, "tenantId", "originator", key, "RSx512")
	if err == nil {
		t.Error("should have seen a key nil error")
	}
}

func TestJWTValidationWithWrongKey(t *testing.T) {
	key1 := []byte("blah")
	key2 := []byte("halb")
	privateKey1, _ := GenerateRSAKey()

	privateKey2, _ := GenerateRSAKey()
	publicKey2 := privateKey2.Public()

	tokenStr1, _ := JWTCreateToken(1, "qa4", "api-gateway", key1, "HS256")
	tokenStr2, _ := JWTCreateToken(1, "qa4", "api-gateway", privateKey1, "RS512")

	_, err := JWTValidateAndGetClaims(&tokenStr1, key2)
	if err == nil {
		t.Error("should have had a validation error")
	}

	_, err = JWTValidateAndGetClaims(&tokenStr2, publicKey2)
	if err == nil {
		t.Error("should have had a validation error")
	}
}

func TestJWTValidationErrors(t *testing.T) {
	key := []byte("blah")
	_, err := JWTValidateAndGetClaims(nil, key)
	if err == nil {
		t.Error("should have had a tokenstring does not exist error")
	}
	ts := ""
	_, err = JWTValidateAndGetClaims(&ts, key)
	if err == nil {
		t.Error("should have had a token string does not exist error")
	}
	//nil key
	ts, _ = JWTCreateToken(1, "qa4", "api-gateway", key, "HS256")
	_, err = JWTValidateAndGetClaims(&ts, nil)
	if err == nil {
		t.Error("should have had a JWT Validation key does not exist error")
	}
}

//For testing the api with the JWTToken
func TestJWTCreate(t *testing.T) {

	//For testing the api
	ts, _ := JWTCreateToken(2, "sym-ac6", "api-gateway", []byte("blah"), "HS256")
	fmt.Println(ts)
}
