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

// GetAdvancedConfiguration - Get Advanced Configuration
func (s *ConfigurationApiService) GetAdvancedConfiguration(ctx context.Context) (apiserver.ImplResponse, error) {
	advCnf, err := conf.GetAdvancedConfig(ctx)

	return apiserver.Response(http.StatusOK, advCnf), err
}

// GetAttributeMapping - Get Attribute Mapping
func (s *ConfigurationApiService) GetAttributeMapping(ctx context.Context) (apiserver.ImplResponse, error) {
	attrMap, err := conf.GetAttributeMapping(ctx)

	return apiserver.Response(http.StatusOK, attrMap), err
}

// GetBasicConfiguration - Get Basic Configurations
func (s *ConfigurationApiService) GetBasicConfiguration(ctx context.Context) (apiserver.ImplResponse, error) {
	basicCnf, err := conf.GetBasicConfig(ctx)

	return apiserver.Response(http.StatusOK, basicCnf), err
}

// GetPermissionMapping - Get Permission Mapping
func (s *ConfigurationApiService) GetPermissionMapping(ctx context.Context) (apiserver.ImplResponse, error) {
	permMap, err := conf.GetPermissionSettings(ctx)

	return apiserver.Response(http.StatusOK, permMap), err
}

// PutAdvancedConfiguration - Creates or Update Advanced Configuration
func (s *ConfigurationApiService) PutAdvancedConfiguration(ctx context.Context, advancedConfiguration apiserver.AdvancedConfiguration) (apiserver.ImplResponse, error) {
	cnfRet, err := conf.SetAdvancedConfig(ctx, &advancedConfiguration)

	return apiserver.Response(http.StatusOK, cnfRet), err
}

// PutAttributeMapping - Creates or Update Attribute Mapping
func (s *ConfigurationApiService) PutAttributeMapping(ctx context.Context, attributeMap apiserver.AttributeMap) (apiserver.ImplResponse, error) {
	mapRet, err := conf.SetAttributeMapping(ctx, &attributeMap)

	return apiserver.Response(http.StatusOK, mapRet), err
}

// PutBasicConfiguration - Creates or Update Basic Configuration
func (s *ConfigurationApiService) PutBasicConfiguration(ctx context.Context, basicConfiguration apiserver.BasicConfiguration) (apiserver.ImplResponse, error) {
	cnfRet, err := conf.SetBasicConfig(ctx, &basicConfiguration)

	return apiserver.Response(http.StatusOK, cnfRet), err
}

// PutPermissionMapping - Creates or Update Permission Mapping Configurations
func (s *ConfigurationApiService) PutPermissionMapping(ctx context.Context, permissions apiserver.Permissions) (apiserver.ImplResponse, error) {
	permRet, err := conf.SetPermissionSettings(ctx, &permissions)

	return apiserver.Response(http.StatusOK, permRet), err
}
