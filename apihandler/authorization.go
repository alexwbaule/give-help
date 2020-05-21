package apihandler

import (
	"context"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/handlers/authorization"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// CheckAPIKeyAuth from Token
func CheckAPIKeyAuth(rt *runtimeApp.Runtime, tokenStr string, roles []string) (*models.LoggedUser, error) {
	var user *models.LoggedUser

	token, isRevoked := authorization.VerifyIDTokenAndCheckRevoked(context.Background(), rt.GetFirebase(), tokenStr)

	if isRevoked {
		return nil, errors.New(401, "Revoked token, please log in again.")
	}

	if token == nil {
		return nil, errors.New(401, "Invalid token, please log in again.")
	}

	name := swag.String(token.Claims["name"].(string))

	if name == nil {
		name = swag.String(token.Claims["email"].(string))
	}

	user = &models.LoggedUser{
		Email:    swag.String(token.Claims["email"].(string)),
		Name:     name,
		Picture:  swag.String(token.Claims["picture"].(string)),
		Provider: &token.Firebase.SignInProvider,
		UserID:   swag.String(token.Claims["user_id"].(string)),
	}

	return user, nil
}
