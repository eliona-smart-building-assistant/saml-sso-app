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

	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Todo: Define anything for configuration like structures and methods to read and process configuration
func InsertDefaultSamlConfiguration(connection db.Connection) error {

	return errors.New("not implemented")
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

func SetBasicConfig(ctx context.Context, basicConfig *apiserver.BasicConfiguration) (*apiserver.BasicConfiguration, error) {

	var (
		err           error              = nil
		basicConfigDb *appdb.BasicConfig = nil

		apiForm any = nil
		dbForm  any = nil
	)

	if basicConfig == nil {
		return nil, errors.New("basic config is nil")
	}

	dbForm, err = ConvertDbToApiForm(basicConfig)
	if err != nil {
		return nil, err
	} else {
		basicConfigDb = dbForm.(*appdb.BasicConfig)
	}

	basicConfigDb.Insert(ctx, getDb(), boil.Infer())
}

func getDb() *sql.DB {
	return db.Database(os.Getenv("APPNAME"))
}
