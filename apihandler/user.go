package apihandler

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations/user"
	handler "github.com/alexwbaule/give-help/v2/handlers/user"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func AddUserHandler(rt *runtimeApp.Runtime) user.AddUserHandler {
	return &addUser{rt: rt}
}

type addUser struct {
	rt *runtimeApp.Runtime
}

func (ctx *addUser) Handle(params user.AddUserParams, principal *models.LoggedUser) middleware.Responder {
	params.Body.RegisterFrom = *principal.Provider
	params.Body.Name = *principal.Name

	if params.Body.Contact == nil {
		contact := &models.Contact{
			Email: *principal.Email,
		}
		params.Body.Contact = contact
	} else {
		params.Body.Contact.Email = *principal.Email
	}

	c := handler.New(ctx.rt.GetDatabase())
	ruser, err := c.Insert(params.Body, *principal.UserID)

	if err != nil {
		return user.NewAddUserInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return user.NewAddUserOK().WithPayload(ruser)
}
