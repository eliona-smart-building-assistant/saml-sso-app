/*
 * App SAML 2.0 SSO API
 *
 * API to access and configure the SAML 2.0 SSO service provider
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

// BasicConfiguration - The Basic Configurations for running a SAML 2.0 Service Provider
type BasicConfiguration struct {

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
}

// AssertBasicConfigurationRequired checks if the required fields are not zero-ed
func AssertBasicConfigurationRequired(obj BasicConfiguration) error {
	return nil
}

// AssertRecurseBasicConfigurationRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of BasicConfiguration (e.g. [][]BasicConfiguration), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseBasicConfigurationRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aBasicConfiguration, ok := obj.(BasicConfiguration)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertBasicConfigurationRequired(aBasicConfiguration)
	})
}
