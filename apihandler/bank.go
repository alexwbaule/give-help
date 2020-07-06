package apihandler

import (
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations/banks"
	handler "github.com/alexwbaule/give-help/v2/handlers/bank"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rafaelfino/metrics"
)

func GetBankListHandler(rt *runtimeApp.Runtime) banks.GetBankListHandler {
	return &getBanks{rt: rt}
}

type getBanks struct {
	rt *runtimeApp.Runtime
}

func (ctx *getBanks) Handle(params banks.GetBankListParams) middleware.Responder {
	start := time.Now()

	c := handler.New(ctx.rt.GetDatabase())
	ret, err := c.LoadBanks()

	defer ctx.rt.GetMetricProcessor().Send(metrics.NewMetric("GetBankListHandler.ElapsedTime", metrics.CounterType, nil, float64(time.Since(start).Milliseconds())))

	if err != nil {
		return banks.NewGetBankListInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return banks.NewGetBankListOK().WithPayload(ret)
}
