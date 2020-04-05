package authentication

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"

	app "github.com/alexwbaule/go-app"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	expiresAt  time.Duration
)

// CustomClaims Auth Claims
type CustomClaims struct {
	*jwt.StandardClaims
	TokenType string
	User      models.LoggedUser
	Scopes    map[string]bool
}

/*
InitToken initializes the token by weither loading the keys from the
filesystem with the loadToken() function or by generating temporarily
ones with the generateToken() function
*/
func InitToken(app app.Application) {
	var err error
	privateKey, publicKey, err = loadToken(app.Config().GetString("auth.Token.PrivateKey"), app.Config().GetString("auth.Token.PublicKey"))
	if err != nil {
		app.Logger().Fatalln(err)
	}
	expiresAt = app.Config().GetDuration("auth.Token.ExpiresAt")
}

// loadToken loads a private and public RSA keys from the filesystem in order to be used for the JWT signature
func loadToken(PrivateKey, PublicKey string) (*rsa.PrivateKey, *rsa.PublicKey, error) {

	if PrivateKey == "" || PublicKey == "" {
		return nil, nil, fmt.Errorf("The paths to the private and public RSA keys were not provided")
	}

	// Read the files from the filesystem
	prv, err := ioutil.ReadFile(PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to open the private key file: %v", err)
	}
	pub, err := ioutil.ReadFile(PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to open the public key file: %v", err)
	}

	// Parse the RSA keys
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(prv)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to parse the private key: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to parse the public key: %v", err)
	}

	return privateKey, publicKey, nil
}

// VerifyJWT extracts and verifies the validity of the JWT
func VerifyJWT(tokenString string) (*CustomClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return publicKey, nil
	})

	if err == nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, err
}
