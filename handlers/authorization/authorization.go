package authorization

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/alexwbaule/give-help/v2/generated/models"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// CheckAPIKeyAuth from Token
func CheckAPIKeyAuth(rt *runtimeApp.Runtime, tokenStr string, roles []string) (*models.LoggedUser, error) {
	var user *models.LoggedUser

	token, isRevoked := verifyIDTokenAndCheckRevoked(context.Background(), rt.GetFirebase(), tokenStr)

	if isRevoked {
		return nil, errors.New(401, "Revoked token, please log in again.")
	}

	user = &models.LoggedUser{
		Email:    swag.String(token.Claims["email"].(string)),
		Name:     swag.String(token.Claims["name"].(string)),
		Picture:  swag.String(token.Claims["picture"].(string)),
		Provider: &token.Firebase.SignInProvider,
		UserID:   swag.String(token.Claims["user_id"].(string)),
	}

	return user, nil
}

func verifyIDTokenAndCheckRevoked(ctx context.Context, app *firebase.App, idToken string) (*auth.Token, bool) {
	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error getting Auth client: %v\n", err)
		return nil, false
	}
	token, err := client.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		if err.Error() == "ID token has been revoked" {
			return nil, true
		} else {
			log.Printf("error verifying ID token: %v\n", err)
			return nil, false
		}
	}
	return token, false
}
