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

package conf

import (
	"encoding/json"
	"saml-sso/apiserver"
	"saml-sso/appdb"
	"strconv"

	"github.com/eliona-smart-building-assistant/go-utils/log"
	"github.com/volatiletech/null/v8"
)

func ConfigApiToDbForm(config *apiserver.Configuration) (*appdb.Config, error) {
	return &appdb.Config{
		ID:                       config.Id,
		Enable:                   config.Enable,
		SPCertificate:            config.ServiceProviderCertificate,
		SPPrivateKey:             config.ServiceProviderPrivateKey,
		IdpMetadataURL:           null.StringFromPtr(config.IdpMetadataUrl),
		MetadataXML:              null.StringFromPtr(config.IdpMetadataXml),
		OwnURL:                   config.OwnUrl,
		UserToArchive:            config.UserToArchive,
		AllowInitializationByIdp: config.AllowInitializationByIdp,
		SignedRequest:            config.SignedRequest,
		ForceAuthn:               config.ForceAuthn,
		EntityID:                 config.EntityId,
		CookieSecure:             config.CookieSecure,
		LoginFailedURL:           config.LoginFailedUrl,
	}, nil
}

func ConfigDbToApiForm(config *appdb.Config) (*apiserver.Configuration, error) {
	return &apiserver.Configuration{
		Id:                         config.ID,
		Enable:                     config.Enable,
		ServiceProviderCertificate: config.SPCertificate,
		ServiceProviderPrivateKey:  config.SPPrivateKey,
		IdpMetadataUrl:             config.IdpMetadataURL.Ptr(),
		IdpMetadataXml:             config.MetadataXML.Ptr(),
		OwnUrl:                     config.OwnURL,
		UserToArchive:              config.UserToArchive,
		AllowInitializationByIdp:   config.AllowInitializationByIdp,
		SignedRequest:              config.SignedRequest,
		ForceAuthn:                 config.ForceAuthn,
		EntityId:                   config.EntityID,
		CookieSecure:               config.CookieSecure,
		LoginFailedUrl:             config.LoginFailedURL,
	}, nil
}

func AttributeMapApiToDbForm(config *apiserver.AttributeMap) (*appdb.AttributeMap, error) {
	return &appdb.AttributeMap{
		ID:        config.Id,
		Email:     config.Email,
		FirstName: null.StringFromPtr(config.FirstName),
		LastName:  null.StringFromPtr(config.LastName),
		Phone:     null.StringFromPtr(config.Phone),
	}, nil
}
func AttributeMapDbToApiForm(config *appdb.AttributeMap) (*apiserver.AttributeMap, error) {
	return &apiserver.AttributeMap{
		Id:        config.ID,
		Email:     config.Email,
		FirstName: config.FirstName.Ptr(),
		LastName:  config.LastName.Ptr(),
		Phone:     config.Phone.Ptr(),
	}, nil
}

func PermissionApiToDbForm(permissions *apiserver.Permissions) (*appdb.Permission, error) {

	sysRoleMap, err := RoleMapToNullableJSON(permissions.SystemRoleMap)
	if err != nil {
		return nil, err
	}

	projRoleMap, err := RoleMapToNullableJSON(permissions.ProjRoleMap)
	if err != nil {
		return nil, err
	}

	langMap, err := RoleMapToNullableJSON(permissions.LanguageMap)
	if err != nil {
		return nil, err
	}

	return &appdb.Permission{
		ID:                      permissions.Id,
		DefaultSystemRole:       permissions.DefaultSystemRole,
		DefaultProjRole:         permissions.DefaultProjRole,
		SystemRoleSamlAttribute: null.StringFromPtr(permissions.SystemRoleSamlAttribute),
		SystemRoleMap:           sysRoleMap,
		ProjRoleSamlAttribute:   null.StringFromPtr(permissions.ProjRoleSamlAttribute),
		ProjRoleMap:             projRoleMap,
		DefaultLanguage:         permissions.DefaultLanguage,
		LanguageSamlAttribute:   null.StringFromPtr(permissions.LanguageSamlAttribute),
		LanguageMap:             langMap,
	}, nil
}

