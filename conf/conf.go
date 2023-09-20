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

package conf

import (
	"context"
	"database/sql"
	"errors"
	"os"

	"saml-sso/apiserver"
	"saml-sso/appdb"
	"saml-sso/utils"

	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	LOG_REGIO = "config"
)

const (
	AUTO_CNF_DEFAULT_ENABLED           = true
	AUTO_CNF_DEFAULT_ALLOW_INIT_BY_IDP = false
	AUTO_CNF_DEFAULT_SIGNING_REQ       = true
	AUTO_CNF_DEFAULT_FORCE_AUTHN       = false
	AUTO_CNF_DEFAULT_COOKIE_SECURE     = false
	AUTO_CNF_DEFAULT_ENTITY_ID         = "{ownUrl}/saml/metadata"
	AUTO_CNF_DEFAULT_LOGIN_FAIL_URL    = "{ownUrl}/noLogin"
	AUTO_CNF_DEFAULT_USERNAME_ATTR     = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"
	AUTO_CNF_DEFAULT_SYSTEM_PERMISSION = "regular"
	AUTO_CNF_DEFAULT_PROJ_PERMISSION   = "operator"
	AUTO_CNF_DEFAULT_CERT_VALID_DAYS   = 3650
	AUTO_CNF_DEFAULT_KEY_LENGTH        = 4096
)

func InsertAutoSamlConfiguration(ctx context.Context) error {

	var (
		basicConfig appdb.BasicConfig = appdb.BasicConfig{
			Enable:         AUTO_CNF_DEFAULT_ENABLED,
			IdpMetadataURL: null.StringFromPtr(nil),
			MetadataXML:    null.StringFromPtr(nil),
		}
		advancedConfig appdb.AdvancedConfig = appdb.AdvancedConfig{
			Enable:                   AUTO_CNF_DEFAULT_ENABLED,
			AllowInitializationByIdp: AUTO_CNF_DEFAULT_ALLOW_INIT_BY_IDP,
			SignedRequest:            AUTO_CNF_DEFAULT_SIGNING_REQ,
			ForceAuthn:               AUTO_CNF_DEFAULT_FORCE_AUTHN,
			EntityID:                 AUTO_CNF_DEFAULT_ENTITY_ID,
			CookieSecure:             AUTO_CNF_DEFAULT_COOKIE_SECURE,
			LoginFailedURL:           AUTO_CNF_DEFAULT_LOGIN_FAIL_URL,
		}
		attributeMapping appdb.AttributeMap = appdb.AttributeMap{
			Enable:    AUTO_CNF_DEFAULT_ENABLED,
			Email:     AUTO_CNF_DEFAULT_USERNAME_ATTR,
			FirstName: null.StringFromPtr(nil),
			LastName:  null.StringFromPtr(nil),
			Phone:     null.StringFromPtr(nil),
		}
		permissionCnf appdb.Permission = appdb.Permission{
			Enable:                  AUTO_CNF_DEFAULT_ENABLED,
			DefaultSystemRole:       AUTO_CNF_DEFAULT_SYSTEM_PERMISSION,
			DefaultProjRole:         AUTO_CNF_DEFAULT_PROJ_PERMISSION,
			SystemRoleSamlAttribute: null.StringFromPtr(nil),
			SystemRoleMap:           null.JSONFromPtr(nil),
			ProjRoleSamlAttribute:   null.StringFromPtr(nil),
			ProjRoleMap:             null.JSONFromPtr(nil),
		}
	)

	basicConfig.OwnURL = "https://" + GetElionaHost()
	certificate, privateKey, err := utils.CreateSelfsignedX509Certificate(AUTO_CNF_DEFAULT_CERT_VALID_DAYS, AUTO_CNF_DEFAULT_KEY_LENGTH, nil, nil)
	if err != nil {
		log.Error(LOG_REGIO, "auto configuration generate x509 certificates: %v", err)
	} else {
		basicConfig.SPCertificate = certificate
		basicConfig.SPPrivateKey = privateKey
	}

	err = basicConfig.Insert(ctx, getDb(), boil.Infer())
	if err != nil {
		return err
	}

	err = advancedConfig.Insert(ctx, getDb(), boil.Infer())
	if err != nil {
		return err
	}

	err = attributeMapping.Insert(ctx, getDb(), boil.Infer())
	if err != nil {
		return err
	}

	err = permissionCnf.Insert(ctx, getDb(), boil.Infer())
	if err != nil {
		return err
	}

	return err
}

func GetBasicConfig(ctx context.Context) (*apiserver.BasicConfiguration, error) {

	var (
		err            error                         = nil
		basicConfigDb  *appdb.BasicConfig            = nil
		basicConfigApi *apiserver.BasicConfiguration = nil

		apiForm any = nil
	)

	basicConfigDb, err = appdb.FindBasicConfigG(ctx, true)
	if err != nil {
		return nil, err
	}

	apiForm, err = ConvertDbToApiForm(basicConfigDb)
	if err != nil {
		return nil, err
	} else {
		basicConfigApi = apiForm.(*apiserver.BasicConfiguration)
	}

	return basicConfigApi, err
}

func SetBasicConfig(ctx context.Context, config *apiserver.BasicConfiguration) (*apiserver.BasicConfiguration, error) {

	var (
		err           error              = nil
		basicConfigDb *appdb.BasicConfig = nil

		apiForm any = nil
		dbForm  any = nil
	)

	if config == nil {
		return nil, errors.New("basic config is nil")
	}

	dbForm, err = ConvertDbToApiForm(config)
	if err != nil {
		return nil, err
	} else {
		basicConfigDb = dbForm.(*appdb.BasicConfig)
	}

	exists, err := basicConfigDb.Exists(ctx, getDb())
	if err != nil {
		return nil, err
	}

	if exists {
		_, err = basicConfigDb.Update(ctx, getDb(), boil.Infer())
	} else {
		err = basicConfigDb.Insert(ctx, getDb(), boil.Infer())
	}

	apiForm, err = ConvertDbToApiForm(basicConfigDb)
	if err != nil {
		return nil, err
	}

	return apiForm.(*apiserver.BasicConfiguration), err
}

func GetAdvancedConfig(ctx context.Context) (*apiserver.AdvancedConfiguration, error) {
	return nil, errors.New("not implemented")
}

func SetAdvancedConfig(ctx context.Context, config *apiserver.AdvancedConfiguration) (*apiserver.AdvancedConfiguration, error) {
	if config != nil {
		return nil, errors.New("advanced config to set is nil")
	}
	return nil, errors.New("not implemented")
}

func GetAttributeMapping(ctx context.Context) (*apiserver.AttributeMap, error) {
	return nil, errors.New("not implemented")
}

func SetAttributeMapping(ctx context.Context, mapping *apiserver.AttributeMap) (*apiserver.AttributeMap, error) {
	if mapping != nil {
		return nil, errors.New("attribute map settings to set is nil")
	}
	return nil, errors.New("not implemented")
}

func GetPermissionSettings(ctx context.Context) (*apiserver.AttributeMap, error) {
	return nil, errors.New("not implemented")
}

func SetPermissionSettings(ctx context.Context, permissions *apiserver.Permissions) (*apiserver.Permissions, error) {
	if permissions != nil {
		return nil, errors.New("permission settings to set is nil")
	}
	return nil, errors.New("not implemented")
}

func GetElionaHost() string {
	var eliDomain string

	db := getDb()
	row := db.QueryRow("SELECT domain_name FROM eliona_config ;")
	row.Scan()

	return eliDomain
}

func getDb() *sql.DB {
	return db.Database(os.Getenv("APPNAME"))
}
