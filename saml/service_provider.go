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
	"net/url"
	"saml-sso/utils"

	"github.com/crewjam/saml/samlsp"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	LOG_REGIO = "service provider"
)

type ServiceProvider struct {
	sp *samlsp.Middleware
}

func NewServiceProvider(certificate string, privateKey string, baseUrl string,
	idpMetadata []byte) (*ServiceProvider, error) {

	return NewServiceProviderAdvanced(certificate, privateKey, baseUrl, idpMetadata, nil, nil,
		nil, nil, nil)
}

func NewServiceProviderAdvanced(certificate string, privateKey string, baseUrl string, idpMetadata []byte,
	entityId *string, allowInitByIdp *bool, signedRequest *bool, forceAuthn *bool,
	cookieSecure *bool) (*ServiceProvider, error) {
	var serviceProvider ServiceProvider = ServiceProvider{}

	rootUrl, err := url.Parse(baseUrl)
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
			"continiue without, but cannot operate with a IdP in current state! ... %v", err)
		// return nil, err
	}

	opts := samlsp.Options{
		URL:         *rootUrl,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMeta,
	}

	if entityId != nil {
		opts.EntityID = utils.SubstituteOwnUrlUrlString(*entityId, baseUrl)
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
	if cookieSecure != nil {
		// opts.CookieSecure: true // option not available any more
		log.Debug(LOG_REGIO, "not implemented")
	}

	serviceProvider.sp, err = samlsp.New(opts)

	return &serviceProvider, err
}

func (s *ServiceProvider) GetMiddleWare() *samlsp.Middleware {
	return s.sp
}
