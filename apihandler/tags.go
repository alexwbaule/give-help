package apihandler

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations/tags"
	handler "github.com/alexwbaule/give-help/v2/handlers/tags"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func GetTagsHandler(rt *runtimeApp.Runtime) tags.GetTagsHandler {
	return &getTags{rt: rt}
}

type getTags struct {
	rt *runtimeApp.Runtime
}

func (ctx *getTags) Handle(params tags.GetTagsParams) middleware.Responder {

	c := handler.New(ctx.rt.GetDatabase())
	ret, err := c.Load()

	if err != nil {
		return tags.NewGetTagsInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return tags.NewGetTagsOK().WithPayload(ret)
}
