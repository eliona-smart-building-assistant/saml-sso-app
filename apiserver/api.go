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
	"context"
	"net/http"
)

// ConfigurationAPIRouter defines the required methods for binding the api requests to a responses for the ConfigurationAPI
// The ConfigurationAPIRouter implementation should parse necessary information from the http request,
// pass the data to a ConfigurationAPIServicer to perform the required actions, then write the service results to the http response.
type ConfigurationAPIRouter interface {
	GetAttributeMapping(http.ResponseWriter, *http.Request)
	GetConfiguration(http.ResponseWriter, *http.Request)
	GetPermissionMapping(http.ResponseWriter, *http.Request)
	PutAttributeMapping(http.ResponseWriter, *http.Request)
	PutConfiguration(http.ResponseWriter, *http.Request)
	PutPermissionMapping(http.ResponseWriter, *http.Request)
}

// GenericSingleSignOnAPIRouter defines the required methods for binding the api requests to a responses for the GenericSingleSignOnAPI
// The GenericSingleSignOnAPIRouter implementation should parse necessary information from the http request,
// pass the data to a GenericSingleSignOnAPIServicer to perform the required actions, then write the service results to the http response.
type GenericSingleSignOnAPIRouter interface {
	GetAuthorizationProcedure(http.ResponseWriter, *http.Request)
	GetSSOActive(http.ResponseWriter, *http.Request)
}

// SAML20APIRouter defines the required methods for binding the api requests to a responses for the SAML20API
// The SAML20APIRouter implementation should parse necessary information from the http request,
// pass the data to a SAML20APIServicer to perform the required actions, then write the service results to the http response.
type SAML20APIRouter interface {
	SamlAcsPost(http.ResponseWriter, *http.Request)
	SamlSloPost(http.ResponseWriter, *http.Request)
}

// VersionAPIRouter defines the required methods for binding the api requests to a responses for the VersionAPI
// The VersionAPIRouter implementation should parse necessary information from the http request,
// pass the data to a VersionAPIServicer to perform the required actions, then write the service results to the http response.
type VersionAPIRouter interface {
	GetOpenAPI(http.ResponseWriter, *http.Request)
	GetVersion(http.ResponseWriter, *http.Request)
}

// ConfigurationAPIServicer defines the api actions for the ConfigurationAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ConfigurationAPIServicer interface {
	GetAttributeMapping(context.Context) (ImplResponse, error)
	GetConfiguration(context.Context) (ImplResponse, error)
	GetPermissionMapping(context.Context) (ImplResponse, error)
	PutAttributeMapping(context.Context, AttributeMap) (ImplResponse, error)
	PutConfiguration(context.Context, Configuration) (ImplResponse, error)
	PutPermissionMapping(context.Context, Permissions) (ImplResponse, error)
}

// GenericSingleSignOnAPIServicer defines the api actions for the GenericSingleSignOnAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type GenericSingleSignOnAPIServicer interface {
	GetAuthorizationProcedure(context.Context) (ImplResponse, error)
	GetSSOActive(context.Context) (ImplResponse, error)
}

// SAML20APIServicer defines the api actions for the SAML20API service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type SAML20APIServicer interface {
	SamlAcsPost(context.Context) (ImplResponse, error)
	SamlSloPost(context.Context) (ImplResponse, error)
}

// VersionAPIServicer defines the api actions for the VersionAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type VersionAPIServicer interface {
	GetOpenAPI(context.Context) (ImplResponse, error)
	GetVersion(context.Context) (ImplResponse, error)
}
