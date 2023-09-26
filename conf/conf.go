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

package conf

import (
	"context"
	"database/sql"
	"errors"

	"saml-sso/apiserver"
	"saml-sso/appdb"
	"saml-sso/utils"

	"github.com/eliona-smart-building-assistant/go-eliona/app"
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
	AUTO_CNF_DEFAULT_USER_TO_ARCHIVE   = false
	AUTO_CNF_DEFAULT_ALLOW_INIT_BY_IDP = false
	AUTO_CNF_DEFAULT_SIGNING_REQ       = true
	AUTO_CNF_DEFAULT_FORCE_AUTHN       = false
	AUTO_CNF_DEFAULT_COOKIE_SECURE     = false
	AUTO_CNF_DEFAULT_ENTITY_ID         = utils.UTILS_OWN_URL_PLACEHOLDER + "/saml/metadata"
	AUTO_CNF_DEFAULT_LOGIN_FAIL_URL    = utils.UTILS_OWN_URL_PLACEHOLDER + "/noLogin"
	AUTO_CNF_DEFAULT_USERNAME_ATTR     = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"
	AUTO_CNF_DEFAULT_SYSTEM_PERMISSION = "regular"
	AUTO_CNF_DEFAULT_PROJ_PERMISSION   = "operator"
	AUTO_CNF_DEFAULT_CERT_VALID_DAYS   = 3650
	AUTO_CNF_DEFAULT_KEY_LENGTH        = 4096
)

