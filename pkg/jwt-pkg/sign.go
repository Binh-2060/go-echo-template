package jwtpkg

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//source: https://github.com/golang-jwt/jwt/blob/main/http_example_test.go

var privateKey []byte

type CustomClaims struct {
	jwt.RegisteredClaims
	//assign your model cliams data
}

const (
	privKeyPath = "app.rsa" // openssl genrsa -out app.rsa keysize
	// pubKeyPath  = "test/sample_key.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

// assign model
func SignToken() (string, error) {
	var err error
	if privateKey == nil {
		privateKey, err = os.ReadFile(privKeyPath)
		if err != nil {
			return "", err
		}
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &CustomClaims{
		jwt.RegisteredClaims{
			//assign token expired
			ExpiresAt: jwt.NewNumericDate(time.Now().Add((2 * time.Hour))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "Go-Echo-App",
			Subject:   "",
			ID:        "",
			Audience:  []string{"Go-Echo-App"},
		},
	}

	return t.SignedString(key)
}
