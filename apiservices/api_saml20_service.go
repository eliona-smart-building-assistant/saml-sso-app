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

// SAML20ApiService is a service that implements the logic for the SAML20ApiServicer
// This service should implement the business logic for every endpoint for the SAML20Api API.
// Include any external packages or services that will be required by this service.
type SAML20ApiService struct {
}

// NewSAML20ApiService creates a default api service
func NewSAML20ApiService() apiserver.SAML20ApiServicer {
	return &SAML20ApiService{}
}

// SamlAcsPost -
func (s *SAML20ApiService) SamlAcsPost(ctx context.Context) (apiserver.ImplResponse, error) {
	// TODO - update SamlAcsPost with the required logic for this service method.
	// Add api_saml20_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(302, {}) or use other options such as http.Ok ...
	//return Response(302, nil),nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("SamlAcsPost method not implemented")
}

// SamlSloPost -
func (s *SAML20ApiService) SamlSloPost(ctx context.Context) (apiserver.ImplResponse, error) {
	// TODO - update SamlSloPost with the required logic for this service method.
	// Add api_saml20_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(302, {}) or use other options such as http.Ok ...
	//return Response(302, nil),nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("SamlSloPost method not implemented")
}
