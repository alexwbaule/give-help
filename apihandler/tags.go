package apihandler

import (
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations/tags"
	handler "github.com/alexwbaule/give-help/v2/handlers/tags"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rafaelfino/metrics"
)

func GetTagsHandler(rt *runtimeApp.Runtime) tags.GetTagsHandler {
	return &getTags{rt: rt}
}

type getTags struct {
	rt *runtimeApp.Runtime
}

func (ctx *getTags) Handle(params tags.GetTagsParams) middleware.Responder {
	start := time.Now()

	c := handler.New(ctx.rt.GetDatabase())
	ret, err := c.Load()

	defer ctx.rt.GetMetricProcessor().Send(metrics.NewMetric("GetTagsHandler.ElapsedTime", metrics.CounterType, nil, float64(time.Since(start).Milliseconds())))

	if err != nil {
		return tags.NewGetTagsInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return tags.NewGetTagsOK().WithPayload(ret)
}
