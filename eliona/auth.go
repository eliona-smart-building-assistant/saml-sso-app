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

package eliona

import (
	"log"
	"net/http"
)

type Authorization struct {
}

func NewAuthorization() *Authorization {
	return &Authorization{}
}

// redirected from the generic sso endpoint sso/auth (handled by app api)
func (a *Authorization) Authorize(w http.ResponseWriter, r *http.Request) {
	log.Printf("ADFS Auth, aufgerufen mit %s", r.Method)

	// login := ""
	// var firstname, lastname, phone interface{}

	// if s.Map.UserName != nil && s.Map.UserName != "" {
	// 	login = samlsp.AttributeFromContext(r.Context(), s.Map.UserName.(string))
	// } else {
	// 	// default without config
	// 	login = samlsp.AttributeFromContext(r.Context(), "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn")
	// }
	// if s.Map.CutMail != nil && s.Map.CutMail.(bool) == true {
	// 	split := strings.Split(login, "@")
	// 	if len(split) >= 1 {
	// 		login = split[0]
	// 	}
	// 	log.Printf("cut, user: %s\r\n", login)
	// }

	// if s.Map.FirstName != nil {
	// 	firstname = samlsp.AttributeFromContext(r.Context(), s.Map.FirstName.(string))
	// }
	// if s.Map.LastName != nil {
	// 	lastname = samlsp.AttributeFromContext(r.Context(), s.Map.LastName.(string))
	// }
	// if s.Map.Email != nil {
	// 	login = samlsp.AttributeFromContext(r.Context(), s.Map.Email.(string))
	// }
	// if s.Map.Phone != nil {
	// 	phone = samlsp.AttributeFromContext(r.Context(), s.Map.Phone.(string))
	// }

	// log.Printf("firstname: %v, lastname: %v, email/login: %v, phone: %v want to login",
	// 	firstname, lastname, login, phone)

	// jwt := s.DB.CheckUser(firstname, lastname, login, phone)

	// log.Printf("Auth: User %s and token %s", login, jwt)
	// if jwt == "" {

	// 	http.Redirect(w, r, s.Config.OwnURL+"/noLogin/", http.StatusFound)
	// } else {
	// 	var myCookie http.Cookie
	// 	myCookie = http.Cookie{
	// 		Name:  "elionaAuthorization",
	// 		Value: jwt,
	// 		Path:  "/"}

	// 	http.SetCookie(w, &myCookie)
	// 	http.Redirect(w, r, s.Config.OwnURL, http.StatusFound)
	// }
}

func (a *Authorization) ADFSAttributeMap(w http.ResponseWriter, r *http.Request) {
	log.Printf("ADFS Attribute Map, aufgerufen mit %s", r.Method)
	// if r.Method == http.MethodGet {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	response := myService.DB.GetAttributeMap()
	// 	payload, _ := json.Marshal(response)
	// 	w.Write(payload)
	// }

	// if r.Method == http.MethodPatch || r.Method == http.MethodPost {
	// 	w.Header().Set("Content-Type", "application/json")

	// 	old := a.DB.GetAttributeMap()
	// 	new := AttributeMap{}

	// 	byteBody, err := ioutil.ReadAll(r.Body)
	// 	if err != nil {
	// 		http.Error(w, "couldn't read request body", http.StatusBadRequest)
	// 	}

	// 	err = json.Unmarshal(byteBody, &new)
	// 	if err != nil {
	// 		http.Error(w, "couldn't parse json", http.StatusBadRequest)
	// 	}

	// 	if new.UserName != nil {
	// 		old.UserName = new.UserName
	// 	}
	// 	if new.FirstName != nil {
	// 		old.FirstName = new.FirstName
	// 	}
	// 	if new.LastName != nil {
	// 		old.LastName = new.LastName
	// 	}
	// 	if new.Email != nil {
	// 		old.Email = new.Email
	// 	}
	// 	if new.Phone != nil {
	// 		old.Phone = new.Phone
	// 	}
	// 	if new.CutMail != nil {
	// 		old.CutMail = new.CutMail
	// 	}

	// 	err = myService.DB.SetAttributeMap(old)
	// 	if err != nil {
	// 		http.Error(w, "couldn't set new config", http.StatusBadRequest)
	// 	}

	// 	response, _ := json.Marshal(old)

	// 	w.Write(response)
	// }
}
