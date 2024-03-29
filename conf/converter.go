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

	return &appdb.Permission{
		ID:                      permissions.Id,
		DefaultSystemRole:       permissions.DefaultSystemRole,
		DefaultProjRole:         permissions.DefaultProjRole,
		SystemRoleSamlAttribute: null.StringFromPtr(permissions.SystemRoleSamlAttribute),
		SystemRoleMap:           sysRoleMap,
		ProjRoleSamlAttribute:   null.StringFromPtr(permissions.ProjRoleSamlAttribute),
		ProjRoleMap:             projRoleMap,
	}, nil
}

func PermissionDbToApiForm(permission *appdb.Permission) (*apiserver.Permissions, error) {

	var (
		systemRoleMap *[]apiserver.RoleMap
		projRoleMap   *[]apiserver.RoleMap
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

	return &apiserver.Permissions{
		Id:                      permission.ID,
		DefaultSystemRole:       permission.DefaultSystemRole,
		DefaultProjRole:         permission.DefaultProjRole,
		SystemRoleSamlAttribute: permission.SystemRoleSamlAttribute.Ptr(),
		SystemRoleMap:           systemRoleMap,
		ProjRoleSamlAttribute:   permission.ProjRoleSamlAttribute.Ptr(),
		ProjRoleMap:             projRoleMap,
	}, nil
}

func RoleMapToNullableJSON(roleMapPtr *[]apiserver.RoleMap) (null.JSON, error) {
	var (
		jsonBytes []byte
		err       error
	)

	if roleMapPtr == nil {
		return null.JSONFromPtr(nil), nil
	}
	jsonBytes, err = json.Marshal(*roleMapPtr)
	return null.JSONFrom(jsonBytes), err
}

func NullableJSONToRoleMapPtr(nullableJson null.JSON) (*[]apiserver.RoleMap, error) {
	var (
		roleMap []apiserver.RoleMap = []apiserver.RoleMap{}
		err     error
	)

	if nullableJson.Ptr() == nil {
		return nil, nil
	}

	err = json.Unmarshal(nullableJson.JSON, &roleMap)
	return &roleMap, err
}
