package cryptoutils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	RS512                 = "RS512"
	RS256                 = "RS256"
	HS256                 = "HS256"
	TOKEN_EXPIRATION_TIME = 3600000
)

type APIGatewayClaims struct {
	TenantId          string `json:"tenantId,omitempty"`
	UserId            int64  `json:"userId"`
	Originator        string `json:"originator,omitempty"`
	CreationTime      int64  `json:"creationTime"`
	SesExpirationTime int64  `json:"sesExpirationTime"`
}

type JWTClaims struct {
	APIGatewayClaims
	jwt.StandardClaims
}

func (claims *JWTClaims) Valid() error {
	now := nowMillis()

	if claims.ExpiresAt < now && claims.ExpiresAt != 0 { //assumed as a non-required field
		return errors.New("Token valid time is expired, exp field")
	}
	if claims.SesExpirationTime < now { // assumed as a required field
		return errors.New("Token valid time is expired, sesExpirationTime field")
	}

	if claims.TenantId == "" {
		return errors.New("Token does not have tenantId in the claims")
	}

	if claims.UserId == 0 { //will never be 0
		return errors.New("Token does not have userId in the claims")
	}

	return nil
}

func JWTValidateAndGetClaims(tokenString *string, key interface{}) (*JWTClaims, error) {
	//jwt.StandardClaims
	if tokenString == nil || *tokenString == "" {
		return nil, errors.New("TokenString does not exist")
	}

	token, err := jwt.ParseWithClaims(*tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if key == nil {
			return nil, errors.New("JWT Validation key does not exist")
		}
		return key, nil

	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		claims, ok := token.Claims.(*JWTClaims)
		if ok {
			return claims, nil
		}
	}
	return nil, errors.New("Error in returning the claims")
}

func nowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func JWTCreateToken(userId int64, tenantId string, originator string, key interface{}, signingMethod string) (string, error) {

	creationTime := nowMillis()
	sesExpTime := creationTime + TOKEN_EXPIRATION_TIME

	var claims JWTClaims = JWTClaims{
		APIGatewayClaims{
			TenantId:          tenantId,
			UserId:            userId,
			Originator:        originator,
			CreationTime:      creationTime,
			SesExpirationTime: sesExpTime,
		},
		jwt.StandardClaims{
			ExpiresAt: sesExpTime,
			Issuer:    originator,
			Subject:   tenantId,
		},
	}
	if signingMethod == "" {
		signingMethod = RS512 // defaults to JWTUtils class in SymphonyLibrary
	}

	if signingMethod != HS256 && signingMethod != RS256 && signingMethod != RS512 {
		return "", errors.New("Signing method should be one of [HS256, RS256, RS512]")
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), &claims)

	return token.SignedString(key)

}
