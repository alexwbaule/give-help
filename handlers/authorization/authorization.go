package authorization

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
)

// CheckAPIKeyAuth from Token
func CheckAPIKeyAuth(rt *runtimeApp.Runtime, token string, roles []string) (*models.LoggedUser, error) {
	var user *models.LoggedUser

	//jwttoken, err := authentication.VerifyJWT(token)

	user = &models.LoggedUser{
		ID: models.ID(common.GetULID()),
	}

	return user, nil
}
