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

package saml

import (
	"crypto/rsa"
	"net/url"
	"saml-sso/utils"

	"github.com/crewjam/saml/samlsp"
)

type ServiceProvider struct {
	sp *samlsp.Middleware
}

func NewServiceProvider(certificate string, privateKey string, baseUrl string, idpMetadata []byte) (*ServiceProvider, error) {
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
		return nil, err
	}

	serviceProvider.sp, err = samlsp.New(samlsp.Options{
		URL:         *rootUrl,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMeta,
	})

	return &serviceProvider, err
}

func (s *ServiceProvider) GetMiddleWare() *samlsp.Middleware {
	return s.sp
}