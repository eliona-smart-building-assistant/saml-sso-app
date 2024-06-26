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

package saml

import (
	"crypto/rsa"
	"net/http"
	"net/url"
	"saml-sso/utils"

	"github.com/crewjam/saml/samlsp"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	LOG_REGIO           = "service provider"
	SP_HANDLE_BASE_PATH = "/saml/" // used from saml middleware for ACS, ACL, etc.
	PUBLIC_BASE_PATH    = "/apps-public/saml-sso"
)

type ServiceProvider struct {
	pubBasePath string
	sp          *samlsp.Middleware
}

func NewServiceProvider(certificate string, privateKey string, pubBaseUrl string,
	idpMetadata []byte) (*ServiceProvider, error) {

	return NewServiceProviderAdvanced(certificate, privateKey, pubBaseUrl, idpMetadata, nil, nil,
		nil, nil, "")
}

func NewServiceProviderAdvanced(certificate string, privateKey string, baseUrl string, idpMetadata []byte,
	entityId *string, allowInitByIdp *bool, signedRequest *bool, forceAuthn *bool, pubBasePath string,
) (*ServiceProvider, error) {
	var serviceProvider ServiceProvider = ServiceProvider{
		pubBasePath: pubBasePath,
	}

	rootUrl, err := url.Parse(baseUrl + pubBasePath + "/")
	if err != nil {
		return nil, err
	}

	keyPair, err := utils.GetCombinedX509Certificate(certificate, privateKey)
	if err != nil {
		return nil, err
	}

	idpMeta, err := samlsp.ParseMetadata(idpMetadata)
	if err != nil {
		log.Warn(LOG_REGIO, "cannot parse metadata. "+
			"continuing without, but cannot operate with a IdP in current state! ... %v", err)
		// return nil, err
	}

	opts := samlsp.Options{
		URL:                *rootUrl,
		Key:                keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:        keyPair.Leaf,
		IDPMetadata:        idpMeta,
		DefaultRedirectURI: PUBLIC_BASE_PATH + "/",
	}

	if entityId != nil {
		opts.EntityID, _ = utils.SubstituteOwnUrlUrlString(*entityId, baseUrl)
	}
	if allowInitByIdp != nil {
		opts.AllowIDPInitiated = *allowInitByIdp
	}
	if signedRequest != nil {
		opts.SignRequest = *signedRequest
	}
	if forceAuthn != nil {
		opts.ForceAuthn = *forceAuthn
	}

	serviceProvider.sp, err = samlsp.New(opts)

	log.Debug(LOG_REGIO, "ACS URL %v", serviceProvider.sp.ServiceProvider.AcsURL)
	log.Info(LOG_REGIO, "Metadata URL %v", serviceProvider.sp.ServiceProvider.MetadataURL)
	log.Debug(LOG_REGIO, "SLO URL %v", serviceProvider.sp.ServiceProvider.SloURL)

	return &serviceProvider, err
}

func (s *ServiceProvider) GetMiddleWare() *samlsp.Middleware {

	return s.sp
}

func (s *ServiceProvider) FixPath(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.URL.Path = s.pubBasePath + r.URL.Path

		next.ServeHTTP(w, r)
	})
}
