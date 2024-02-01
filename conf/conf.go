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

	bs, err := GetConfig(context.Background())
	if err == nil && bs != nil {
		log.Info(LOG_REGIO, "config already exists. skip insert default config.")
		return nil
	}

	var (
		basicConfig appdb.Config = appdb.Config{
			Enable:                   AUTO_CNF_DEFAULT_ENABLED,
			IdpMetadataURL:           null.StringFromPtr(nil),
			MetadataXML:              null.StringFromPtr(nil),
			UserToArchive:            AUTO_CNF_DEFAULT_USER_TO_ARCHIVE,
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

	basicConfig.OwnURL = "https://" + getElionaHost()
	certificate, privateKey, err := utils.CreateSelfsignedX509Certificate(
		AUTO_CNF_DEFAULT_CERT_VALID_DAYS, AUTO_CNF_DEFAULT_KEY_LENGTH, nil, nil)
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

func GetConfig(ctx context.Context) (*apiserver.Configuration, error) {
	basicConfigDb, err := appdb.Configs().One(ctx, getDb())
	if err != nil {
		return nil, err
	}

	basicConfigApi, err := ConfigDbToApiForm(basicConfigDb)
	if err != nil {
		return nil, err
	}

	return basicConfigApi, nil
}

func SetConfig(ctx context.Context, config *apiserver.Configuration) (*apiserver.Configuration, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	config.Id = 1 // Enforced by table definition.

	ConfigDb, err := ConfigApiToDbForm(config)
	if err != nil {
		return nil, err
	}

	exists, err := appdb.Configs().Exists(ctx, getDb())
	if err != nil {
		return nil, err
	}

	log.Debug(LOG_REGIO, "config exists %v", exists)

	if exists {
		_, err = ConfigDb.Update(ctx, getDb(),
			boil.Blacklist(appdb.ConfigColumns.ID),
		)
	} else {
		err = ConfigDb.Insert(ctx, getDb(),
			boil.Greylist(appdb.ConfigColumns.Enable, appdb.ConfigColumns.SignedRequest),
		)
	}

	if err != nil {
		return nil, err
	}

	apiForm, err := ConfigDbToApiForm(ConfigDb)
	if err != nil {
		return nil, err
	}

	return apiForm, err
}

func DeleteConfig(ctx context.Context) (int, error) {
	ans, err := appdb.Configs().DeleteAll(ctx, getDb())

	return int(ans), err
}

func GetAttributeMapping(ctx context.Context) (*apiserver.AttributeMap, error) {
	attrMapDb, err := appdb.AttributeMaps().One(ctx, getDb())
	if err != nil {
		return nil, err
	}

	attrMapApi, err := AttributeMapDbToApiForm(attrMapDb)
	if err != nil {
		return nil, err
	}

	return attrMapApi, err
}

func SetAttributeMapping(ctx context.Context, mapping *apiserver.AttributeMap) (*apiserver.AttributeMap, error) {
	if mapping == nil {
		return nil, errors.New("attribute mapping is nil")
	}
	mapping.Id = 1 // Enforced by table definition.

	attributeMappingDb, err := AttributeMapApiToDbForm(mapping)
	if err != nil {
		return nil, err
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

	if err != nil {
		return nil, err
	}

	apiForm, err := AttributeMapDbToApiForm(attributeMappingDb)
	if err != nil {
		return nil, err
	}

	return apiForm, err
}

func DeleteAttributeMapping(ctx context.Context) (int, error) {
	ans, err := appdb.AttributeMaps().DeleteAll(ctx, getDb())

	return int(ans), err
}

func GetPermissionMapping(ctx context.Context) (*apiserver.Permissions, error) {
	permDb, err := appdb.Permissions().One(ctx, getDb())
	if err != nil {
		return nil, err
	}

	permApi, err := PermissionDbToApiForm(permDb)
	if err != nil {
		return nil, err
	}

	return permApi, err
}

func SetPermissionMapping(ctx context.Context, permissions *apiserver.Permissions) (*apiserver.Permissions, error) {
	if permissions == nil {
		return nil, errors.New("permissions nil")
	} else {
		permissions.Id = 1 // set id to 1, because it must
	}

	permissionsDb, err := PermissionApiToDbForm(permissions)
	if err != nil {
		return nil, err
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

	if err != nil {
		return nil, err
	}

	apiForm, err := PermissionDbToApiForm(permissionsDb)
	if err != nil {
		return nil, err
	}

	return apiForm, err
}

func DeletePermissionMapping(ctx context.Context) (int, error) {
	ans, err := appdb.Permissions().DeleteAll(ctx, getDb())

	return int(ans), err
}

func DeleteAllConfigurations(ctx context.Context) error {
	_, err := DeletePermissionMapping(ctx)
	if err != nil {
		return err
	}
	_, err = DeleteAttributeMapping(ctx)
	if err != nil {
		return err
	}
	_, err = DeleteConfig(ctx)
	return err
}

func getElionaHost() string {
	var eliDomain string

	db := getDb()
	row := db.QueryRow("SELECT domain_name FROM eliona_config ;")
	err := row.Scan(&eliDomain)
	if err != nil {
		log.Error(LOG_REGIO, "scan getElionaHost: %v", err)
	}

	return eliDomain
}

func getDb() *sql.DB {
	return db.Database(app.AppName())
}

//
// functions for test purposes only!

func DropOwnSchema() error {
	log.Warn(LOG_REGIO, "Do not use function DropOwnSchema() in the Application. "+
		"Only for test purposes!")
	db := getDb()
	_, err := db.Exec("DROP SCHEMA IF EXISTS saml_sp CASCADE")
	return err
}
func UserLeicomInit() error {
	log.Warn(LOG_REGIO, "Do not use function UserLeicomInit() in the Application. "+
		"Only for test purposes!")

	var (
		userLeicom string = "leicom"
		userRet    string
	)

	db := getDb()
	// check, if user exists (do not exist only on in the mock db)
	row := db.QueryRow("SELECT usename FROM pg_user WHERE usename = $1", userLeicom)
	if err := row.Err(); err != nil {
		return err
	}
	err := row.Scan(&userRet)

	if err != nil || userRet != userLeicom {
		// user not exist
		_, err = db.Exec("CREATE USER leicom NOLOGIN")
	}

	return err
}
