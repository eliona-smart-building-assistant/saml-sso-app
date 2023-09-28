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

package eliona

import (
	"context"
	"errors"
	"io"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

type EliApiV2 struct {
	authCtx context.Context
	client  *api.APIClient
}

func NewEliApiV2() *EliApiV2 {
	return &EliApiV2{
		authCtx: client.AuthenticationContext(),
		client:  client.NewClient(),
	}
}

func (e *EliApiV2) GetApiVersion() (map[string]interface{}, error) {
	ver, _, err := e.client.VersionAPI.GetVersion(e.authCtx).Execute()

	return ver, err
}

func (e *EliApiV2) GetUserIfExists(email string) (*api.User, error) {
	users, resp, err := e.client.UsersAPI.GetUsers(e.authCtx).Execute()
	if err != nil {
		msg := err.Error()

		if resp != nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				msg += string(body)
			}
		}
		log.Error(LOG_REGIO, msg)
		return nil, err
	}

	for _, user := range users {
		if user.Email == email {
			return &user, err
		}
	}

	return nil, errors.New("user not exist")
}

func (e *EliApiV2) AddUser(user *api.User) (*api.User, error) {
	if user == nil {
		return nil, errors.New("user to add is nil")
	}
	userRet, resp, err := e.client.UsersAPI.PutUser(e.authCtx).User(*user).Execute()
	if err != nil {
		msg := err.Error()
		if resp != nil {
			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				msg += string(respBody)
			}
		}
		log.Error(LOG_REGIO, msg)
	}
	return userRet, err
}
