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
	"fmt"
	"strings"

	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

type ElionaJwt struct {
	Jwt string
}

const (
	OTHERS_GET_JWT_QUERY_V10 = "(SELECT public.make_jwt(jwt,secret) " +
		"FROM  public.eliona_user u JOIN public.eliona_secret " +
		"USING (schema), public.claim_jwt(role, now() + validity,user_id,null) jwt " +
		"WHERE lower(u.email) = lower($1) AND NOT u.archived)"

	OTHERS_GET_JWT_QUERY_V11 = "(SELECT public.make_jwt(jwt,secret) " +
		"FROM  public.eliona_user u " +
		"JOIN public.acl_role r ON (u.role_id = r.role_id) " +
		"JOIN public.eliona_secret " +
		"USING (schema), public.claim_jwt(role, now() + validity,user_id,null) jwt " +
		"WHERE lower(u.email) = lower($1) AND NOT u.archived)"

	OTHERS_GET_JWT_QUERY_V12 = "(SELECT public.make_jwt(jwt,secret) " +
		"FROM  public.eliona_user u " +
		"JOIN public.acl_role r ON (u.role_id = r.role_id) " +
		"JOIN public.project_user USING (user_id) " +
		"JOIN public.eliona_secret USING (schema), " +
		"public.claim_jwt(role, now() + validity,user_id,proj_id,r.role_id::text,schema) jwt " +
		"WHERE lower(u.email) = lower($1) AND NOT u.archived)"
)

func GetElionaJsonWebToken(email string) (*string, error) {

	var (
		err      error
		version  string
		jwt      ElionaJwt = ElionaJwt{}
		jwtQuery string
	)

	db := getDb()

	// find version
	row := db.QueryRow("SELECT version FROM versioning.latest_version WHERE app_name = 'public'")
	if row == nil {
		return nil, row.Err()
	}
	err = row.Scan(&version)
	if row == nil {
		return nil, err
	}

	// before v10 docker matching images are available.
	if strings.Contains(version, "v10.") {
		log.Debug(LOG_REGIO, "eliona v10")
		jwtQuery = OTHERS_GET_JWT_QUERY_V10
	} else if strings.Contains(version, "v11.") &&
		!strings.Contains(version, "v11.1.5") {
		log.Debug(LOG_REGIO, "eliona v11")
		jwtQuery = OTHERS_GET_JWT_QUERY_V11
	} else {
		// assume, that the version is newer (with ACL)
		jwtQuery = OTHERS_GET_JWT_QUERY_V12
	}

	row = db.QueryRow(jwtQuery, email)

	if row == nil {
		return nil, errors.New("returned row is nil")
	} else if row.Err() != nil {
		return nil, row.Err()
	}

	err = row.Scan(&jwt.Jwt)

	return &jwt.Jwt, err
}

func UpdateElionaUserArchivedPhone(email string, phone *string, archived bool) error {

	_, err := getDb().Exec("UPDATE eliona_user SET archived = $1, phone = $2 WHERE email = $3",
		archived, phone, email)

	return err
}

func GetFirstProjectId() (projectId string, err error) {
	row := getDb().QueryRow("SELECT proj_id FROM eliona_project ORDER BY proj_id LIMIT 1")

	if err = row.Err(); err != nil {
		return
	}

	err = row.Scan(&projectId)
	return
}

func GetACLRoleMap() (roleMap map[string]int, err error) {

	roleMap = make(map[string]int)

	rows, err := getDb().Query("SELECT role_id, displayname FROM acl_role")

	if err != nil {
		err = fmt.Errorf("cannot get acl roles: %v", err)
		return
	}

	for rows.Next() {
		var (
			roleId       int
			roleDispName string
		)
		err = rows.Scan(&roleId, &roleDispName)
		if err != nil {
			return
		}
		roleMap[roleDispName] = roleId
	}

	return
}

func SetUserPermissions(userId *string, systemRoleId int, language string) (err error) {
	_, err = getDb().Exec("UPDATE public.eliona_user SET role_id = $1, language = $2 WHERE user_id = $3",
		systemRoleId, language, userId)
	return
}

func SetProjectUser(projectId string, userId *string, roleId int) (err error) {

	_, err = getDb().Exec("INSERT INTO project_user (proj_id, user_id, role_id) "+
		"VALUES ($1, $2, $3)", projectId, userId, roleId)

	return
}

func getDb() *sql.DB {
	return db.Database(app.AppName())
}
