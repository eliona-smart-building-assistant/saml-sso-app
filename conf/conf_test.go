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

package conf_test

import (
	"context"
	"fmt"
	"saml-sso/conf"
	"saml-sso/utils"
	"strings"
	"testing"
)

// Needs a DB with schema and exported CONNECTION_STRING
func TestConf_LoadAutoConfig(t *testing.T) {

	err := conf.DeleteAllConfigurations(context.Background())
	if err != nil {
		t.Error("Prepare DB: ", err)
	}

	err = conf.InsertAutoSamlConfiguration(context.Background())
	if err != nil {
		t.Error(err)
	}

	basicConfig, err := conf.GetBasicConfig(context.Background())
	if err != nil {
		t.Error(err)
	}
	if basicConfig.Enable != conf.AUTO_CNF_DEFAULT_ENABLED ||
		!strings.Contains(basicConfig.OwnUrl, "https://") ||
		basicConfig.IdpMetadataUrl != nil ||
		basicConfig.IdpMetadataXml != nil ||
		!strings.Contains(basicConfig.ServiceProviderCertificate,
			"-----BEGIN CERTIFICATE-----") ||
		!strings.Contains(basicConfig.ServiceProviderPrivateKey,
			"-----BEGIN RSA PRIVATE KEY-----") {
		t.Error("mismatch auto config")
	}
	_, err = utils.GetCombinedX509Certificate(basicConfig.ServiceProviderCertificate,
		basicConfig.ServiceProviderPrivateKey)
	if err != nil {
		t.Error("auto gen certificate")
	}
	fmt.Println(basicConfig)

	advancedConfig, err := conf.GetAdvancedConfig(context.Background())
	if err != nil {
		t.Error(err)
	}
	if advancedConfig.Enable != conf.AUTO_CNF_DEFAULT_ENABLED ||
		advancedConfig.AllowInitializationByIdp != conf.AUTO_CNF_DEFAULT_ALLOW_INIT_BY_IDP ||
		advancedConfig.CookieSecure != conf.AUTO_CNF_DEFAULT_COOKIE_SECURE ||
		advancedConfig.EntityId != conf.AUTO_CNF_DEFAULT_ENTITY_ID ||
		advancedConfig.ForceAuthn != conf.AUTO_CNF_DEFAULT_FORCE_AUTHN ||
		advancedConfig.LoginFailedUrl != conf.AUTO_CNF_DEFAULT_LOGIN_FAIL_URL ||
		advancedConfig.SignedRequest != conf.AUTO_CNF_DEFAULT_SIGNING_REQ {
		t.Error("invalid advanced auto conf")
	}

	attrMap, err := conf.GetAttributeMapping(context.Background())
	if err != nil {
		t.Error(err)
	}
	if attrMap.Email != conf.AUTO_CNF_DEFAULT_USERNAME_ATTR ||
		attrMap.Enable != conf.AUTO_CNF_DEFAULT_ENABLED ||
		attrMap.FirstName != nil ||
		attrMap.LastName != nil ||
		attrMap.Phone != nil {
		t.Error("invalid attribute map auto conf")
	}

	perms, err := conf.GetPermissionSettings(context.Background())
	if err != nil {
		t.Error(err)
	}
	if perms.DefaultProjRole != conf.AUTO_CNF_DEFAULT_PROJ_PERMISSION ||
		perms.DefaultSystemRole != conf.AUTO_CNF_DEFAULT_SYSTEM_PERMISSION ||
		perms.Enable != conf.AUTO_CNF_DEFAULT_ENABLED ||
		perms.ProjRoleMap != nil ||
		perms.ProjRoleSamlAttribute != nil ||
		perms.SystemRoleMap != nil ||
		perms.SystemRoleSamlAttribute != nil {
		t.Error("invalid permissions auto conf")
	}
}

func TestConf_InsertUpdateConfig(t *testing.T) {

}
