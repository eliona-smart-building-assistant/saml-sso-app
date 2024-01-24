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

// VersionApiService is a service that implements the logic for the VersionApiServicer
// This service should implement the business logic for every endpoint for the VersionApi API.
// Include any external packages or services that will be required by this service.
type VersionApiService struct {
}

// NewVersionApiService creates a default api service
func NewVersionApiService() apiserver.VersionAPIServicer {
	return &VersionApiService{}
}

// GetOpenAPI - OpenAPI specification for this API version
func (s *VersionApiService) GetOpenAPI(ctx context.Context) (apiserver.ImplResponse, error) {
	// TODO - update GetOpenAPI with the required logic for this service method.
	// Add api_version_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, map[string]interface{}{}) or use other options such as http.Ok ...
	//return Response(200, map[string]interface{}{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("GetOpenAPI method not implemented")
}

// GetVersion - Version of the API
func (s *VersionApiService) GetVersion(ctx context.Context) (apiserver.ImplResponse, error) {
	// TODO - update GetVersion with the required logic for this service method.
	// Add api_version_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, map[string]interface{}{}) or use other options such as http.Ok ...
	//return Response(200, map[string]interface{}{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("GetVersion method not implemented")
}
