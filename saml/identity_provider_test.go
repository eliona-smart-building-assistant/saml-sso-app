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

package saml_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"saml-sso/saml"
	"saml-sso/utils"
	"strconv"
	"strings"
	"testing"

	"github.com/crewjam/saml/samlsp"
	"golang.org/x/net/html"
)

const (
	TEST_IDP_SCHEME        = "http://"
	TEST_IDP_HOST          = "localhost"
	TEST_IDP_PORT_SP       = 8080
	TEST_SP_RESOURCE_EP    = "/test"
	TEST_IDP_USER          = "eliona"
	TEST_IDP_PW            = "pwEliona"
	TEST_IDP_USER_MAIL     = "eliona@eliona.io"
	TEST_IDP_USER_NAME     = "Eliona"
	TEST_IDP_USER_LASTNAME = "Anoile"
	TEST_RESOURCE_CONTENT  = "Hello SAML user!"
)

func TestIdentityProvider(t *testing.T) {

	idp, err := saml.NewIdentityProvider(TEST_IDP_SCHEME+TEST_IDP_HOST+":"+
		strconv.Itoa(saml.IDENTITY_PROVIDER_PORT), nil, nil)
	if err != nil {
		t.Log(err)
	}

	err = idp.AddUser(TEST_IDP_USER, TEST_IDP_PW, TEST_IDP_USER_MAIL,
		TEST_IDP_USER_NAME, TEST_IDP_USER_LASTNAME)
	if err != nil {
		t.Error(err)
	}

	go idp.ServeForever()

	responseIdpMeta, err := http.Get(TEST_IDP_SCHEME + TEST_IDP_HOST + ":" +
		strconv.Itoa(saml.IDENTITY_PROVIDER_PORT) + saml.IDENTITY_PROVIDER_METADATA_PATH)
	if err != nil {
		t.Error(err)
	}

	bodyIdpMeta, err := ioutil.ReadAll(responseIdpMeta.Body)
	responseIdpMeta.Body.Close()

	if !strings.Contains(string(bodyIdpMeta), "<EntityDescriptor") || err != nil {
		t.Error("get identity providers metadata, ", err)
	}

	cert, key, err := utils.CreateSelfsignedX509Certificate(2, 2048, nil, nil)
	if err != nil {
		t.Error("create cert for sp")
	}

	sp, err := saml.NewServiceProvider(cert, key, TEST_IDP_SCHEME+TEST_IDP_HOST+":"+
		strconv.Itoa(TEST_IDP_PORT_SP), bodyIdpMeta)
	if err != nil {
		t.Error("setup sp: ", err)
	}

	app := http.HandlerFunc(test)
	http.Handle(TEST_SP_RESOURCE_EP, sp.GetMiddleWare().RequireAccount(app))
	http.Handle("/saml/", sp.GetMiddleWare())
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(TEST_IDP_PORT_SP), nil)
		t.Error(err)
	}()

	// idp and sp up, add sp to idp
	responseSpMeta, err := http.Get(TEST_IDP_SCHEME + TEST_IDP_HOST + ":" +
		strconv.Itoa(TEST_IDP_PORT_SP) + "/saml/" + saml.IDENTITY_PROVIDER_METADATA_PATH)
	if err != nil {
		t.Error(err)
	}
	bodySpMeta, err := ioutil.ReadAll(responseSpMeta.Body)
	if err != nil {
		t.Error(err)
	}
	responseSpMeta.Body.Close()
	err = idp.AddServiceProvider("http://localhost:8080/saml/metadata", string(bodySpMeta))

	// call resource
	jar, err := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	req, err := http.NewRequest("GET", TEST_IDP_SCHEME+TEST_IDP_HOST+":"+
		strconv.Itoa(TEST_IDP_PORT_SP)+TEST_SP_RESOURCE_EP, nil)
	if err != nil {
		t.Error(err)
	}
	response, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	response.Body.Close()

	// should be redirected to sso login page
	if !strings.Contains(string(body), "action=\"http://localhost:8000/sso") {
		t.Error("not redirected to the sso loginpage")
	}
	// login
	actionLoginUrl, samlReq, relayState, err := extractFormAction(string(body))
	if err != nil {
		t.Error("cannot extract login url")
	}
	formData := url.Values{
		"user":        {TEST_IDP_USER},
		"password":    {TEST_IDP_PW},
		"SAMLRequest": {samlReq},
		"RelayState":  {relayState},
	}

	response, err = client.PostForm(actionLoginUrl, formData)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Error("login @ ", actionLoginUrl)
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	response.Body.Close()

	// send the saml response from the IdP to the ACS of SP to login
	actionURL, samlResponse, err := extextractFormSamlResponse(string(body))
	formData = url.Values{
		"SAMLResponse": {samlResponse},
	}
	response, err = client.PostForm(actionURL, formData)
	if err != nil {
		t.Error(err)
	}
	// if response.StatusCode != http.StatusOK {
	// 	t.Error("returning the saml response to acs @ ", actionURL)
	// }
	// body, err = ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	t.Error(err)
	// }
	response.Body.Close()
	// fmt.Println(string(body))

	// get resource again
	response, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	response.Body.Close()
	fmt.Println(string(body))

	if !strings.Contains(string(body), TEST_RESOURCE_CONTENT) {
		t.Error("getting resource after login")
	}

}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, TEST_RESOURCE_CONTENT+", %s!",
		samlsp.AttributeFromContext(r.Context(), "displayName"))
}

func extractFormAction(htmlContent string) (string, string, string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", "", "", err
	}

	var actionURL string
	var samlReq string
	var relayState string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "form" {
			for _, attr := range n.Attr {
				if attr.Key == "action" {
					actionURL = attr.Val
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "input" {
			isRelayState := false
			isSAMLReq := false
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == "SAMLRequest" {
					isSAMLReq = true
					continue
				}
				if attr.Key == "name" && attr.Val == "relayState" {
					isRelayState = true
					continue
				}
				if attr.Key == "value" {
					if isRelayState {
						relayState = attr.Val
						fmt.Println("relayState: ", relayState)
						isRelayState = false
						continue
					}
					if isSAMLReq {
						samlReq = attr.Val
						fmt.Println("SAMLReq: ", samlReq)
						isSAMLReq = false
						continue
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return actionURL, samlReq, relayState, err
}

func extextractFormSamlResponse(htmlContent string) (string, string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", "", err
	}

	var samlResponse string
	var actionURL string

	var isSamlResponse bool = false

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "form" {
			for _, attr := range n.Attr {
				if attr.Key == "action" {
					actionURL = attr.Val
					break
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, attr := range n.Attr {

				if attr.Key == "name" && attr.Val == "SAMLResponse" {
					isSamlResponse = true
					continue
				}
				if isSamlResponse {
					if attr.Key == "value" {
						samlResponse = attr.Val
					}
					isSamlResponse = false
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return actionURL, samlResponse, err
}
