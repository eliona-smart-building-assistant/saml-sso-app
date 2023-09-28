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
	"encoding/json"
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

const (
	ENDPOINT_SSO_GENERIC_VERIFICATION = "/sso/auth"
	ENDPOINT_SSO_GENERIC_ACTIVE       = "/sso/active"
)

type SingleSignOn struct {
	baseUrl         string
	redirectNoLogin string
	userToArchive   bool
	eliApi          *EliApiV2
}

func NewSingleSignOn(baseUrl string, userToArchive bool,
	redirectNoLogin string) *SingleSignOn {

	return &SingleSignOn{
		baseUrl:         baseUrl,
		userToArchive:   userToArchive,
		eliApi:          NewEliApiV2(),
		redirectNoLogin: utils.SubstituteOwnUrlUrlString(redirectNoLogin, baseUrl),
	}
}

func (s *SingleSignOn) ActiveHandle(w http.ResponseWriter, r *http.Request) {

	var (
		err          error
		responseCode int    = http.StatusMethodNotAllowed
		responseMsg  []byte = []byte("not allowed")
	)

	log.Debug(LOG_REGIO, "active handle called")

	if r.Method == http.MethodGet {
		active := apiserver.Active{
			Active: true,
		}
		responseMsg, err = json.Marshal(active)
		if err == nil {
			responseCode = http.StatusOK
		} else {
			responseCode = http.StatusInternalServerError
			responseMsg = []byte(err.Error())
		}
	}

	w.WriteHeader(responseCode)
	_, err = w.Write(responseMsg)
	if err != nil {
		log.Error(LOG_REGIO, "write internal server error: %v", err)
	}
}

func (s *SingleSignOn) Authentication(w http.ResponseWriter, r *http.Request) {
	log.Info(LOG_REGIO, "authentication handle called [%s]", r.Method)

	var (
		err error

		mapping *apiserver.AttributeMap

		loginEmail, userIp         string
		firstname, lastname, phone string

		user       *api.User
		jwt        *string
		setCookies http.Cookie

		errorMessage []byte
	)

	// try to optain real user ip
	userIp = r.Header.Get("X-Forwarded-For")
	if userIp == "" {
		userIp = r.Header.Get("X-Real-Ip")
	}
	log.Debug(LOG_REGIO, "user from %s called authentication ep", userIp)

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
	user, err = s.eliApi.GetUserIfExists(loginEmail)
	if err != nil {
		log.Info(LOG_REGIO, "user doesn't exist. creating now user...")
		user, err = s.eliApi.AddUser(&api.User{
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

	err = UpdateElionaUserArchivedPhone(user.Email, &phone, s.userToArchive)
	if err != nil {
		log.Error(LOG_REGIO, "cannot set phone and archive flag: %v", err)
		errorMessage = []byte(err.Error())
		goto internalServerError
	}

	err = s.setUserPermissions(user.Email)
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
		goto notAuthenticated
	}

	goto authenticated

authenticated:
	log.Info(LOG_REGIO, "authenticated user login: %s", loginEmail)
	setCookies = http.Cookie{
		Name:  "elionaAuthorization",
		Value: *jwt,
		Path:  "/"}

	http.SetCookie(w, &setCookies)
	http.Redirect(w, r, s.baseUrl, http.StatusFound)
	return

notAuthenticated:
	log.Info(LOG_REGIO, "not authenticated user tried to login: %s, %s",
		loginEmail, userIp)
	// reset eliona cookies
	setCookies = http.Cookie{
		Name:  "elionaAuthorization",
		Value: "invalid",
		Path:  "/"}

	http.SetCookie(w, &setCookies)
	http.Redirect(w, r, s.redirectNoLogin, http.StatusFound)
	return

internalServerError:
	log.Warn(LOG_REGIO, "internal servererror occured while auth")
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(errorMessage)
	if err != nil {
		log.Error(LOG_REGIO, "write internal server error: %v", err)
	}
}

func (s *SingleSignOn) setUserPermissions(email string) error {

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
