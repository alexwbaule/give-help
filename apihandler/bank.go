package apihandler

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations/banks"
	handler "github.com/alexwbaule/give-help/v2/handlers/bank"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func GetBankListHandler(rt *runtimeApp.Runtime) banks.GetBankListHandler {
	return &getBanks{rt: rt}
}

type getBanks struct {
	rt *runtimeApp.Runtime
}

func (ctx *getBanks) Handle(params banks.GetBankListParams) middleware.Responder {
	c := handler.New(ctx.rt.GetDatabase())
	ret, err := c.LoadBanks()

	if err != nil {
		return banks.NewGetBankListInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return banks.NewGetBankListOK().WithPayload(ret)
}
