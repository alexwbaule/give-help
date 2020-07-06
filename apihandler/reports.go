package apihandler

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations/reports"
	handler "github.com/alexwbaule/give-help/v2/handlers/reports"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func ReportsGetProposalReportHandler(rt *runtimeApp.Runtime) reports.GetProposalReportHandler {
	return &getProposalReportHandler{rt: rt}
}

type getProposalReportHandler struct {
	rt *runtimeApp.Runtime
}

func (ctx *getProposalReportHandler) Handle(params reports.GetProposalReportParams) middleware.Responder {

	c := handler.New(ctx.rt.GetDatabase())
	ret, err := c.LoadViews()

	if err != nil {
		return reports.NewGetProposalReportInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return reports.NewGetProposalReportOK().WithPayload(ret)
}

func ReportsGetProposalReportcsvHandler(rt *runtimeApp.Runtime) reports.GetProposalReportcsvHandler {
	return &getProposalReportcsvHandler{rt: rt}
}

type getProposalReportcsvHandler struct {
	rt *runtimeApp.Runtime
}

func (ctx *getProposalReportcsvHandler) Handle(params reports.GetProposalReportcsvParams) middleware.Responder {

	c := handler.New(ctx.rt.GetDatabase())
	ret, err := c.LoadViewsCSV()

	if err != nil {
		return reports.NewGetProposalReportcsvInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return reports.NewGetProposalReportcsvOK().WithPayload(ret)
}
