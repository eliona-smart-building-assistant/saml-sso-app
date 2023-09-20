//  This file is part of the eliona project.
//  Copyright © 2023 LEICOM iTEC AG. All Rights Reserved.
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

// GenericSingleSignOnApiService is a service that implements the logic for the GenericSingleSignOnApiServicer
// This service should implement the business logic for every endpoint for the GenericSingleSignOnApi API.
// Include any external packages or services that will be required by this service.
type GenericSingleSignOnApiService struct {
}

// NewGenericSingleSignOnApiService creates a default api service
func NewGenericSingleSignOnApiService() apiserver.GenericSingleSignOnApiServicer {
	return &GenericSingleSignOnApiService{}
}

// GetAuthorizationProcedure - Begin authorization / login procedure
func (s *GenericSingleSignOnApiService) GetAuthorizationProcedure(ctx context.Context) (apiserver.ImplResponse, error) {
	// TODO - update GetAuthorizationProcedure with the required logic for this service method.
	// Add api_generic_single_sign_on_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(302, {}) or use other options such as http.Ok ...
	//return Response(302, nil),nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("GetAuthorizationProcedure method not implemented")
}

// GetSSOActive - Check, if a SSO service is available and configured
func (s *GenericSingleSignOnApiService) GetSSOActive(ctx context.Context) (apiserver.ImplResponse, error) {
	// TODO - update GetSSOActive with the required logic for this service method.
	// Add api_generic_single_sign_on_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Active{}) or use other options such as http.Ok ...
	//return Response(200, Active{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("GetSSOActive method not implemented")
}
