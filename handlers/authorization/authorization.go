package authorization

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

func VerifyIDTokenAndCheckRevoked(ctx context.Context, app *firebase.App, idToken string) (*auth.Token, bool) {
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
