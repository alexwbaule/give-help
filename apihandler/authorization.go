package apihandler

import (
	"context"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/handlers/authorization"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/errors"
)

// CheckAPIKeyAuth from Token
func CheckAPIKeyAuth(rt *runtimeApp.Runtime, tokenStr string, roles []string) (*models.LoggedUser, error) {
	var user *models.LoggedUser
	var name string
	var email string
	var picture string
	var userID string

	token, isRevoked := authorization.VerifyIDTokenAndCheckRevoked(context.Background(), rt.GetFirebase(), tokenStr)

	if isRevoked {
		return nil, errors.New(401, "Revoked token, please log in again.")
	}

	if token == nil {
		return nil, errors.New(401, "Invalid token, please log in again.")
	}

	if value, ok := token.Claims["email"]; ok {
		email = value.(string)
	}

	if value, ok := token.Claims["name"]; ok {
		name = value.(string)
	}

	if value, ok := token.Claims["picture"]; ok {
		picture = value.(string)
	}

	if value, ok := token.Claims["user_id"]; ok {
		userID = value.(string)
	} else {
		return nil, errors.New(401, "Invalid token, missing UserID please log in again.")
	}

	user = &models.LoggedUser{
		Email:    &email,
		Name:     &name,
		Picture:  &picture,
		Provider: &token.Firebase.SignInProvider,
		UserID:   &userID,
	}

	return user, nil
}
