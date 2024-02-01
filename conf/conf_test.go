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

package conf_test

import (
	"context"
	"saml-sso/apiserver"
	"saml-sso/conf"
	"saml-sso/utils"
	"strings"
	"testing"

	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/go-test/deep"
)

// Needs a DB with exported CONNECTION_STRING
func TestApp_Conf_InitDB(t *testing.T) {
	err := conf.UserLeicomInit()
	if err != nil {
		t.Log("user leicom, ", err)
	}
	err = conf.DropOwnSchema()
	if err != nil {
		// no error, if schema not exist
		t.Log("drop schema, ", err)
	}

	execFunc := app.ExecSqlFile("init.sql")
	err = execFunc(db.NewConnection())
	if err != nil {
		t.Error("init.sql failed, ", err)
	}
}

func TestApp_Conf_LoadAutoConfig(t *testing.T) {
	err := conf.DeleteAllConfigurations(context.Background())
	if err != nil {
		t.Error("Prepare DB (delete all config): ", err)
	}

	err = conf.InsertAutoSamlConfiguration(context.Background())
	if err != nil {
		t.Error(err)
	}

	config, err := conf.GetConfig(context.Background())
	if err != nil {
		t.Error(err)
	}
	if config.Enable != conf.AUTO_CNF_DEFAULT_ENABLED ||
		config.UserToArchive != conf.AUTO_CNF_DEFAULT_USER_TO_ARCHIVE ||
		!strings.Contains(config.OwnUrl, "https://") ||
		config.IdpMetadataUrl != nil ||
		config.IdpMetadataXml != nil ||
		!strings.Contains(config.ServiceProviderCertificate,
			"-----BEGIN CERTIFICATE-----") ||
		!strings.Contains(config.ServiceProviderPrivateKey,
			"-----BEGIN RSA PRIVATE KEY-----") {
		t.Error("mismatch auto config")
	}
	_, err = utils.GetCombinedX509Certificate(config.ServiceProviderCertificate,
		config.ServiceProviderPrivateKey)
	if err != nil {
		t.Error("auto gen certificate")
	}

	if config.AllowInitializationByIdp != conf.AUTO_CNF_DEFAULT_ALLOW_INIT_BY_IDP ||
		config.CookieSecure != conf.AUTO_CNF_DEFAULT_COOKIE_SECURE ||
		config.EntityId != conf.AUTO_CNF_DEFAULT_ENTITY_ID ||
		config.ForceAuthn != conf.AUTO_CNF_DEFAULT_FORCE_AUTHN ||
		config.LoginFailedUrl != conf.AUTO_CNF_DEFAULT_LOGIN_FAIL_URL ||
		config.SignedRequest != conf.AUTO_CNF_DEFAULT_SIGNING_REQ {
		t.Error("invalid advanced auto conf")
	}

	attrMap, err := conf.GetAttributeMapping(context.Background())
	if err != nil {
		t.Error(err)
	}
	if attrMap.Email != conf.AUTO_CNF_DEFAULT_USERNAME_ATTR ||
		attrMap.Id != 1 ||
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
		perms.Id != 1 ||
		perms.ProjRoleMap != nil ||
		perms.ProjRoleSamlAttribute != nil ||
		perms.SystemRoleMap != nil ||
		perms.SystemRoleSamlAttribute != nil {
		t.Error("invalid permissions auto conf")
	}
}

func TestApp_Conf_InsertUpdateConfig(t *testing.T) {
	for i := 0; i < 5; i++ {
		var (
			config1      apiserver.Configuration = utils.CreateRandomApiConfig()
			attrMapping1 apiserver.AttributeMap  = utils.CreateRandomApiAttributeMap()
			perms1       apiserver.Permissions   = utils.CreateRandomApiPermissions()

			config2      apiserver.Configuration = utils.CreateRandomApiConfig()
			attrMapping2 apiserver.AttributeMap  = utils.CreateRandomApiAttributeMap()
			perms2       apiserver.Permissions   = utils.CreateRandomApiPermissions()
		)
		config1.Enable = false
		config2.Enable = true

		err := conf.DeleteAllConfigurations(context.Background())
		if err != nil {
			t.Error("Prepare DB (delete all config): ", err)
		}

		// Basic Config Test
		confRet1, err := conf.SetConfig(context.Background(), &config1)
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&config1, confRet1); diff != nil {
			t.Error("missmatch config 1: ", diff)
		}
		confRet1, err = conf.GetConfig(context.Background())
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&config1, confRet1); diff != nil {
			t.Error("missmatch config 1_1: ", diff)
		}
		confRet2, err := conf.SetConfig(context.Background(), &config2)
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&config2, confRet2); diff != nil {
			t.Error("missmatch config 2: ", diff)
		}
		confRet2, err = conf.GetConfig(context.Background())
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&config2, confRet2); diff != nil {
			t.Error("missmatch config 2_2: ", diff)
		}

		// Attribute Mapping Test
		attrMapRet1, err := conf.SetAttributeMapping(context.Background(), &attrMapping1)
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&attrMapping1, attrMapRet1); diff != nil {
			t.Error("missmatch attribute mapping 1: ", diff)
		}
		attrMapRet1, err = conf.GetAttributeMapping(context.Background())
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&attrMapping1, attrMapRet1); diff != nil {
			t.Error("missmatch attribute mapping 1_1: ", diff)
		}
		attrMapRet2, err := conf.SetAttributeMapping(context.Background(), &attrMapping2)
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&attrMapping2, attrMapRet2); diff != nil {
			t.Error("missmatch attribute mapping 2: ", diff)
		}
		attrMapRet2, err = conf.GetAttributeMapping(context.Background())
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&attrMapping2, attrMapRet2); diff != nil {
			t.Error("missmatch attribute mapping 2_2: ", diff)
		}

		// Permission Cnf Test
		permsRet1, err := conf.SetPermissionSettings(context.Background(), &perms1)
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&perms1, permsRet1); diff != nil {
			t.Error("missmatch permission config 1: ", diff)
		}
		permsRet1, err = conf.GetPermissionSettings(context.Background())
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&perms1, permsRet1); diff != nil {
			t.Error("missmatch permission config 1_1: ", diff)
		}
		permsRet2, err := conf.SetPermissionSettings(context.Background(), &perms2)
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&perms2, permsRet2); diff != nil {
			t.Error("missmatch permission config 2: ", diff)
		}
		permsRet2, err = conf.GetPermissionSettings(context.Background())
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(&perms2, permsRet2); diff != nil {
			t.Error("missmatch permission config 2_2: ", diff)
		}
	}
}
