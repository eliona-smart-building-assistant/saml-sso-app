/*
 * App SAML 2.0 SSO API
 *
 * API to access and configure the SAML 2.0 SSO service provider
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

// Active - If the service is active
type Active struct {
	Active bool `json:"active,omitempty"`
}

// AssertActiveRequired checks if the required fields are not zero-ed
func AssertActiveRequired(obj Active) error {
	return nil
}

// AssertRecurseActiveRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Active (e.g. [][]Active), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseActiveRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aActive, ok := obj.(Active)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertActiveRequired(aActive)
	})
}
