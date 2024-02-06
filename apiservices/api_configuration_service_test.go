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

package apiservices_test

import (
	"context"
	"database/sql"
	"fmt"
	"saml-sso/apiservices"
	"saml-sso/conf"
	"testing"

	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-utils/db"
)

func TestApp_AppApi_Configuration_InitDB(t *testing.T) {
	err := conf.UserLeicomInit()
	if err != nil {
		t.Log("user leicom, ", err)
	}
	err = conf.DropOwnSchema()
	if err != nil {
		// no error, if schema not exist
		t.Log("drop schema, ", err)
	}

	execFunc := app.ExecSqlFile("../conf/init.sql")
	err = execFunc(db.NewConnection())
	if err != nil {
		t.Error("init.sql failed, ", err)
	}
}

func TestApp_AppApi_Configuration_GetBasicConfig(t *testing.T) {
	apiService := apiservices.NewConfigurationApiService()
	cnf, err := apiService.GetBasicConfiguration(context.Background())
	if err != sql.ErrNoRows {
		t.Error(err)
	}
	fmt.Println(cnf)
}

func TestApp_AppApi_Configuration_GetAdvancedConfig(t *testing.T) {

}

func TestApp_AppApi_Configuration_GetAttributeMapping(t *testing.T) {

}

func TestApp_AppApi_Configuration_GetPermissionMap(t *testing.T) {

}

func TestApp_AppApi_Configuration_PutBasicConfig(t *testing.T) {

}

func TestApp_AppApi_Configuration_PutAdvancedConfig(t *testing.T) {

}

func TestApp_AppApi_Configuration_PutAttributeMapping(t *testing.T) {

}

func TestApp_AppApi_Configuration_PutPermissionMap(t *testing.T) {

}
