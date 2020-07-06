package apihandler

import (
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations/etc"
	handler "github.com/alexwbaule/give-help/v2/handlers/etc"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rafaelfino/metrics"
)

func GetEtcListHandler(rt *runtimeApp.Runtime) etc.GetEtcListHandler {
	return &getEtcListHandler{rt: rt}
}

type getEtcListHandler struct {
	rt *runtimeApp.Runtime
}

func (ctx *getEtcListHandler) Handle(params etc.GetEtcListParams) middleware.Responder {
	start := time.Now()

	c := handler.New(ctx.rt.GetDatabase())
	ret, err := c.Load()

	defer ctx.rt.GetMetricProcessor().Send(metrics.NewMetric("GetEtcListHandler.ElapsedTime", metrics.CounterType, nil, float64(time.Since(start).Milliseconds())))

	if err != nil {
		return etc.NewGetEtcListInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return etc.NewGetEtcListOK().WithPayload(ret)
}
