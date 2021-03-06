package apihandler

import (
	"time"

	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations/proposal"
	handler "github.com/alexwbaule/give-help/v2/handlers/proposal"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func AddProposalHandler(rt *runtimeApp.Runtime) proposal.AddProposalHandler {
	return &addProposal{rt: rt}
}

type addProposal struct {
	rt *runtimeApp.Runtime
}

func (ctx *addProposal) Handle(params proposal.AddProposalParams, principal *models.LoggedUser) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	pid, err := p.Insert(params.Body)

	if err != nil {
		return proposal.NewAddProposalInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewAddProposalOK().WithPayload(pid)
}

func AddProposalImagesHandler(rt *runtimeApp.Runtime) proposal.AddProposalImagesHandler {
	return &addProposalImages{rt: rt}
}

type addProposalImages struct {
	rt *runtimeApp.Runtime
}

func (ctx *addProposalImages) Handle(params proposal.AddProposalImagesParams, principal *models.LoggedUser) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	err := p.AddImages(params.ProposalID, params.Body)

	if err != nil {
		return proposal.NewAddProposalImagesInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewAddProposalImagesOK()
}

func ChangeProposalImagesHandler(rt *runtimeApp.Runtime) proposal.ChangeProposalImagesHandler {
	return &changeProposalImages{rt: rt}
}

type changeProposalImages struct {
	rt *runtimeApp.Runtime
}

func (ctx *changeProposalImages) Handle(params proposal.ChangeProposalImagesParams, principal *models.LoggedUser) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	err := p.ChangeImages(params.ProposalID, params.Body)

	if err != nil {
		return proposal.NewChangeProposalImagesInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewChangeProposalImagesOK()
}

func AddProposalTagsHandler(rt *runtimeApp.Runtime) proposal.AddProposalTagsHandler {
	return &addProposalTag{rt: rt}
}

type addProposalTag struct {
	rt *runtimeApp.Runtime
}

func (ctx *addProposalTag) Handle(params proposal.AddProposalTagsParams, principal *models.LoggedUser) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	err := p.AddTags(params.ProposalID, params.Body)

	if err != nil {
		return proposal.NewAddProposalTagsInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewAddProposalTagsOK()
}

func ChangeProposalStateHandler(rt *runtimeApp.Runtime) proposal.ChangeProposalStateHandler {
	return &changeProposalState{rt: rt}
}

type changeProposalState struct {
	rt *runtimeApp.Runtime
}

func (ctx *changeProposalState) Handle(params proposal.ChangeProposalStateParams, principal *models.LoggedUser) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	err := p.ChangeValidStatus(params.ProposalID, params.ProposalState)

	if err != nil {
		return proposal.NewChangeProposalStateInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewChangeProposalStateOK()
}

func ChangeProposalTextHandler(rt *runtimeApp.Runtime) proposal.ChangeProposalTextHandler {
	return &changeProposalText{rt: rt}
}

type changeProposalText struct {
	rt *runtimeApp.Runtime
}

func (ctx *changeProposalText) Handle(params proposal.ChangeProposalTextParams, principal *models.LoggedUser) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	err := p.ChangeText(params.ProposalID, *params.Body.Title, *params.Body.Description)

	if err != nil {
		return proposal.NewChangeProposalStateInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewChangeProposalStateOK()
}

func ChangeProposalValidateHandler(rt *runtimeApp.Runtime) proposal.ChangeProposalValidateHandler {
	return &changeProposalValidDate{rt: rt}
}

type changeProposalValidDate struct {
	rt *runtimeApp.Runtime
}

func (ctx *changeProposalValidDate) Handle(params proposal.ChangeProposalValidateParams, principal *models.LoggedUser) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	err := p.ChangeValidate(params.ProposalID, time.Time(*params.Body.Date))

	if err != nil {
		return proposal.NewChangeProposalValidateInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewChangeProposalValidateOK()
}

func SearchProposalsHandler(rt *runtimeApp.Runtime) proposal.SearchProposalsHandler {
	return &searchProposalsHandler{rt: rt}
}

type searchProposalsHandler struct {
	rt *runtimeApp.Runtime
}

func (ctx *searchProposalsHandler) Handle(params proposal.SearchProposalsParams) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	result, err := p.LoadFromFilter(params.Body)

	if err != nil {
		return proposal.NewSearchProposalsInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewSearchProposalsOK().WithPayload(result)
}

func GetProposalByIDHandler(rt *runtimeApp.Runtime) proposal.GetProposalByIDHandler {
	return &getProposalByID{rt: rt}
}

type getProposalByID struct {
	rt *runtimeApp.Runtime
}

func (ctx *getProposalByID) Handle(params proposal.GetProposalByIDParams) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	oneProposal, err := p.LoadFromID(params.ProposalID)

	if err != nil {
		return proposal.NewGetProposalByIDInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewGetProposalByIDOK().WithPayload(oneProposal)
}

func GetProposalByUserIDHandler(rt *runtimeApp.Runtime) proposal.GetProposalByUserIDHandler {
	return &getProposalByUser{rt: rt}
}

type getProposalByUser struct {
	rt *runtimeApp.Runtime
}

func (ctx *getProposalByUser) Handle(params proposal.GetProposalByUserIDParams, principal *models.LoggedUser) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	proposals, err := p.LoadFromUser(*principal.UserID)

	if err != nil {
		return proposal.NewGetProposalByUserIDInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewGetProposalByUserIDOK().WithPayload(proposals)
}

func GetProposalShareDataIDHandler(rt *runtimeApp.Runtime) proposal.GetProposalShareDataIDHandler {
	return &getProposalShareData{rt: rt}
}

type getProposalShareData struct {
	rt *runtimeApp.Runtime
}

func (ctx *getProposalShareData) Handle(params proposal.GetProposalShareDataIDParams) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	shareData, err := p.GetUserDataToShare(params.ProposalID)

	if err != nil {
		return proposal.NewGetProposalShareDataIDInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewGetProposalShareDataIDOK().WithPayload(shareData)
}

func AddProposalComplaintHandler(rt *runtimeApp.Runtime) proposal.AddProposalComplaintHandler {
	return &addProposalComplaintHandler{rt: rt}
}

type addProposalComplaintHandler struct {
	rt *runtimeApp.Runtime
}

func (ctx *addProposalComplaintHandler) Handle(params proposal.AddProposalComplaintParams) middleware.Responder {

	p := handler.New(ctx.rt.GetDatabase(), ctx.rt.GetCache())
	err := p.InsertComplaint(params.Body)

	if err != nil {
		return proposal.NewAddProposalComplaintInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return proposal.NewAddProposalComplaintOK()
}
