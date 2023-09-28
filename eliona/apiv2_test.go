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

package eliona_test

import (
	"saml-sso/eliona"
	"testing"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/go-test/deep"
)

// need a running apiv2 and exported API_TOKEN / API_ENDPOINT

func TestApp_EliApi_Version(t *testing.T) {
	eApi := eliona.NewEliApiV2()
	ver, err := eApi.GetApiVersion()
	if err != nil {
		t.Error("get version, ", err)
	}
	t.Log("APIv2 Version: ", ver)
}

func TestApp_EliApi_AddUser(t *testing.T) {
	fistName := "myFirstName"
	lastName := "myLastName"
	email := "my@email.net"

	eApi := eliona.NewEliApiV2()
	userToAdd := api.User{
		Firstname: *api.NewNullableString(&fistName),
		Lastname:  *api.NewNullableString(&lastName),
		Email:     email,
	}
	addedUser, err := eApi.AddUser(&userToAdd)

	userToAdd.Id = addedUser.Id

	if diff := deep.Equal(*addedUser, userToAdd); diff != nil {
		t.Error("added user has diff to returned one: ", diff)
	}

	if err != nil {
		t.Error("add user, ", err)
	}
}

func TestApp_EliApi_CheckUserExists(t *testing.T) {

	eApi := eliona.NewEliApiV2()
	myUser, err := eApi.GetUserIfExists("my@email.net")
	if err != nil || myUser == nil {
		t.Error(err)
	}
	if *myUser.Firstname.Get() != "myFirstName" || *myUser.Lastname.Get() != "myLastName" ||
		myUser.Email != "my@email.net" {
		t.Error("wrong user content")
	}

	noUser, err := eApi.GetUserIfExists("absolutUnrealisticUser@nobodyknows.anywhere")
	if err == nil || noUser != nil {
		t.Error("this user should not exists. ;)", err)
	}
}
