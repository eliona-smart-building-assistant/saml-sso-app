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
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"saml-sso/utils"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/logger"
	"github.com/crewjam/saml/samlidp"
	"github.com/crewjam/saml/samlsp"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"github.com/zenazn/goji"
	"golang.org/x/crypto/bcrypt"
)

type ServiceProviderEntry struct {
	EntityID string
	MetaData string
}

// just for testing purposes!
type IdentityProvider struct {
	idp                     *samlidp.Server
	allowedServiceProviders []ServiceProviderEntry
}

const (
	IDENTITY_PROVIDER_METADATA_PATH = "/metadata"
	IDENTITY_PROVIDER_PORT          = 8000
)

func NewIdentityProvider(baseUrl string) (*IdentityProvider, error) {
	cn := "localhost"
	ip := "127.0.0.1"

	sCert, pKey, err := utils.CreateSelfsignedX509Certificate(10, 4096, &cn, &ip)
	if err != nil {
		return nil, err
	}

	certB, _ := pem.Decode([]byte(sCert))
	cert, err := x509.ParseCertificate(certB.Bytes)
	if err != nil {
		return nil, err
	}

	keyB, _ := pem.Decode([]byte(pKey))
	key, _ := x509.ParsePKCS1PrivateKey(keyB.Bytes)
	if err != nil {
		return nil, err
	}

	baseURLstr := flag.String("idp", baseUrl, "IdP Server IP or URL")
	flag.Parse()
	baseURL, err := url.Parse(*baseURLstr)
	if err != nil {
		return nil, err
	}

	idp, err := samlidp.New(samlidp.Options{
		URL:         *baseURL,
		Key:         key,
		Logger:      logger.DefaultLogger,
		Certificate: cert,
		Store:       &samlidp.MemoryStore{},
	})

	return &IdentityProvider{
		idp:                     idp,
		allowedServiceProviders: []ServiceProviderEntry{},
	}, err
}

func (i *IdentityProvider) AddUser(user string, password string, email string, name string, lastname string) error {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return i.idp.Store.Put("/users/"+user, samlidp.User{
		Name:           user,
		HashedPassword: hashedPassword,
		Groups:         []string{"Administrators", "Users"},
		Email:          email,
		CommonName:     name + " " + lastname,
		Surname:        lastname,
		GivenName:      name,
	})
}

func (i *IdentityProvider) AddServiceProvider(entityID string, metadata string) error {
	i.idp.IDP.ServiceProviderProvider = i

	i.allowedServiceProviders = append(i.allowedServiceProviders, ServiceProviderEntry{
		EntityID: entityID,
		MetaData: metadata,
	})

	return nil
}

func (i *IdentityProvider) ServeForever() {
	goji.Handle("/*", i.idp)
	goji.Serve()
}

func (i *IdentityProvider) GetServiceProvider(r *http.Request, serviceProviderID string) (*saml.EntityDescriptor, error) {

	for _, sp := range i.allowedServiceProviders {
		if serviceProviderID == sp.EntityID {
			log.Info("Identity Provider", "Service Provider allowed")
			fmt.Println(string(sp.MetaData))
			return samlsp.ParseMetadata([]byte(sp.MetaData))
		}
	}

	return nil, errors.New("Unknown Service Provider")
}
