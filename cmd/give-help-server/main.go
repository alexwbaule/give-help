// MIT License
//
// Copyright (c) 2020 Alex W. Baulé
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package main

import (
	"log"

	"github.com/alexwbaule/give-help/v2/apihandler"
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/generated/restapi"
	"github.com/alexwbaule/give-help/v2/generated/restapi/operations"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"

	app "github.com/alexwbaule/go-app"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/flagext"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/cors"
)

func main() {

	app, err := app.New("give-help-service")
	cfg := app.Config()

	cfg.SetDefault("service.Host", "127.0.0.1")
	cfg.SetDefault("service.Port", "8081")
	cfg.SetDefault("service.TLSWriteTimeout", "15m")
	cfg.SetDefault("service.WriteTimeout", "15m")

	rt, err := runtimeApp.NewRuntime(app)
	if err != nil {
		log.Fatal(err.Error())
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err.Error())
	}

	api := operations.NewGiveHelpAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()
	defer rt.CloseDatabase()

	server.EnabledListeners = app.Config().GetStringSlice("service.EnabledListeners")
	server.Host = app.Config().GetString("service.Host")
	server.Port = app.Config().GetInt("service.Port")
	server.ListenLimit = app.Config().GetInt("service.ListenLimit")
	server.TLSHost = app.Config().GetString("service.TLSHost")
	server.TLSPort = app.Config().GetInt("service.TLSPort")
	server.TLSListenLimit = app.Config().GetInt("service.TLSListenLimit")

	server.CleanupTimeout = app.Config().GetDuration("service.CleanupTimeout")
	server.TLSKeepAlive = app.Config().GetDuration("service.TLSKeepAlive")
	server.TLSReadTimeout = app.Config().GetDuration("service.TLSReadTimeout")
	server.TLSWriteTimeout = app.Config().GetDuration("service.TLSWriteTimeout")
	server.KeepAlive = app.Config().GetDuration("service.KeepAlive")
	server.ReadTimeout = app.Config().GetDuration("service.ReadTimeout")
	server.WriteTimeout = app.Config().GetDuration("service.WriteTimeout")
	server.MaxHeaderSize = flagext.ByteSize(app.Config().GetSizeInBytes("service.MaxHeaderSize"))

	server.SocketPath = flags.Filename(app.Config().GetString("service.SocketPath"))
	server.TLSCertificate = flags.Filename(app.Config().GetString("service.TLSCertificate"))
	server.TLSCertificateKey = flags.Filename(app.Config().GetString("service.TLSCertificateKey"))
	server.TLSCACertificate = flags.Filename(app.Config().GetString("service.TLSCACertificate"))

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "give-help-service"
	parser.LongDescription = swaggerSpec.Spec().Info.Description

	/*
	 * App Handlers
	 */

	// Applies when the "x-api-token" header is set
	api.APIKeyHeaderAuth = func(token string, roles []string) (*models.LoggedUser, error) {
		return apihandler.CheckAPIKeyAuth(rt, token, roles)
	}

	/* API Proposal */
	api.ProposalAddProposalHandler = apihandler.AddProposalHandler(rt)
	api.ProposalAddProposalImagesHandler = apihandler.AddProposalImagesHandler(rt)
	api.ProposalChangeProposalImagesHandler = apihandler.ChangeProposalImagesHandler(rt)
	api.ProposalAddProposalTagsHandler = apihandler.AddProposalTagsHandler(rt)
	api.ProposalChangeProposalStateHandler = apihandler.ChangeProposalStateHandler(rt)
	api.ProposalChangeProposalTextHandler = apihandler.ChangeProposalTextHandler(rt)
	api.ProposalChangeProposalValidateHandler = apihandler.ChangeProposalValidateHandler(rt)
	api.ProposalSearchProposalsHandler = apihandler.SearchProposalsHandler(rt)
	api.ProposalGetProposalByIDHandler = apihandler.GetProposalByIDHandler(rt)
	api.ProposalGetProposalByUserIDHandler = apihandler.GetProposalByUserIDHandler(rt)
	api.ProposalGetProposalShareDataIDHandler = apihandler.GetProposalShareDataIDHandler(rt)
	api.ProposalAddProposalComplaintHandler = apihandler.AddProposalComplaintHandler(rt)

	/* API User */
	api.UserAddUserHandler = apihandler.AddUserHandler(rt)
	api.UserUpdateUserByIDHandler = apihandler.UpdateUserByIDHandler(rt)
	api.UserGetUserByIDHandler = apihandler.GetUserByIDHandler(rt)

	/* API Transaction */
	api.TransactionAddTransactionHandler = apihandler.AddTransactionHandler(rt)
	api.TransactionChangeTransactionStatusHandler = apihandler.ChangeTransactionStatusHandler(rt)
	api.TransactionGetTransactionByIDHandler = apihandler.GetTransactionByIDHandler(rt)
	api.TransactionGetTransactionByProposalIDHandler = apihandler.GetTransactionByProposalIDHandler(rt)
	api.TransactionGetTransactionByUserIDHandler = apihandler.GetTransactionByUserIDHandler(rt)
	api.TransactionTransactionGiverReviewHandler = apihandler.TransactionGiverReviewHandler(rt)
	api.TransactionTransactionTakerReviewHandler = apihandler.TransactionTakerReviewHandler(rt)

	api.TransactionAcceptTransactionHandler = apihandler.TransactionAcceptTransactionHandler(rt)
	api.TransactionFinishTransactionHandler = apihandler.TransactionFinishTransactionHandler(rt)
	api.TransactionCancelTransactionHandler = apihandler.TransactionCancelTransactionHandler(rt)

	/* API Tags */
	api.TagsGetTagsHandler = apihandler.GetTagsHandler(rt)

	/* API Banks */
	api.BanksGetBankListHandler = apihandler.GetBankListHandler(rt)

	/* API Etc */
	api.EtcGetEtcListHandler = apihandler.GetEtcListHandler(rt)

	/* API Terms */
	api.TermsPutUserAcceptHandler = apihandler.TermsPutUserAcceptHandler(rt)
	api.TermsGetTermsHandler = apihandler.TermsGetTermsHandler(rt)
	api.TermsGetUserAcceptedHandler = apihandler.TermsGetUserAcceptedHandler(rt)

	/* API Reports */

	c := cors.New(cors.Options{
		Debug:              true,
		AllowedHeaders:     []string{"*"},
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "OPTIONS", "DELETE", "PATCH"},
		MaxAge:             1000,
		OptionsPassthrough: false,
	})

	handler := c.Handler(api.Serve(nil))
	server.SetHandler(handler)

	if err := server.Serve(); err != nil {
		log.Fatal(err.Error())
	}
}
