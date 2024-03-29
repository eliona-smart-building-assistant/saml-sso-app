/*
 * App SAML 2.0 SSO API
 *
 * API to access and configure the SAML 2.0 SSO service provider
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

// Permissions - Sets default user permissions and optionaly maps SAML Attributes and Content to eliona's roles
type Permissions struct {

	// Configuration Id refer to config's id. Can only be 1
	Id int32 `json:"id,omitempty"`

	DefaultSystemRole string `json:"default_system_role,omitempty"`

	DefaultProjRole string `json:"default_proj_role,omitempty"`

	SystemRoleSamlAttribute *string `json:"system_role_saml_attribute,omitempty"`

	SystemRoleMap *[]RoleMap `json:"system_role_map,omitempty"`

	ProjRoleSamlAttribute *string `json:"proj_role_saml_attribute,omitempty"`

	ProjRoleMap *[]RoleMap `json:"proj_role_map,omitempty"`
}

// AssertPermissionsRequired checks if the required fields are not zero-ed
func AssertPermissionsRequired(obj Permissions) error {
	if obj.SystemRoleMap != nil {
		for _, el := range *obj.SystemRoleMap {
			if err := AssertRoleMapRequired(el); err != nil {
				return err
			}
		}
	}
	if obj.ProjRoleMap != nil {
		for _, el := range *obj.ProjRoleMap {
			if err := AssertRoleMapRequired(el); err != nil {
				return err
			}
		}
	}
	return nil
}

// AssertPermissionsConstraints checks if the values respects the defined constraints
func AssertPermissionsConstraints(obj Permissions) error {
	return nil
}
