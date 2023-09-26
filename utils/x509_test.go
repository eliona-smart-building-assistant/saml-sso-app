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

package utils_test

import (
	"saml-sso/utils"
	"testing"
)

func TestX509_Certificate(t *testing.T) {

	var (
		cn *string
		ip *string
	)

	for i := 0; i < 5; i++ {

		cn = nil
		ip = nil

		if utils.RandomBoolean() {
			commonName := utils.RandomCharacter(utils.RandomInt(1, 131), false)
			cn = &commonName
			ipS := "192.0.2.1"
			ip = &ipS
		}

		cert, key, err := utils.CreateSelfsignedX509Certificate(utils.RandomInt(1, 10000), 2048, cn, ip)
		if err != nil {
			t.Error("generate a self signed certificate: ", err)
		}

		_, err = utils.GetCombinedX509Certificate(cert, key)
		if err != nil {
			t.Error("combine cert to a tls cert, maybe mismatch key <-> cert: ", err)
		}
	}
}
