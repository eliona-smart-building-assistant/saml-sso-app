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

package apiservices

import (
	"context"
	"errors"
	"net/http"
	"saml-sso/apiserver"
)

// ConfigurationApiService is a service that implements the logic for the ConfigurationAPIServicer
// This service should implement the business logic for every endpoint for the ConfigurationApi API.
// Include any external packages or services that will be required by this service.
type ConfigurationApiService struct {
}

// NewConfigurationApiService creates a default api service
func NewConfigurationApiService() apiserver.ConfigurationAPIServicer {
	return &ConfigurationApiService{}
}

// GetAttributeMapping - Get Attribute Mapping
func (s *ConfigurationApiService) GetAttributeMapping(ctx context.Context) (apiserver.ImplResponse, error) {
	// TODO - update GetAttributeMapping with the required logic for this service method.
	// Add api_configuration_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, AttributeMap{}) or use other options such as http.Ok ...
	//return Response(200, AttributeMap{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("GetAttributeMapping method not implemented")
}

// GetConfiguration - Get Basic Configurations
func (s *ConfigurationApiService) GetConfiguration(ctx context.Context) (apiserver.ImplResponse, error) {
	// TODO - update GetConfiguration with the required logic for this service method.
	// Add api_configuration_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Configuration{}) or use other options such as http.Ok ...
	//return Response(200, Configuration{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("GetConfiguration method not implemented")
}

// GetPermissionMapping - Get Permission Mapping
func (s *ConfigurationApiService) GetPermissionMapping(ctx context.Context) (apiserver.ImplResponse, error) {
	// TODO - update GetPermissionMapping with the required logic for this service method.
	// Add api_configuration_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Permissions{}) or use other options such as http.Ok ...
	//return Response(200, Permissions{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("GetPermissionMapping method not implemented")
}

// PutAttributeMapping - Creates or Update Attribute Mapping
func (s *ConfigurationApiService) PutAttributeMapping(ctx context.Context, attributeMap apiserver.AttributeMap) (apiserver.ImplResponse, error) {
	// TODO - update PutAttributeMapping with the required logic for this service method.
	// Add api_configuration_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, AttributeMap{}) or use other options such as http.Ok ...
	//return Response(200, AttributeMap{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("PutAttributeMapping method not implemented")
}

// PutConfiguration - Creates or Update Basic Configuration
func (s *ConfigurationApiService) PutConfiguration(ctx context.Context, basicConfiguration apiserver.Configuration) (apiserver.ImplResponse, error) {
	// TODO - update PutConfiguration with the required logic for this service method.
	// Add api_configuration_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Configuration{}) or use other options such as http.Ok ...
	//return Response(200, Configuration{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("PutConfiguration method not implemented")
}

// PutPermissionMapping - Creates or Update Permission Mapping Configurations
func (s *ConfigurationApiService) PutPermissionMapping(ctx context.Context, permissions apiserver.Permissions) (apiserver.ImplResponse, error) {
	// TODO - update PutPermissionMapping with the required logic for this service method.
	// Add api_configuration_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Permissions{}) or use other options such as http.Ok ...
	//return Response(200, Permissions{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("PutPermissionMapping method not implemented")
}
