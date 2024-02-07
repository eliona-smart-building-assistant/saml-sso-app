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
	"os"
	"saml-sso/eliona"
	"testing"

	"github.com/eliona-smart-building-assistant/go-utils/log"
)

func TestApp_Others(t *testing.T) {

	// this test needs a real eliona db to due missing tables
	// in the test db
	_, realDb := os.LookupEnv("REAL_DB")
	if !realDb {
		log.Warn("TestApp_Others", "test disabled because missing env var REAL_DB")
		t.Log("TestApp_Others: test disabled because missing env var REAL_DB")
		return
	}

	token, err := eliona.GetElionaJsonWebToken("su#@eliona.io")
	if err != nil {
		t.Error(err)
	}
	if token != nil {
		log.Debug("test eliona/others", "jwt: %v", *token)
	} else {
		t.Error("token of sv is nil")
	}
}
