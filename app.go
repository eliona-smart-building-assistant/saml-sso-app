//  This file is part of the eliona project.
//  Copyright © 2023 LEICOM iTEC AG. All Rights Reserved.
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

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	LOG_REGIO       = "app"
	API_SERVER_PORT = 3000
)

func run() {

	var (
		err      error
		metadata []byte
	)

	basicConfig, err := conf.GetBasicConfig(context.Background())
	if err != nil {
		log.Fatal(LOG_REGIO, "cannot load basic config")
	}

	advancedConfig, err := conf.GetAdvancedConfig(context.Background())
	if err != nil {
		log.Fatal(LOG_REGIO, "cannot load advanced config")
	}

	if basicConfig.IdpMetadataUrl != nil && *basicConfig.IdpMetadataUrl != "" {
		// fetch metadata
		metadataResp, err := http.Get(*basicConfig.IdpMetadataUrl)
		if err != nil {
			log.Error(LOG_REGIO, "cannot fetch IdP metadata from url: %v", err)
		}
		defer metadataResp.Body.Close()
		metaB, err := io.ReadAll(metadataResp.Body)
		if err != nil {
			log.Error(LOG_REGIO, "cannot read metadata response from IdP: %v", err)
		}
		metadata = metaB
	} else if basicConfig.IdpMetadataXml != nil {
		metadata = []byte(*basicConfig.IdpMetadataXml)
	} else {
		log.Error(LOG_REGIO, "not able to set IdP Metadata")
	}

	sp, err := saml.NewServiceProviderAdvanced(basicConfig.ServiceProviderCertificate,
		basicConfig.ServiceProviderPrivateKey, basicConfig.OwnUrl,
		[]byte(metadata), &advancedConfig.EntityId,
		&advancedConfig.AllowInitializationByIdp, &advancedConfig.SignedRequest,
		&advancedConfig.ForceAuthn, &advancedConfig.CookieSecure)
	if err != nil {
		log.Fatal(LOG_REGIO, "cannot initialize saml service provider")
	}

	elionaAuth := eliona.NewAuthorization(basicConfig.OwnUrl, basicConfig.UserToArchive, advancedConfig.LoginFailedUrl)

	app := http.HandlerFunc(elionaAuth.Authorize)
	http.Handle("/sso/evaluate", sp.GetMiddleWare().RequireAccount(app))
	http.Handle("/saml/", sp.GetMiddleWare())

	err = http.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", strconv.Itoa(API_SERVER_PORT)), apiserver.NewRouter(
		apiserver.NewConfigurationApiController(apiservices.NewConfigurationApiService()),
		apiserver.NewVersionApiController(apiservices.NewVersionApiService()),
		apiserver.NewGenericSingleSignOnApiController(apiservices.NewGenericSingleSignOnApiService()),
		// apiserver.NewSAML20ApiController(apiservices.NewSAML20ApiService()), // managed over thirdparty lib crewjam/saml
	))
	log.Fatal(LOG_REGIO, "API server: %v", err)
}
