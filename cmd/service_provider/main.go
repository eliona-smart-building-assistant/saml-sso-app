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

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"saml-sso/saml"

	"github.com/crewjam/saml/samlsp"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	SP_CONFIG = "./cmd/service_provider/cnf.json"
)

type SPConfig struct {
	BaseURL     string `json:"baseUrl"`
	IdpMetadata string `json:"idpMetadata"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}

func main() {
	var config SPConfig = SPConfig{}

	cnfContent, err := os.ReadFile(SP_CONFIG)
	if err != nil {
		log.Fatal("sp app", "reading config file: %v", err)
	}

	err = json.Unmarshal(cnfContent, &config)
	if err != nil {
		log.Fatal("sp app", "parsing config file: %v", err)
	}

	spUrl, err := url.Parse(config.BaseURL)
	if err != nil {
		log.Fatal("sp app", "parsing base url: %v", err)
	}

	cert, err := os.ReadFile(config.Certificate)
	if err != nil {
		log.Fatal("sp app", "read cert: %v", err)
	}
	key, err := os.ReadFile(config.PrivateKey)
	if err != nil {
		log.Fatal("sp app", "read priv key: %v", err)
	}
	metaIdp, err := os.ReadFile(config.IdpMetadata)
	if err != nil {
		log.Fatal("sp app", "read metadata: %v", err)
	}

	sp, err := saml.NewServiceProvider(string(cert), string(key), config.BaseURL,
		metaIdp)
	if err != nil {
		log.Fatal("sp app", "init service provider: %v", err)
	}

	app := http.HandlerFunc(test)
	http.Handle("/test", sp.GetMiddleWare().RequireAccount(app))
	http.Handle("/saml/", sp.GetMiddleWare())

	log.Info("sp app", "started @ %s", config.BaseURL)
	err = http.ListenAndServe(":"+spUrl.Port(), nil)

	if err != nil {
		log.Error("sp app", "exiting due to an error: %v", err)
	} else {
		log.Info("sp app", "exited")
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world! "+", %s!",
		samlsp.AttributeFromContext(r.Context(), "Name"))
}
