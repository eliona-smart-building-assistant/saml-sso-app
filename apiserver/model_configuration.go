/*
 * App SAML 2.0 SSO API
 *
 * API to access and configure the SAML 2.0 SSO service provider
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

// Configuration - The Configurations for running a SAML 2.0 Service Provider
type Configuration struct {

	// Configuration Id. Can only be 1
	Id int32 `json:"id,omitempty"`

	// If the configuration is enabled or not
	Enable bool `json:"enable,omitempty"`

	// The Certificate of this SAML Service Provider (SP). Can be a self-signed x509 certificate.
	ServiceProviderCertificate string `json:"serviceProviderCertificate,omitempty"`

	// The Private Key matching the Certificate of this SAML Service Provider (SP). Can be the Private Key of a self-signed x509 certificate. DO NOT use rsa key length lower than 2048
	ServiceProviderPrivateKey string `json:"serviceProviderPrivateKey,omitempty"`

	// The Metadata URL of the Identity Provider (IdP) if available. Otherwise use the metadataXml to provide Metadata of IdP directly and leave this null
	IdpMetadataUrl *string `json:"idpMetadataUrl,omitempty"`

	// Provide the IdP Metadata XML directly, if you have not the idpMetadataUrl accessable
	IdpMetadataXml *string `json:"idpMetadataXml,omitempty"`

	// The own URL of this Eliona instance
	OwnUrl string `json:"ownUrl,omitempty"`

	// If enabled, the new created user is archived and cannot login until a admin has activated it.
	UserToArchive bool `json:"userToArchive,omitempty"`

	// If the configuration is enabled or not
	AllowInitializationByIdp bool `json:"allowInitializationByIdp,omitempty"`

	// If the SP should make a signed SAML Authn-Request or not
	SignedRequest bool `json:"signedRequest,omitempty"`

	// Normaly this value is set to false for a SP. If set to true the user has to re-authenticate (login at IdP) even it has a valid session to the IdP.
	ForceAuthn bool `json:"forceAuthn,omitempty"`

	// If you have to use a customized Entity Id, you can overwrite it here. Normally the default value can be left as it is.
	EntityId string `json:"entityId,omitempty"`

	// only send cookies over encrypted connection (HTTPS)
	CookieSecure bool `json:"cookieSecure,omitempty"`

	// The url to redirect if the login failed. If this value is null the default page /noLogin will showed up
	LoginFailedUrl string `json:"loginFailedUrl,omitempty"`
}

// AssertConfigurationRequired checks if the required fields are not zero-ed
func AssertConfigurationRequired(obj Configuration) error {
	return nil
}

// AssertConfigurationConstraints checks if the values respects the defined constraints
func AssertConfigurationConstraints(obj Configuration) error {
	return nil
}
