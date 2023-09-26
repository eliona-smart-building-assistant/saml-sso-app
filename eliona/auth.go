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

package eliona

import (
	"context"
	"errors"
	"net/http"
	"saml-sso/apiserver"
	"saml-sso/conf"
	"saml-sso/utils"

	"github.com/crewjam/saml/samlsp"
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	LOG_REGIO = "eliona"
)

type Authorization struct {
	baseUrl         string
	redirectNoLogin string
	userToArchive   bool
	eliApi          *EliApiV2
}

func NewAuthorization(baseUrl string, userToArchive bool,
	redirectNoLogin string) *Authorization {

	return &Authorization{
		baseUrl:         baseUrl,
		userToArchive:   userToArchive,
		eliApi:          NewEliApiV2(),
		redirectNoLogin: utils.SubstituteOwnUrlUrlString(redirectNoLogin, baseUrl),
	}
}

// redirected from the generic sso endpoint sso/auth (handled by app api)
func (a *Authorization) Authorize(w http.ResponseWriter, r *http.Request) {
	log.Info(LOG_REGIO, "SAML Auth [%s]", r.Method)

	var (
		err error

		mapping *apiserver.AttributeMap

		loginEmail                 string = ""
		firstname, lastname, phone string

		user       *api.User
		jwt        *string
		setCookies http.Cookie

		errorMessage []byte
	)

	mapping, err = conf.GetAttributeMapping(context.Background())
	if err != nil {
		log.Error(LOG_REGIO, "cannot get attribute mapping. skip auth. %v", err)
		errorMessage = []byte(err.Error())
		goto internalServerError
	}

	if mapping.Email != "" {
		loginEmail = samlsp.AttributeFromContext(r.Context(), mapping.Email)
	} else {
		// default without config (fallback)
		loginEmail = samlsp.AttributeFromContext(r.Context(),
			"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn")
	}

	if mapping.FirstName != nil && *mapping.FirstName != "" {
		firstname = samlsp.AttributeFromContext(r.Context(), *mapping.FirstName)
	}
	if mapping.LastName != nil && *mapping.LastName != "" {
		lastname = samlsp.AttributeFromContext(r.Context(), *mapping.LastName)
	}
	if mapping.Phone != nil && *mapping.Phone != "" {
		phone = samlsp.AttributeFromContext(r.Context(), *mapping.Phone)
	}

	log.Info(LOG_REGIO, "User with firstname: %v, lastname: %v, email/login: "+
		"%v, phone: %v want to login",
		firstname, lastname, loginEmail, phone)

	// get or create user
	user, err = a.eliApi.GetUserIfExists(loginEmail)
	if err != nil {
		log.Info(LOG_REGIO, "user doesn't exist. creating now user...")
		user, err = a.eliApi.AddUser(&api.User{
			Email:     loginEmail,
			Firstname: *api.NewNullableString(&firstname),
			Lastname:  *api.NewNullableString(&lastname),
			// Phone: *api.NewNullableString(&phone), 	// not possible over APIv2
			// Archived: a.userToArchive,				// not possible over APIv2
		})
		if err != nil {
			log.Error(LOG_REGIO, "creating user: %v", err)
			errorMessage = []byte(err.Error())
			goto internalServerError
		}
	}

	err = UpdateElionaUserArchivedPhone(user.Email, &phone, a.userToArchive)
	if err != nil {
		log.Error(LOG_REGIO, "cannot set phone and archive flag: %v", err)
		errorMessage = []byte(err.Error())
		goto internalServerError
	}

	err = a.SetUserPermissions(user.Email)
	if err != nil {
		log.Error(LOG_REGIO, "cannot set user permissions")
	}

	// obtain a jwt to login via cookies
	jwt, err = GetElionaJsonWebToken(user.Email)
	if err != nil {
		log.Error(LOG_REGIO, "cannot obtain a JWT")
		errorMessage = []byte("cannot obtain a JWT")
		goto internalServerError
	}

	log.Debug(LOG_REGIO, "User %s with token %v", user.Email, jwt)
	if jwt == nil || *jwt == "" {
		goto notAuthorized
	}

	goto authorized

authorized:
	setCookies = http.Cookie{
		Name:  "elionaAuthorization",
		Value: *jwt,
		Path:  "/"}

	http.SetCookie(w, &setCookies)
	http.Redirect(w, r, a.baseUrl, http.StatusFound)
	return

notAuthorized:
	http.Redirect(w, r, a.redirectNoLogin, http.StatusFound)
	return

internalServerError:
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(errorMessage)
	return
}

func (a *Authorization) SetUserPermissions(email string) error {

	var (
		err         error
		permissions *apiserver.Permissions
	)

	permissions, err = conf.GetPermissionSettings(context.Background())
	if err != nil {
		return err
	}

	log.Info(LOG_REGIO,
		"ToDo: add user to a project and set permissions according the configurations. %v",
		permissions)

	err = errors.New("not implemented")
	return err
}
