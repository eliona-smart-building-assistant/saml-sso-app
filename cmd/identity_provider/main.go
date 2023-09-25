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
	"encoding/json"
	"os"
	"saml-sso/saml"

	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	IDP_CONFIG = "./cmd/identity_provider/cnf.json"
)

type IdpConfig struct {
	BaseURL          string `json:"baseUrl"`
	IdpCertificate   string `json:"idpCertificate"`
	IdpPrivateKey    string `json:"idpPrivateKey"`
	ServiceProviders []struct {
		EntityID string `json:"entityId"`
		Metadata string `json:"metadata"`
	} `json:"serviceProviders"`
	Users []struct {
		Login     string `json:"login"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Password  string `json:"password"`
	} `json:"users"`
}

func main() {
	var cnf IdpConfig = IdpConfig{}

	fileContent, err := os.ReadFile(IDP_CONFIG)
	if err != nil {
		log.Fatal("idp app", "reading config: %v", err)
	}

	err = json.Unmarshal(fileContent, &cnf)
	if err != nil {
		log.Fatal("idp app", "parsing config: %v", err)
	}

	idpCert, err := os.ReadFile(cnf.IdpCertificate)
	if err != nil {
		log.Fatal("idp app", "reading cert: %v", err)
	}
	idpCertStr := string(idpCert)

	idpKey, err := os.ReadFile(cnf.IdpPrivateKey)
	if err != nil {
		log.Fatal("idp app", "reading key: %v", err)
	}
	idpKeyStr := string(idpKey)

	idp, err := saml.NewIdentityProvider(cnf.BaseURL, &idpCertStr, &idpKeyStr)
	if err != nil {
		log.Fatal("idp app", "creating idp: %v", err)
	}

	if err != nil {
		log.Fatal("idp app", "reading config: %v", err)
	}

	for _, sp := range cnf.ServiceProviders {
		metaSp, err := os.ReadFile(sp.Metadata)
		if err != nil {
			log.Fatal("idp app", "read sp meta")
		}
		err = idp.AddServiceProvider(sp.EntityID, string(metaSp))
		if err != nil {
			log.Fatal("idp app", "adding service providers: %v", err)
		}
	}

	for _, user := range cnf.Users {
		err = idp.AddUser(user.Login, user.Password, user.Email, user.FirstName, user.LastName)
		if err != nil {
			log.Fatal("idp app", "adding user: %v", err)
		}
	}

	log.Info("idp app", "starting and serve identity provider @ :8000")
	idp.ServeForever()
	log.Warn("idp app", "exited")
}