func InsertAutoSamlConfiguration(ctx context.Context) error {

	bs, err := GetBasicConfig(context.Background())
	if err == nil && bs != nil {
		log.Info(LOG_REGIO, "config already exists. skip insert default config.")
		return nil
	}

	var (
		basicConfig appdb.BasicConfig = appdb.BasicConfig{
			Enable:         AUTO_CNF_DEFAULT_ENABLED,
			IdpMetadataURL: null.StringFromPtr(nil),
			MetadataXML:    null.StringFromPtr(nil),
			UserToArchive:  AUTO_CNF_DEFAULT_USER_TO_ARCHIVE,
		}
		advancedConfig appdb.AdvancedConfig = appdb.AdvancedConfig{
			AllowInitializationByIdp: AUTO_CNF_DEFAULT_ALLOW_INIT_BY_IDP,
			SignedRequest:            AUTO_CNF_DEFAULT_SIGNING_REQ,
			ForceAuthn:               AUTO_CNF_DEFAULT_FORCE_AUTHN,
			EntityID:                 AUTO_CNF_DEFAULT_ENTITY_ID,
			CookieSecure:             AUTO_CNF_DEFAULT_COOKIE_SECURE,
			LoginFailedURL:           AUTO_CNF_DEFAULT_LOGIN_FAIL_URL,
		}
		attributeMapping appdb.AttributeMap = appdb.AttributeMap{
			Email:     AUTO_CNF_DEFAULT_USERNAME_ATTR,
			FirstName: null.StringFromPtr(nil),
			LastName:  null.StringFromPtr(nil),
			Phone:     null.StringFromPtr(nil),
		}
		permissionCnf appdb.Permission = appdb.Permission{
			DefaultSystemRole:       AUTO_CNF_DEFAULT_SYSTEM_PERMISSION,
			DefaultProjRole:         AUTO_CNF_DEFAULT_PROJ_PERMISSION,
			SystemRoleSamlAttribute: null.StringFromPtr(nil),
			SystemRoleMap:           null.JSONFromPtr(nil),
			ProjRoleSamlAttribute:   null.StringFromPtr(nil),
			ProjRoleMap:             null.JSONFromPtr(nil),
		}
	)

	basicConfig.OwnURL = "https://" + GetElionaHost()
	certificate, privateKey, err := utils.CreateSelfsignedX509Certificate(
		AUTO_CNF_DEFAULT_CERT_VALID_DAYS, AUTO_CNF_DEFAULT_KEY_LENGTH, nil, nil)
	if err != nil {
		log.Error(LOG_REGIO, "auto configuration generate x509 certificates: %v",
			err)
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

	basicConfigDb, err = appdb.BasicConfigs().One(ctx, getDb())
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

func SetBasicConfig(ctx context.Context, config *apiserver.BasicConfiguration) (
	*apiserver.BasicConfiguration, error) {

	var (
		err           error              = nil
		basicConfigDb *appdb.BasicConfig = nil

		apiForm any = nil
		dbForm  any = nil
	)

	if config == nil {
		return nil, errors.New("basic config is nil")
	} else {
		config.Id = 1 // set id to 1, because it must
	}

	dbForm, err = ConvertApiToDbForm(config)
	if err != nil {
		return nil, err
	} else {
		basicConfigDb = dbForm.(*appdb.BasicConfig)
	}

	exists, err := appdb.BasicConfigs().Exists(ctx, getDb())
	if err != nil {
		return nil, err
	}

	log.Debug(LOG_REGIO, "basic config exists %v", exists)

	if exists {
		_, err = basicConfigDb.Update(ctx, getDb(),
			boil.Blacklist(appdb.BasicConfigColumns.ID))
	} else {
		err = basicConfigDb.Insert(ctx, getDb(),
			boil.Greylist(appdb.BasicConfigColumns.Enable))
	}

	apiForm, err = ConvertDbToApiForm(basicConfigDb)
	if err != nil {
		return nil, err
	}

	return apiForm.(*apiserver.BasicConfiguration), err
}

func DeleteBasicConfig(ctx context.Context) (int, error) {
	ans, err := appdb.BasicConfigs().DeleteAll(ctx, getDb())

	return int(ans), err
}

func GetAdvancedConfig(ctx context.Context) (*apiserver.AdvancedConfiguration, error) {

	var (
		err          error                            = nil
		advConfigDb  *appdb.AdvancedConfig            = nil
		advConfigApi *apiserver.AdvancedConfiguration = nil

		apiForm any = nil
	)

	advConfigDb, err = appdb.AdvancedConfigs().One(ctx, getDb())
	if err != nil {
		return nil, err
	}

	apiForm, err = ConvertDbToApiForm(advConfigDb)
	if err != nil {
		return nil, err
	} else {
		advConfigApi = apiForm.(*apiserver.AdvancedConfiguration)
	}

	return advConfigApi, err
}

func SetAdvancedConfig(ctx context.Context, config *apiserver.AdvancedConfiguration) (
	*apiserver.AdvancedConfiguration, error) {

	var (
		err              error                 = nil
		advancedConfigDb *appdb.AdvancedConfig = nil

		apiForm any = nil
		dbForm  any = nil
	)

	if config == nil {
		return nil, errors.New("advanced config is nil")
	} else {
		config.Id = 1 // set id to 1, because it must
	}

	dbForm, err = ConvertApiToDbForm(config)
	if err != nil {
		return nil, err
	} else {
		advancedConfigDb = dbForm.(*appdb.AdvancedConfig)
	}

	exists, err := appdb.AdvancedConfigs().Exists(ctx, getDb())
	if err != nil {
		return nil, err
	}

	log.Debug(LOG_REGIO, "advanced config exists %v", exists)

	if exists {
		_, err = advancedConfigDb.Update(ctx, getDb(),
			boil.Greylist(appdb.AdvancedConfigColumns.SignedRequest))
	} else {
		err = advancedConfigDb.Insert(ctx, getDb(),
			boil.Greylist(appdb.AdvancedConfigColumns.SignedRequest))
	}

	apiForm, err = ConvertDbToApiForm(advancedConfigDb)
	if err != nil {
		return nil, err
	}

	return apiForm.(*apiserver.AdvancedConfiguration), err
}

func DeleteAdvancedConfig(ctx context.Context) (int, error) {
	ans, err := appdb.AdvancedConfigs().DeleteAll(ctx, getDb())

	return int(ans), err
}

func GetAttributeMapping(ctx context.Context) (*apiserver.AttributeMap, error) {

	var (
		err        error                   = nil
		attrMapDb  *appdb.AttributeMap     = nil
		attrMapApi *apiserver.AttributeMap = nil

		apiForm any = nil
	)

	attrMapDb, err = appdb.AttributeMaps().One(ctx, getDb())
	if err != nil {
		return nil, err
	}

	apiForm, err = ConvertDbToApiForm(attrMapDb)
	if err != nil {
		return nil, err
	} else {
		attrMapApi = apiForm.(*apiserver.AttributeMap)
	}

	return attrMapApi, err
}

func SetAttributeMapping(ctx context.Context, mapping *apiserver.AttributeMap) (
	*apiserver.AttributeMap, error) {

	var (
		err                error               = nil
		attributeMappingDb *appdb.AttributeMap = nil

		apiForm any = nil
		dbForm  any = nil
	)

	if mapping == nil {
		return nil, errors.New("attribute mapping is nil")
	} else {
		mapping.Id = 1 // set id to 1, because it must
	}

	dbForm, err = ConvertApiToDbForm(mapping)
	if err != nil {
		return nil, err
	} else {
		attributeMappingDb = dbForm.(*appdb.AttributeMap)
	}

	exists, err := appdb.AttributeMaps().Exists(ctx, getDb())
	if err != nil {
		return nil, err
	}

	log.Debug(LOG_REGIO, "attribute map exists %v", exists)

	if exists {
		_, err = attributeMappingDb.Update(ctx, getDb(),
			boil.Blacklist(appdb.AttributeMapColumns.ID))
	} else {
		err = attributeMappingDb.Insert(ctx, getDb(),
			boil.Blacklist(appdb.AttributeMapColumns.ID))
	}

	apiForm, err = ConvertDbToApiForm(attributeMappingDb)
	if err != nil {
		return nil, err
	}

	return apiForm.(*apiserver.AttributeMap), err
}

func DeleteAttributeMapping(ctx context.Context) (int, error) {
	ans, err := appdb.AttributeMaps().DeleteAll(ctx, getDb())

	return int(ans), err
}

func GetPermissionSettings(ctx context.Context) (*apiserver.Permissions, error) {

	var (
		err     error                  = nil
		permDb  *appdb.Permission      = nil
		permApi *apiserver.Permissions = nil

		apiForm any = nil
	)

	permDb, err = appdb.Permissions().One(ctx, getDb())
	if err != nil {
		return nil, err
	}

	apiForm, err = ConvertDbToApiForm(permDb)
	if err != nil {
		return nil, err
	} else {
		permApi = apiForm.(*apiserver.Permissions)
	}

	return permApi, err
}

func SetPermissionSettings(ctx context.Context, permissions *apiserver.Permissions) (
	*apiserver.Permissions, error) {

	var (
		err           error             = nil
		permissionsDb *appdb.Permission = nil

		apiForm any = nil
		dbForm  any = nil
	)

	if permissions == nil {
		return nil, errors.New("permissions nil")
	} else {
		permissions.Id = 1 // set id to 1, because it must
	}

	dbForm, err = ConvertApiToDbForm(permissions)
	if err != nil {
		return nil, err
	} else {
		permissionsDb = dbForm.(*appdb.Permission)
	}

	exists, err := appdb.Permissions().Exists(ctx, getDb())
	if err != nil {
		return nil, err
	}

	log.Debug(LOG_REGIO, "permissions exists %v", exists)

	if exists {
		_, err = permissionsDb.Update(ctx, getDb(),
			boil.Blacklist(appdb.PermissionColumns.ID))
	} else {
		err = permissionsDb.Insert(ctx, getDb(),
			boil.Blacklist(appdb.PermissionColumns.ID))
	}

	apiForm, err = ConvertDbToApiForm(permissionsDb)
	if err != nil {
		return nil, err
	}

	return apiForm.(*apiserver.Permissions), err
}

func DeletePermissionSettings(ctx context.Context) (int, error) {
	ans, err := appdb.Permissions().DeleteAll(ctx, getDb())

	return int(ans), err
}

func DeleteAllConfigurations(ctx context.Context) error {
	_, err := DeletePermissionSettings(ctx)
	if err != nil {
		return err
	}
	_, err = DeleteAttributeMapping(ctx)
	if err != nil {
		return err
	}
	_, err = DeleteAdvancedConfig(ctx)
	if err != nil {
		return err
	}
	_, err = DeleteBasicConfig(ctx)

	return err
}

func GetElionaHost() string {
	var eliDomain string

	db := getDb()
	row := db.QueryRow("SELECT domain_name FROM eliona_config ;")
	row.Scan(&eliDomain)

	return eliDomain
}

func DropOwnSchema() error {
	db := getDb()
	_, err := db.Exec("DROP SCHEMA IF EXISTS saml_sp CASCADE")
	return err
}

func getDb() *sql.DB {
	return db.Database(app.AppName())
}
