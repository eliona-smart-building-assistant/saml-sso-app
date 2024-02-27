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
	"fmt"
	"net/http"
	"os"
	"saml-sso/apiserver"
	"saml-sso/conf"
	"saml-sso/utils"

	"github.com/crewjam/saml/samlsp"
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	LOG_REGIO = "eliona"

	DefaultLang = "en"
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

	w.Header().Add("Content-Type", "application/json")
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

		loginEmail, userIp  string
		firstname, lastname string
		phone               *string

		user       *api.User
		jwt        *string
		setCookies http.Cookie

		errorMessage []byte
	)

	// Try to obtain real user IP.
	userIp = r.Header.Get("X-Forwarded-For")
	if userIp == "" {
		userIp = r.Header.Get("X-Real-Ip")
	}
	if userIp == "" {
		userIp = r.RemoteAddr
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
		phoneS := samlsp.AttributeFromContext(r.Context(), *mapping.Phone)
		phone = &phoneS
	}

	log.Info(LOG_REGIO, "User with firstname: %v, lastname: %v, email/login: "+
		"%v, phone: %v from %v want to login",
		firstname, lastname, loginEmail, phone, userIp)

	// get or create user
	user, err = s.eliApi.GetUserIfExists(loginEmail)
	if err != nil {
		log.Info(LOG_REGIO, "user doesn't exist. creating now user...")

		projectId, err := s.getProjectId()
		if err != nil {
			errorMessage = []byte(err.Error())
			goto internalServerError
		}

		sysRoleId, projRoleId, lang, err := s.getPermissionsAndLang(r.Context())
		if err != nil {
			log.Error(LOG_REGIO, "mapping failed: sysRoleId:%v, projRoleId:%v, lang:%v",
				sysRoleId, projRoleId, lang)
			goto notAuthenticated
		}

		// cannot set role over api
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

		err = UpdateElionaUserArchivedPhone(user.Email, phone, s.userToArchive)
		if err != nil {
			log.Error(LOG_REGIO, "cannot set phone and archive flag: %v", err)
			errorMessage = []byte(err.Error())
			goto internalServerError
		}

		err = SetUserPermissions(user.Id.Get(), sysRoleId, lang)
		if err != nil {
			log.Error(LOG_REGIO, "cannot set user permissions: %v", err)
			goto notAuthenticated
		}
		err = SetProjectUser(projectId, user.Id.Get(), projRoleId)
		if err != nil {
			errorMessage = []byte(err.Error())
			goto internalServerError
		}
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
	log.Warn(LOG_REGIO, "internal servererror occured while auth: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(errorMessage)
	if err != nil {
		log.Error(LOG_REGIO, "write internal server error: %v", err)
	}
}

func (s *SingleSignOn) getPermissionsAndLang(samlCtx context.Context) (sysRoleId int,
	projRoleId int, lang string, err error) {

	var (
		permissions *apiserver.Permissions
		aclRoleMap  map[string]int
	)

	sysRoleId, projRoleId = -1, -1
	lang = DefaultLang

	permissions, err = conf.GetPermissionMapping(context.Background())
	if err != nil {
		return
	}
	if permissions == nil {
		err = errors.New("cannot load permission config. <nil>")
		return
	}

	aclRoleMap, err = GetACLRoleMap()
	if err != nil {
		return
	}

	// get defaults
	sysRoleId = conf.StringToRoleId(permissions.DefaultSystemRole, aclRoleMap)

	projRoleId = conf.StringToRoleId(permissions.DefaultProjRole, aclRoleMap)

	lang = samlsp.AttributeFromContext(samlCtx, permissions.DefaultLanguage)

	// if configured, map permissions and lang
	if permissions.SystemRoleSamlAttribute != nil &&
		*permissions.SystemRoleSamlAttribute != "" &&
		permissions.SystemRoleMap != nil {

		systemRoleMap := conf.ApiRoleMapToGolangMap(*permissions.SystemRoleMap)

		samlValue := samlsp.AttributeFromContext(samlCtx, *permissions.SystemRoleSamlAttribute)

		elionaRoleOrId := systemRoleMap[samlValue]
		sysRoleId = conf.AnyToRoleId(elionaRoleOrId, aclRoleMap)
	}
	if permissions.ProjRoleSamlAttribute != nil &&
		*permissions.ProjRoleSamlAttribute != "" &&
		permissions.ProjRoleMap != nil {

		projectRoleMap := conf.ApiRoleMapToGolangMap(*permissions.ProjRoleMap)

		samlValue := samlsp.AttributeFromContext(samlCtx, *permissions.ProjRoleSamlAttribute)

		elionaRoleOrId := projectRoleMap[samlValue]
		projRoleId = conf.AnyToRoleId(elionaRoleOrId, aclRoleMap)
	}
	if permissions.LanguageSamlAttribute != nil &&
		*permissions.LanguageSamlAttribute != "" &&
		permissions.LanguageMap != nil {

		langMap := conf.ApiRoleMapToGolangMap(*permissions.LanguageMap)

		samlValue := samlsp.AttributeFromContext(samlCtx, *permissions.LanguageSamlAttribute)

		elionaLang := langMap[samlValue]
		switch l := elionaLang.(type) {
		case string:
			lang = l
		default:
			log.Warn(LOG_REGIO, "language after map invalid type: %T, %v", lang, lang)
		}
	}

	if sysRoleId <= 0 || projRoleId <= 0 || lang == "" {
		err = errors.New("mapping unsuccessful")
	}

	return
}

func (s *SingleSignOn) getProjectId() (projectId string, err error) {

	projectId = os.Getenv("PROJID")
	if projectId == "" {
		projectId, err = GetFirstProjectId()
	}
	if err != nil {
		err = fmt.Errorf("cannot look up project id: %v", err)
	}

	return
}
