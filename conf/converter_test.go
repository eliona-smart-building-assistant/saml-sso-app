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
	"encoding/json"
	"errors"
	"saml-sso/apiserver"
	"saml-sso/appdb"
	"saml-sso/conf"
	"saml-sso/utils"
	"testing"

	"github.com/go-test/deep"
)

func TestApp_Converter_ConvertConfig(t *testing.T) {
	for i := 0; i < 5; i++ {
		apiBasicCnf := utils.CreateRandomApiConfig()
		dbBasicCnf, err := conf.ConfigApiToDbForm(&apiBasicCnf)
		if err != nil {
			t.Error("convert basic config from api to db form using general converter: ", err)
		}
		err = compareConfig(dbBasicCnf, &apiBasicCnf)
		if err != nil {
			t.Error("convert basic config from api to db wrong content: ", err)
		}

		apiBasicCnfReturned, err := conf.ConfigDbToApiForm(dbBasicCnf)
		if err != nil {
			t.Error("convert basic config from db to api form using general converter: ", err)
		}
		diff := deep.Equal(apiBasicCnfReturned, &apiBasicCnf)
		if diff != nil {
			t.Error("convert basic config from db to api wrong content (not origin): ", err)
		}
	}
}

func TestApp_Converter_ConvertAttributeMapping(t *testing.T) {
	for i := 0; i < 5; i++ {
		apiAttrMap := utils.CreateRandomApiAttributeMap()
		dbAttrMap, err := conf.AttributeMapApiToDbForm(&apiAttrMap)
		if err != nil {
			t.Error("convert attribute map from api to db form using general converter: ", err)
		}
		err = compareAttributeMaps(dbAttrMap, &apiAttrMap)
		if err != nil {
			t.Error("convert attribute map config from api to db wrong content: ", err)
		}

		apiAttrMapReturned, err := conf.AttributeMapDbToApiForm(dbAttrMap)
		if err != nil {
			t.Error("convert attribute map config from db to api form using general converter: ", err)
		}
		diff := deep.Equal(apiAttrMapReturned, &apiAttrMap)
		if diff != nil {
			t.Error("convert attribute map config from db to api wrong content (not origin): ", err)
		}
	}
}

func TestApp_Converter_ConvertPermissionCnf(t *testing.T) {
	for i := 0; i < 5; i++ {
		apiPerm := utils.CreateRandomApiPermissions()
		dbPerm, err := conf.PermissionApiToDbForm(&apiPerm)
		if err != nil {
			t.Error("convert permissions from api to db form using general converter: ", err)
		}
		err = comparePermissions(dbPerm, &apiPerm)
		if err != nil {
			t.Error("convert permissions config from api to db wrong content: ", err)
		}

		apiPermReturned, err := conf.PermissionDbToApiForm(dbPerm)
		if err != nil {
			t.Error("convert permissions config from db to api form using general converter: ", err)
		}
		diff := deep.Equal(apiPermReturned, &apiPerm)
		if diff != nil {
			t.Error("convert permissions config from db to api wrong content (not origin): ", diff)
		}
	}
}

func compareConfig(db *appdb.Config, api *apiserver.Configuration) error {
	if db.ID != api.Id {
		return errors.New("id")
	}

	if db.Enable != api.Enable {
		return errors.New("enable")
	}

	if deep.Equal(db.IdpMetadataURL.Ptr(), api.IdpMetadataUrl) != nil {
		return errors.New("metadata url")
	}

	if deep.Equal(db.MetadataXML.Ptr(), api.IdpMetadataXml) != nil {
		return errors.New("metadata xml")
	}

	if db.OwnURL != api.OwnUrl {
		return errors.New("own url")
	}

	if db.SPCertificate != api.ServiceProviderCertificate {
		return errors.New("certificate")
	}

	if db.SPPrivateKey != api.ServiceProviderPrivateKey {
		return errors.New("private key")
	}

	if db.UserToArchive != api.UserToArchive {
		return errors.New("user to archive")
	}

	if db.AllowInitializationByIdp != api.AllowInitializationByIdp {
		return errors.New("allow initialisation by idp")
	}

	if db.CookieSecure != api.CookieSecure {
		return errors.New("cookie secure")
	}

	if db.EntityID != api.EntityId {
		return errors.New("entity id")
	}

	if db.ForceAuthn != api.ForceAuthn {
		return errors.New("force authn")
	}

	if db.LoginFailedURL != api.LoginFailedUrl {
		return errors.New("login failed url")
	}

	if db.SignedRequest != api.SignedRequest {
		return errors.New("signed request")
	}

	return nil
}

func compareAttributeMaps(db *appdb.AttributeMap, api *apiserver.AttributeMap) error {
	if db.Email != api.Email {
		return errors.New("email")
	}
	if db.ID != api.Id {
		return errors.New("id")
	}
	if deep.Equal(db.FirstName.Ptr(), api.FirstName) != nil {
		return errors.New("first name")
	}
	if deep.Equal(db.LastName.Ptr(), api.LastName) != nil {
		return errors.New("last name")
	}
	if deep.Equal(db.Phone.Ptr(), api.Phone) != nil {
		return errors.New("phone number")
	}
	return nil
}

func comparePermissions(db *appdb.Permission, api *apiserver.Permissions) error {

	if db.DefaultProjRole != api.DefaultProjRole {
		return errors.New("default project role")
	}

	if db.DefaultSystemRole != api.DefaultSystemRole {
		return errors.New("default system role")
	}

	if db.ID != api.Id {
		return errors.New("id")
	}

	dbJson := []apiserver.RoleMap{}
	if !db.ProjRoleMap.IsZero() {
		dbJsonB := db.ProjRoleMap.JSON
		err := json.Unmarshal(dbJsonB, &dbJson)
		if err != nil {
			return err
		}
	}
	if !(db.ProjRoleMap.Ptr() == nil && api.ProjRoleMap == nil) && (deep.Equal(&dbJson, api.ProjRoleMap) != nil) {
		return errors.New("project role map")
	}

	if deep.Equal(db.ProjRoleSamlAttribute.Ptr(), api.ProjRoleSamlAttribute) != nil {
		return errors.New("project role saml attribute")
	}

	dbJson = []apiserver.RoleMap{}
	if !db.SystemRoleMap.IsZero() {
		dbJsonB := db.SystemRoleMap.JSON
		err := json.Unmarshal(dbJsonB, &dbJson)
		if err != nil {
			return err
		}
	}
	if !(db.SystemRoleMap.Ptr() == nil && api.SystemRoleMap == nil) && (deep.Equal(&dbJson, api.SystemRoleMap) != nil) {
		return errors.New("system role map")
	}

	if deep.Equal(db.SystemRoleSamlAttribute.Ptr(), api.SystemRoleSamlAttribute) != nil {
		return errors.New("system role saml attribute")
	}

	return nil
}
