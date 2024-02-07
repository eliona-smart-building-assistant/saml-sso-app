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
	"net/http"
	"saml-sso/apiserver"
	"saml-sso/conf"
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
	am, err := conf.GetAttributeMapping(ctx)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusOK, am), nil
}

// GetConfiguration - Get Configuration
func (s *ConfigurationApiService) GetConfiguration(ctx context.Context) (apiserver.ImplResponse, error) {
	config, err := conf.GetConfig(ctx)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusOK, config), nil
}

// GetPermissionMapping - Get Permission Mapping
func (s *ConfigurationApiService) GetPermissionMapping(ctx context.Context) (apiserver.ImplResponse, error) {
	pm, err := conf.GetPermissionMapping(ctx)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusOK, pm), nil
}

// PutAttributeMapping - Creates or Update Attribute Mapping
func (s *ConfigurationApiService) PutAttributeMapping(ctx context.Context, attributeMap apiserver.AttributeMap) (apiserver.ImplResponse, error) {
	upsertedAttrMap, err := conf.SetAttributeMapping(ctx, &attributeMap)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusCreated, upsertedAttrMap), nil
}

// PutConfiguration - Creates or Update Basic Configuration
func (s *ConfigurationApiService) PutConfiguration(ctx context.Context, config apiserver.Configuration) (apiserver.ImplResponse, error) {
	upsertedConfig, err := conf.SetConfig(ctx, &config)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusCreated, upsertedConfig), nil
}

// PutPermissionMapping - Creates or Update Permission Mapping Configurations
func (s *ConfigurationApiService) PutPermissionMapping(ctx context.Context, permissions apiserver.Permissions) (apiserver.ImplResponse, error) {
	upsertedPerms, err := conf.SetPermissionMapping(ctx, &permissions)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusCreated, upsertedPerms), nil
}
