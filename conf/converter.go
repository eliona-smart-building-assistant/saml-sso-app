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

package conf

import (
	"encoding/json"
	"errors"
	"saml-sso/apiserver"
	"saml-sso/appdb"

	"github.com/volatiletech/null/v8"
)

func ConvertApiToDbForm(apiForm any) (any, error) {

	var (
		err    error = nil
		dbForm any   = nil
	)

	switch a := apiForm.(type) {
	case *apiserver.BasicConfiguration:
		dbForm, err = BasicConfigApiToDbForm(a)
	case *apiserver.AdvancedConfiguration:
		dbForm, err = AdvancedConfigApiToDbForm(a)
	case *apiserver.AttributeMap:
		dbForm, err = AttributeMapApiToDbForm(a)
	case *apiserver.Permissions:
		dbForm, err = PermissionApiToDbForm(a)
	default:
		err = errors.New("unknown api datatype")
	}

	return dbForm, err
}

func ConvertDbToApiForm(dbForm any) (any, error) {

	var (
		err     error = nil
		apiForm any   = nil
	)

	switch d := dbForm.(type) {
	case *appdb.BasicConfig:
		apiForm, err = BasicConfigDbToApiForm(d)
	case *appdb.AdvancedConfig:
		apiForm, err = AdvancedConfigDbToApiForm(d)
	case *appdb.AttributeMap:
		apiForm, err = AttributeMapDbToApiForm(d)
	case *appdb.Permission:
		apiForm, err = PermissionDbToApiForm(d)
	default:
		err = errors.New("unknown db datatype")
	}

	return apiForm, err
}

func BasicConfigApiToDbForm(config *apiserver.BasicConfiguration) (*appdb.BasicConfig, error) {

	return &appdb.BasicConfig{
		Enable:         config.Enable,
		SPCertificate:  config.ServiceProviderCertificate,
		SPPrivateKey:   config.ServiceProviderPrivateKey,
		IdpMetadataURL: null.StringFromPtr(config.IdpMetadataUrl),
		MetadataXML:    null.StringFromPtr(config.IdpMetadataXml),
		OwnURL:         config.OwnUrl,
	}, nil
}
func BasicConfigDbToApiForm(config *appdb.BasicConfig) (*apiserver.BasicConfiguration, error) {

	return &apiserver.BasicConfiguration{
		Enable:                     config.Enable,
		ServiceProviderCertificate: config.SPCertificate,
		ServiceProviderPrivateKey:  config.SPPrivateKey,
		IdpMetadataUrl:             config.IdpMetadataURL.Ptr(),
		IdpMetadataXml:             config.MetadataXML.Ptr(),
		OwnUrl:                     config.OwnURL,
	}, nil
}

func AdvancedConfigApiToDbForm(config *apiserver.AdvancedConfiguration) (*appdb.AdvancedConfig, error) {
	return &appdb.AdvancedConfig{
		Enable:                   config.Enable,
		AllowInitializationByIdp: config.AllowInitializationByIdp,
		SignedRequest:            config.SignedRequest,
		ForceAuthn:               config.ForceAuthn,
		EntityID:                 config.EntityId,
		CookieSecure:             config.CookieSecure,
		LoginFailedURL:           config.LoginFailedUrl,
	}, nil
}
func AdvancedConfigDbToApiForm(config *appdb.AdvancedConfig) (*apiserver.AdvancedConfiguration, error) {
	return &apiserver.AdvancedConfiguration{
		Enable:                   config.Enable,
		AllowInitializationByIdp: config.AllowInitializationByIdp,
		SignedRequest:            config.SignedRequest,
		ForceAuthn:               config.ForceAuthn,
		EntityId:                 config.EntityID,
		CookieSecure:             config.CookieSecure,
		LoginFailedUrl:           config.LoginFailedURL,
	}, nil
}

func AttributeMapApiToDbForm(config *apiserver.AttributeMap) (*appdb.AttributeMap, error) {
	return &appdb.AttributeMap{
		Enable:    config.Enable,
		Email:     config.Email,
		FirstName: null.StringFromPtr(config.FirstName),
		LastName:  null.StringFromPtr(config.LastName),
		Phone:     null.StringFromPtr(config.Phone),
	}, nil
}
func AttributeMapDbToApiForm(config *appdb.AttributeMap) (*apiserver.AttributeMap, error) {
	return &apiserver.AttributeMap{
		Enable:    config.Enable,
		Email:     config.Email,
		FirstName: config.FirstName.Ptr(),
		LastName:  config.LastName.Ptr(),
		Phone:     config.Phone.Ptr(),
	}, nil
}

func PermissionApiToDbForm(permissions *apiserver.Permissions) (*appdb.Permission, error) {
	return &appdb.Permission{
		Enable:                  permissions.Enable,
		DefaultSystemRole:       permissions.DefaultSystemRole,
		DefaultProjRole:         permissions.DefaultProjRole,
		SystemRoleSamlAttribute: null.StringFromPtr(permissions.SystemRoleSamlAttribute),
		SystemRoleMap:           null.JSONFrom(marshal(permissions.SystemRoleMap)),
		ProjRoleSamlAttribute:   null.StringFromPtr(permissions.ProjRoleSamlAttribute),
		ProjRoleMap:             null.JSONFrom(marshal(permissions.ProjRoleMap)),
	}, nil
}

func PermissionDbToApiForm(permission *appdb.Permission) (*apiserver.Permissions, error) {

	var (
		systemRoleMap []apiserver.RoleMap
		projRoleMap   []apiserver.RoleMap
	)

	if err := unmarshal(permission.SystemRoleMap.JSON, &systemRoleMap); err != nil {
		systemRoleMap = nil
	}

	if err := unmarshal(permission.ProjRoleMap.JSON, &projRoleMap); err != nil {
		projRoleMap = nil
	}

	return &apiserver.Permissions{
		Enable:                  permission.Enable,
		DefaultSystemRole:       permission.DefaultSystemRole,
		DefaultProjRole:         permission.DefaultProjRole,
		SystemRoleSamlAttribute: getNullableString(permission.SystemRoleSamlAttribute),
		SystemRoleMap:           &systemRoleMap,
		ProjRoleSamlAttribute:   getNullableString(permission.ProjRoleSamlAttribute),
		ProjRoleMap:             &projRoleMap,
	}, nil
}

func marshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

func unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func getNullableString(s null.String) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}
