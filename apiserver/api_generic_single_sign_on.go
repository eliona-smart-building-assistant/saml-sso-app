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

// GenericSingleSignOnAPIController binds http requests to an api service and writes the service results to the http response
type GenericSingleSignOnAPIController struct {
	service      GenericSingleSignOnAPIServicer
	errorHandler ErrorHandler
}

// GenericSingleSignOnAPIOption for how the controller is set up.
type GenericSingleSignOnAPIOption func(*GenericSingleSignOnAPIController)

// WithGenericSingleSignOnAPIErrorHandler inject ErrorHandler into controller
func WithGenericSingleSignOnAPIErrorHandler(h ErrorHandler) GenericSingleSignOnAPIOption {
	return func(c *GenericSingleSignOnAPIController) {
		c.errorHandler = h
	}
}

// NewGenericSingleSignOnAPIController creates a default api controller
func NewGenericSingleSignOnAPIController(s GenericSingleSignOnAPIServicer, opts ...GenericSingleSignOnAPIOption) Router {
	controller := &GenericSingleSignOnAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the GenericSingleSignOnAPIController
func (c *GenericSingleSignOnAPIController) Routes() Routes {
	return Routes{
		"GetAuthorizationProcedure": Route{
			strings.ToUpper("Get"),
			"/v1/sso/auth",
			c.GetAuthorizationProcedure,
		},
		"GetSSOActive": Route{
			strings.ToUpper("Get"),
			"/v1/sso/active",
			c.GetSSOActive,
		},
	}
}

// GetAuthorizationProcedure - Begin authorization / login procedure
func (c *GenericSingleSignOnAPIController) GetAuthorizationProcedure(w http.ResponseWriter, r *http.Request) {
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
func (c *GenericSingleSignOnAPIController) GetSSOActive(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetSSOActive(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
