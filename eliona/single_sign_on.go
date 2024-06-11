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
	"net/url"
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
	htmlContent     bool
	userToArchive   bool
	eliApi          *EliApiV2
}

func NewSingleSignOn(baseUrl string, userToArchive bool,
	redirectNoLogin string) *SingleSignOn {

	sso := &SingleSignOn{
		baseUrl:       baseUrl,
		userToArchive: userToArchive,
		eliApi:        NewEliApiV2(),
	}

	sso.redirectNoLogin, sso.htmlContent =
		utils.SubstituteOwnUrlUrlString(redirectNoLogin, baseUrl)

	return sso
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

		user *api.User
		jwt  *string
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
		s.authFailed(true, loginEmail, userIp, fmt.Sprintf("cannot get attribute mapping: %v", err), w, r)
		return
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
			s.authFailed(true, loginEmail, userIp, fmt.Sprintf("cannot obtain project id: %v", err), w, r)
			return
		}

		sysRoleId, projRoleId, lang, err := s.getPermissionsAndLang(r.Context())
		if err != nil {
			log.Info(LOG_REGIO, "mapping failed: sysRoleId:%v, projRoleId:%v, lang:%v",
				sysRoleId, projRoleId, lang)
			// a wrong, uncomplete mapping is not a internal error
			// maybe even wanted for some user groups
			s.authFailed(false, loginEmail, userIp, fmt.Sprintf("cannot map saml attributes: %v", err), w, r)
			return
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
			s.authFailed(true, loginEmail, userIp, fmt.Sprintf("cannot add user: %v", err), w, r)
			return
		}

		err = UpdateElionaUserArchivedPhone(user.Email, phone, s.userToArchive)
		if err != nil {
			s.authFailed(true, loginEmail, userIp, fmt.Sprintf("failed to set phone number and archived flag: %v", err), w, r)
			return
		}

		err = SetUserPermissions(user.Id.Get(), sysRoleId, lang)
		if err != nil {
			s.authFailed(true, loginEmail, userIp, fmt.Sprintf("cannot set system user cnf: %v", err), w, r)
			return
		}
		err = SetProjectUser(projectId, user.Id.Get(), projRoleId)
		if err != nil {
			s.authFailed(true, loginEmail, userIp, fmt.Sprintf("cannot set project user cnf: %v", err), w, r)
			return
		}
	}

	// obtain a jwt to login via cookies
	jwt, err = GetElionaJsonWebToken(user.Email)
	if err != nil {
		s.authFailed(true, loginEmail, userIp, "cannot obtain a JWT", w, r)
		return
	}

	log.Debug(LOG_REGIO, "User %s with token %v", user.Email, jwt)
	if jwt == nil || *jwt == "" {
		s.authFailed(true, loginEmail, userIp, "invalid JWT returned", w, r)
		return
	} else {
		s.authSuccessful(loginEmail, jwt, w, r)
	}
}

func (s *SingleSignOn) authSuccessful(login string, jwt *string, w http.ResponseWriter, r *http.Request) {
	log.Info(LOG_REGIO, "authenticated user login: %s", login)
	setCookies := http.Cookie{
		Name:  "elionaAuthorization",
		Value: *jwt,
		Path:  "/"}

	http.SetCookie(w, &setCookies)
	http.Redirect(w, r, s.baseUrl, http.StatusFound)
}

func (s *SingleSignOn) authFailed(intError bool, login string, ip string, errorMsg string,
	w http.ResponseWriter, r *http.Request) {

	log.Info(LOG_REGIO, "not authenticated user tried to login: %s, %s",
		login, ip)

	if intError {
		log.Warn(LOG_REGIO, "internal server error occured: %s", errorMsg)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// reset eliona cookies
	setCookies := http.Cookie{
		Name:  "elionaAuthorization",
		Value: "invalid",
		Path:  "/"}
	http.SetCookie(w, &setCookies)

	if s.redirectNoLogin == "" {
		// fallback
		u, err := url.Parse(s.baseUrl + "/adfs/error.html")
		if err != nil {
			log.Error(LOG_REGIO, "cannot parse fallback redirect url: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMsg + ":" + err.Error()))
			return
		}
		queries := url.Values{}

		queries.Add("title", "Login Error")
		queries.Add("message", errorMsg)
		queries.Add("details", errorMsg)
		queries.Add("linkText", "Go To Login")
		queries.Add("link", s.baseUrl)

		u.RawQuery = queries.Encode()

		http.Redirect(w, r, u.String(), http.StatusFound)

	} else if !s.htmlContent {
		// redirect
		http.Redirect(w, r, s.redirectNoLogin, http.StatusFound)
	} else {
		// write html content
		w.Write([]byte(utils.SubstituteError(s.redirectNoLogin, []byte(errorMsg))))
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
