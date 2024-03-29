//  This file is part of the eliona project.
//  Copyright © 2023 Eliona by IoTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"context"
	"io"
	"net/http"
	"saml-sso/apiserver"
	"saml-sso/apiservices"
	"saml-sso/conf"
	"saml-sso/eliona"
	"saml-sso/saml"
	"strconv"

	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	LOG_REGIO       = "app"
	API_SERVER_PORT = 3000
	SSO_SERVER_PORT = 8081 // Publicly accessible without auth. See wiki.

	SAML_SPECIFIC_ENDPOINT_PATH = "/saml/"
)

func initialize() {
	ctx := context.Background()

	// Necessary to close used init resources
	conn := db.NewInitConnectionWithContextAndApplicationName(ctx, app.AppName())
	defer conn.Close(ctx)

	// Init the app before the first run.
	app.Init(conn, app.AppName(),
		app.ExecSqlFile("conf/init.sql"),
	)
}

func run() {

	var (
		err      error
		metadata []byte
	)

	err = conf.InsertAutoSamlConfiguration(context.Background())
	if err != nil {
		log.Debug(LOG_REGIO, "insert default config: %v", err)
	}

	config, err := conf.GetConfig(context.Background())
	if err != nil {
		log.Fatal(LOG_REGIO, "cannot load basic config")
	}

	if config.IdpMetadataUrl != nil && *config.IdpMetadataUrl != "" {
		// fetch metadata
		metadataResp, err := http.Get(*config.IdpMetadataUrl)
		if err != nil {
			log.Error(LOG_REGIO, "cannot fetch IdP metadata from url: %v", err)
		} else {
			defer metadataResp.Body.Close()
			metaB, err := io.ReadAll(metadataResp.Body)
			if err != nil {
				log.Error(LOG_REGIO, "cannot read metadata response from IdP: %v", err)
			}
			metadata = metaB
		}
	} else if config.IdpMetadataXml != nil {
		metadata = []byte(*config.IdpMetadataXml)
	} else {
		log.Error(LOG_REGIO, "not able to set IdP Metadata")
	}

	apiPort := common.Getenv("API_SERVER_PORT", strconv.Itoa(API_SERVER_PORT))
	samlSpPort := common.Getenv("SSO_SERVER_PORT", strconv.Itoa(SSO_SERVER_PORT))

	log.Debug(LOG_REGIO, "own url: %v, api port: %v", config.OwnUrl, apiPort)

	sp, err := saml.NewServiceProviderAdvanced(
		config.ServiceProviderCertificate,
		config.ServiceProviderPrivateKey,
		config.OwnUrl,
		metadata,
		&config.EntityId,
		&config.AllowInitializationByIdp,
		&config.SignedRequest,
		&config.ForceAuthn,
		&config.CookieSecure,
	)
	if err != nil {
		log.Fatal(LOG_REGIO, "cannot initialize saml service provider: %v", err)
	}

	// app api handle to router
	router := apiserver.NewRouter(
		apiserver.NewConfigurationAPIController(apiservices.NewConfigurationApiService()),
		apiserver.NewVersionAPIController(apiservices.NewVersionApiService()),

		// separate handled in eliona/sso_handles.go (no RESTful)
		// apiserver.NewGenericSingleSignOnApiController(apiservices.NewGenericSingleSignOnApiService()),

		// managed over thirdparty lib crewjam/saml (no RESTful)
		// apiserver.NewSAML20ApiController(apiservices.NewSAML20ApiService()),
	)

	go func() {
		err := http.ListenAndServe(":"+apiPort, router)
		if err != nil {
			log.Fatal(LOG_REGIO, "app api server: %v", err)
		}
	}()

	// saml specific handle (no RESTful) to router
	elionaAuth := eliona.NewSingleSignOn(config.OwnUrl,
		config.UserToArchive, config.LoginFailedUrl)

	activeHandleFunc := http.HandlerFunc(elionaAuth.ActiveHandle)
	http.Handle(eliona.ENDPOINT_SSO_GENERIC_ACTIVE, activeHandleFunc)
	authHandleFunc := http.HandlerFunc(elionaAuth.Authentication) // TODO: Not completely implemented.
	http.Handle(eliona.ENDPOINT_SSO_GENERIC_VERIFICATION,
		sp.GetMiddleWare().RequireAccount(authHandleFunc))
	http.Handle(SAML_SPECIFIC_ENDPOINT_PATH, sp.GetMiddleWare())

	// for backwards compatibility, can be removed when the frontend is reworked to the new generic /sso/* endpoints
	http.Handle("/adfs/active/", activeHandleFunc)
	http.Handle("/adfs/auth/",
		sp.GetMiddleWare().RequireAccount(authHandleFunc))

	log.Info(LOG_REGIO, "started @ %v", samlSpPort)
	err = http.ListenAndServe(":"+samlSpPort, nil)
	if err != nil {
		log.Error("sp app", "exiting due to an error: %v", err)
	} else {
		log.Info("sp app", "exited")
	}

	log.Fatal(LOG_REGIO, "API server: %v", err)
}
