package apihandler

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations/category"
	handler "github.com/alexwbaule/give-help/v2/handlers/category"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func GetCategoryHandler(rt *runtimeApp.Runtime) category.GetCategoryHandler {
	return &getCategory{rt: rt}
}

type getCategory struct {
	rt *runtimeApp.Runtime
}

func (ctx *getCategory) Handle(params category.GetCategoryParams, principal *models.LoggedUser) middleware.Responder {
	c := handler.New(ctx.rt.GetDatabase())
	categories, err := c.Load()

	if err != nil {
		return category.NewGetCategoryInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return category.NewGetCategoryOK().WithPayload(categories)
}
