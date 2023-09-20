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
		err = errors.New("not implemented")
	case *apiserver.AttributeMap:
		err = errors.New("not implemented")
	case *apiserver.Permissions:
		err = errors.New("not implemented")
	default:
		err = errors.New("unknown datatype")
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
		err = errors.New("not implemented")
	case *appdb.AttributeMap:
		err = errors.New("not implemented")
	case *appdb.Permission:
		err = errors.New("not implemented")
	default:
		err = errors.New("unknown datatype")
	}

	return apiForm, err
}

func BasicConfigApiToDbForm(config *apiserver.BasicConfiguration) (*appdb.BasicConfig, error) {

	return &appdb.BasicConfig{
		Enable:         config.Enable,
		SPCertificate:  config.ServiceProviderCertificate,
		SPPrivateKey:   config.ServiceProviderPrivateKey,
		IdpMetadataURL: null.NewString(*config.IdpMetadataUrl, config.IdpMetadataUrl != nil),
		MetadataXML:    null.NewString(*config.IdpMetadataXml, config.IdpMetadataXml != nil),
		OwnURL:         null.NewString(config.OwnUrl, config.OwnUrl != ""),
	}, nil
}
func BasicConfigDbToApiForm(config *appdb.BasicConfig) (*apiserver.BasicConfiguration, error) {

	return &apiserver.BasicConfiguration{
		Enable:                     config.Enable,
		ServiceProviderCertificate: config.SPCertificate,
		ServiceProviderPrivateKey:  config.SPPrivateKey,
		IdpMetadataUrl:             config.IdpMetadataURL.Ptr(),
		IdpMetadataXml:             config.MetadataXML.Ptr(),
		OwnUrl:                     config.OwnURL.String,
	}, nil
}

func AdvancedConfigApiToDbForm(config *apiserver.AdvancedConfiguration) (*appdb.AdvancedConfig, error) {
	return &appdb.AdvancedConfig{
		Enable:                   config.Enable,
		AllowInitializationByIdp: config.AllowInitializationByIdp,
		SignedRequest:            config.SignedRequest,
		ForceAuthn:               config.ForceAuthn,
		EntityID:                 config.EntityId,
		CookieSecure:             config.CookieSecure != nil,
		LoginFailedURL:           "",
	}, nil
}
func AdvancedConfigDbToApiForm(config *appdb.AdvancedConfig) (*apiserver.AdvancedConfiguration, error) {
	return &apiserver.AdvancedConfiguration{
		Enable:                   config.Enable,
		AllowInitializationByIdp: config.AllowInitializationByIdp,
		SignedRequest:            config.SignedRequest,
		ForceAuthn:               config.ForceAuthn,
		EntityId:                 config.EntityID,
		CookieSecure:             nil,
		LoginFailedUrl:           nil,
	}, nil
}

func AttributeMapApiToDbForm(config *apiserver.AttributeMap) (*appdb.AttributeMap, error) {
	return &appdb.AttributeMap{
		Enable:    config.Enable,
		Email:     config.Email,
		FirstName: null.NewString("", config.FirstName != nil),
		LastName:  null.NewString("", config.LastName != nil),
		Phone:     null.NewString("", config.Phone != nil),
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
		SystemRoleSamlAttribute: null.NewString(*permissions.SystemRoleSamlAttribute, permissions.SystemRoleSamlAttribute != nil),
		SystemRoleMap:           null.JSONFrom(marshal(permissions.SystemRoleMap)),
		ProjRoleSamlAttribute:   null.NewString(*permissions.ProjRoleSamlAttribute, permissions.ProjRoleSamlAttribute != nil),
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
