package apihandler

import (
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations/transaction"
	handler "github.com/alexwbaule/give-help/v2/handlers/transaction"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func AddTransactionHandler(rt *runtimeApp.Runtime) transaction.AddTransactionHandler {
	return &addTransaction{rt: rt}
}

type addTransaction struct {
	rt *runtimeApp.Runtime
}

func (ctx *addTransaction) Handle(params transaction.AddTransactionParams, principal *models.LoggedUser) middleware.Responder {

	t := handler.New(ctx.rt.GetDatabase())
	tid, err := t.Insert(params.Body)

	if err != nil {
		return transaction.NewAddTransactionInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return transaction.NewAddTransactionOK().WithPayload(tid)
}

func ChangeTransactionStatusHandler(rt *runtimeApp.Runtime) transaction.ChangeTransactionStatusHandler {
	return &changeTransactionStatus{rt: rt}
}

type changeTransactionStatus struct {
	rt *runtimeApp.Runtime
}

func (ctx *changeTransactionStatus) Handle(params transaction.ChangeTransactionStatusParams, principal *models.LoggedUser) middleware.Responder {

	t := handler.New(ctx.rt.GetDatabase())
	err := t.ChangeTransactionStatus(params.TransactionID, params.Body)

	if err != nil {
		return transaction.NewChangeTransactionStatusInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return transaction.NewChangeTransactionStatusOK()
}

func GetTransactionByIDHandler(rt *runtimeApp.Runtime) transaction.GetTransactionByIDHandler {
	return &getTransactionByID{rt: rt}
}

type getTransactionByID struct {
	rt *runtimeApp.Runtime
}

func (ctx *getTransactionByID) Handle(params transaction.GetTransactionByIDParams, principal *models.LoggedUser) middleware.Responder {

	t := handler.New(ctx.rt.GetDatabase())
	transactions, err := t.Load(params.TransactionID)

	if err != nil {
		return transaction.NewGetTransactionByIDInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return transaction.NewGetTransactionByIDOK().WithPayload(transactions)
}

func GetTransactionByProposalIDHandler(rt *runtimeApp.Runtime) transaction.GetTransactionByProposalIDHandler {
	return &getTransactionByProposalID{rt: rt}
}

type getTransactionByProposalID struct {
	rt *runtimeApp.Runtime
}

func (ctx *getTransactionByProposalID) Handle(params transaction.GetTransactionByProposalIDParams, principal *models.LoggedUser) middleware.Responder {

	t := handler.New(ctx.rt.GetDatabase())
	transactions, err := t.LoadByProposalID(params.ProposalID)

	if err != nil {
		return transaction.NewGetTransactionByProposalIDInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return transaction.NewGetTransactionByProposalIDOK().WithPayload(transactions)
}

func GetTransactionByUserIDHandler(rt *runtimeApp.Runtime) transaction.GetTransactionByUserIDHandler {
	return &getTransactionByUserID{rt: rt}
}

type getTransactionByUserID struct {
	rt *runtimeApp.Runtime
}

func (ctx *getTransactionByUserID) Handle(params transaction.GetTransactionByUserIDParams, principal *models.LoggedUser) middleware.Responder {

	t := handler.New(ctx.rt.GetDatabase())
	transactions, err := t.LoadByUserID(*principal.UserID)

	if err != nil {
		return transaction.NewGetTransactionByUserIDInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return transaction.NewGetTransactionByUserIDOK().WithPayload(transactions)
}

func TransactionGiverReviewHandler(rt *runtimeApp.Runtime) transaction.TransactionGiverReviewHandler {
	return &transactionReviewGiver{rt: rt}
}

type transactionReviewGiver struct {
	rt *runtimeApp.Runtime
}

func (ctx *transactionReviewGiver) Handle(params transaction.TransactionGiverReviewParams, principal *models.LoggedUser) middleware.Responder {

	t := handler.New(ctx.rt.GetDatabase())
	err := t.InsertGiverReview(params.TransactionID, params.Body)

	if err != nil {
		return transaction.NewTransactionGiverReviewInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return transaction.NewTransactionGiverReviewOK()
}

func TransactionTakerReviewHandler(rt *runtimeApp.Runtime) transaction.TransactionTakerReviewHandler {
	return &transactionReviewTaker{rt: rt}
}

type transactionReviewTaker struct {
	rt *runtimeApp.Runtime
}

func (ctx *transactionReviewTaker) Handle(params transaction.TransactionTakerReviewParams, principal *models.LoggedUser) middleware.Responder {

	t := handler.New(ctx.rt.GetDatabase())
	err := t.InsertTakerReview(params.TransactionID, params.Body)

	if err != nil {
		return transaction.NewTransactionTakerReviewInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return transaction.NewTransactionTakerReviewOK()
}

func TransactionAcceptTransactionHandler(rt *runtimeApp.Runtime) transaction.AcceptTransactionHandler {
	return &transactionAccept{rt: rt}
}

type transactionAccept struct {
	rt *runtimeApp.Runtime
}

func (ctx *transactionAccept) Handle(params transaction.AcceptTransactionParams, principal *models.LoggedUser) middleware.Responder {

	t := handler.New(ctx.rt.GetDatabase())
	err := t.Accept(params.TransactionID)

	if err != nil {
		return transaction.NewAcceptTransactionInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return transaction.NewAcceptTransactionOK()
}

func TransactionFinishTransactionHandler(rt *runtimeApp.Runtime) transaction.FinishTransactionHandler {
	return &transactionFinish{rt: rt}
}

type transactionFinish struct {
	rt *runtimeApp.Runtime
}

func (ctx *transactionFinish) Handle(params transaction.FinishTransactionParams, principal *models.LoggedUser) middleware.Responder {

	t := handler.New(ctx.rt.GetDatabase())
	err := t.Finish(params.TransactionID)

	if err != nil {
		return transaction.NewFinishTransactionInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return transaction.NewFinishTransactionOK()
}

func TransactionCancelTransactionHandler(rt *runtimeApp.Runtime) transaction.CancelTransactionHandler {
	return &transactionCancel{rt: rt}
}

type transactionCancel struct {
	rt *runtimeApp.Runtime
}

func (ctx *transactionCancel) Handle(params transaction.CancelTransactionParams, principal *models.LoggedUser) middleware.Responder {

	t := handler.New(ctx.rt.GetDatabase())
	err := t.Cancel(params.TransactionID, params.UserID)

	if err != nil {
		return transaction.NewChangeTransactionStatusInternalServerError().WithPayload(&models.APIError{Message: "An unexpected error occurred"})
	}

	return transaction.NewCancelTransactionOK()
}
