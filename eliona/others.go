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

// @note: contains functions, not covered over APIv2

package eliona

import (
	"database/sql"
	"errors"

	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-utils/db"
)

type ElionaJwt struct {
	Jwt string
}

func GetElionaJsonWebToken(email string) (*string, error) {

	var jwt ElionaJwt = ElionaJwt{}

	row := getDb().QueryRow("(SELECT public.make_jwt(jwt,secret) "+
		"FROM  public.eliona_user u JOIN public.eliona_secret "+
		"USING (schema), public.claim_jwt(role, now() + validity,user_id,null) jwt "+
		"WHERE lower(u.email) = lower($1) AND NOT u.archived)", email)

	if row == nil {
		return nil, errors.New("returned row is nil")
	} else if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(&jwt.Jwt)

	return &jwt.Jwt, err
}

func UpdateElionaUserArchivedPhone(email string, phone *string, archived bool) error {

	_, err := getDb().Exec("UPDATE eliona_user SET archived = $1, phone = $2 WHERE email = $3",
		archived, phone, email)

	return err
}

func getDb() *sql.DB {
	return db.Database(app.AppName())
}
