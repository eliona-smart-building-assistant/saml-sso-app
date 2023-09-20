/*
 * App SAML 2.0 SSO API
 *
 * API to access and configure the SAML 2.0 SSO service provider
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

import (
	"net/http"
	"strings"
)

// GenericSingleSignOnApiController binds http requests to an api service and writes the service results to the http response
type GenericSingleSignOnApiController struct {
	service      GenericSingleSignOnApiServicer
	errorHandler ErrorHandler
}

// GenericSingleSignOnApiOption for how the controller is set up.
type GenericSingleSignOnApiOption func(*GenericSingleSignOnApiController)

// WithGenericSingleSignOnApiErrorHandler inject ErrorHandler into controller
func WithGenericSingleSignOnApiErrorHandler(h ErrorHandler) GenericSingleSignOnApiOption {
	return func(c *GenericSingleSignOnApiController) {
		c.errorHandler = h
	}
}

// NewGenericSingleSignOnApiController creates a default api controller
func NewGenericSingleSignOnApiController(s GenericSingleSignOnApiServicer, opts ...GenericSingleSignOnApiOption) Router {
	controller := &GenericSingleSignOnApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the GenericSingleSignOnApiController
func (c *GenericSingleSignOnApiController) Routes() Routes {
	return Routes{
		{
			"GetAuthorizationProcedure",
			strings.ToUpper("Get"),
			"/v1/sso/auth",
			c.GetAuthorizationProcedure,
		},
		{
			"GetSSOActive",
			strings.ToUpper("Get"),
			"/v1/sso/active",
			c.GetSSOActive,
		},
	}
}

// GetAuthorizationProcedure - Begin authorization / login procedure
func (c *GenericSingleSignOnApiController) GetAuthorizationProcedure(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetAuthorizationProcedure(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetSSOActive - Check, if a SSO service is available and configured
func (c *GenericSingleSignOnApiController) GetSSOActive(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetSSOActive(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
