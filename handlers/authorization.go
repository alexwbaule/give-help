package handlers

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/swag"
)

// CheckAPIKeyAuth from Token
func CheckAPIKeyAuth(rt *runtimeApp.Runtime, token string, roles []string) (*models.LoggedUser, error) {
	var user *models.LoggedUser

	//jwttoken, err := authentication.VerifyJWT(token)

	user = &models.LoggedUser{
		Email: swag.String("email"),
		Name:  swag.String("name"),
	}

	return user, nil
}