func PermissionDbToApiForm(permission *appdb.Permission) (*apiserver.Permissions, error) {

	var (
		systemRoleMap *[]apiserver.RoleMap
		projRoleMap   *[]apiserver.RoleMap
		langMap       *[]apiserver.RoleMap
		err           error
	)

	systemRoleMap, err = NullableJSONToRoleMapPtr(permission.ProjRoleMap)
	if err != nil {
		return nil, err
	}

	projRoleMap, err = NullableJSONToRoleMapPtr(permission.SystemRoleMap)
	if err != nil {
		return nil, err
	}

	langMap, err = NullableJSONToRoleMapPtr(permission.LanguageMap)
	if err != nil {
		return nil, err
	}

	return &apiserver.Permissions{
		Id:                      permission.ID,
		DefaultSystemRole:       permission.DefaultSystemRole,
		DefaultProjRole:         permission.DefaultProjRole,
		SystemRoleSamlAttribute: permission.SystemRoleSamlAttribute.Ptr(),
		SystemRoleMap:           systemRoleMap,
		ProjRoleSamlAttribute:   permission.ProjRoleSamlAttribute.Ptr(),
		ProjRoleMap:             projRoleMap,
		DefaultLanguage:         permission.DefaultLanguage,
		LanguageSamlAttribute:   permission.LanguageSamlAttribute.Ptr(),
		LanguageMap:             langMap,
	}, nil
}

// convert [{"ElionaRole":"roleIdOrRoleName", "SamlValue":"samlValue"}] to {"samlValue":"roleIdOrRoleName"}
// to use it as map[string]any
func RoleMapToNullableJSON(roleMapPtr *[]apiserver.RoleMap) (null.JSON, error) {
	var (
		jsonBytes []byte
		err       error

		roleMap map[string]any
	)

	if roleMapPtr == nil {
		return null.JSONFromPtr(nil), nil
	}

	roleMap = ApiRoleMapToGolangMap(*roleMapPtr)

	jsonBytes, err = json.Marshal(roleMap)
	return null.JSONFrom(jsonBytes), err
}

func NullableJSONToRoleMapPtr(nullableJson null.JSON) (*[]apiserver.RoleMap, error) {
	var (
		roleMapApi []apiserver.RoleMap
		roleMapDb  map[string]any
		err        error
	)

	if nullableJson.Ptr() == nil {
		return nil, nil
	}

	err = json.Unmarshal(nullableJson.JSON, &roleMapDb)

	roleMapApi = DBRoleToApiRole(roleMapDb)

	return &roleMapApi, err
}

func ApiRoleMapToGolangMap(roleMap []apiserver.RoleMap) (gRoleMap map[string]any) {

	gRoleMap = make(map[string]any)

	for _, m := range roleMap {
		var elionaRole any = m.ElionaRole

		if i, err := strconv.Atoi(m.ElionaRole); err == nil {
			elionaRole = i
		}

		gRoleMap[m.SamlValue] = elionaRole
	}

	return
}

func DBRoleToApiRole(roleMap map[string]any) (apiRoleMap []apiserver.RoleMap) {

	apiRoleMap = make([]apiserver.RoleMap, 0)

	for samlValue, elionaRole := range roleMap {
		var elionaRoleS string

		switch v := elionaRole.(type) {
		case string:
			elionaRoleS = v
		case int:
			elionaRoleS = strconv.Itoa(v)
		case float64:
			elionaRoleS = strconv.Itoa(int(v))
		default:
			log.Warn(LOG_REGIO, "unknown type for eliona role: %T, %v", v, v)
		}

		var singleApiMap apiserver.RoleMap = apiserver.RoleMap{
			ElionaRole: elionaRoleS,
			SamlValue:  samlValue,
		}

		apiRoleMap = append(apiRoleMap, singleApiMap)
	}

	return
}

func AnyToRoleId(roleNameOrId any, aclRoles map[string]int) (roleId int) {
	roleId = -1

	switch v := roleNameOrId.(type) {
	case string:
		roleId = StringToRoleId(v, aclRoles)
	case int:
		roleId = v
	case float64:
		roleId = int(v)
	default:
		log.Warn(LOG_REGIO, "unknown type for roleNameOrId: %T, %v", v, v)
	}

	return
}

func StringToRoleId(roleNameOrId string, aclRoles map[string]int) (roleId int) {
	var err error

	roleId, err = strconv.Atoi(roleNameOrId)
	if err != nil {
		roleId = aclRoles[roleNameOrId]
	}

	return
}
