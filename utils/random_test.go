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
	"unicode"
)

func TestRandom_URL(t *testing.T) {
	oldUrl := ""

	for i := 0; i < 5; i++ {
		url := utils.RandomUrl()
		if oldUrl == url {
			t.Error("url is not random")
		}
		err := utils.ValidateUrl(url)
		if err != nil {
			t.Error("random url", err)
		}
		oldUrl = url
	}
}

func TestRandom_Characters(t *testing.T) {
	cap10_1 := utils.RandomCharacter(10, true)
	if !checkCapitalCharacters(cap10_1) {
		t.Error("non capital letters in capital random string")
	}
	if len(cap10_1) != 10 {
		t.Error("len mismatch in random string")
	}
	cap10_2 := utils.RandomCharacter(10, true)
	if !checkCapitalCharacters(cap10_2) {
		t.Error("non capital letters in capital random string")
	}
	if len(cap10_2) != 10 {
		t.Error("len mismatch in random string")
	}
	if cap10_1 == cap10_2 {
		t.Error("random characters not random")
	}
	low32 := utils.RandomCharacter(32, false)
	if !checkLowerCharacters(low32) {
		t.Error("non lower case characters in lower case random string")
	}
	if len(low32) != 32 {
		t.Error("len mismatch in random string")
	}
}

func TestRandom_Bolean(t *testing.T) {
	hasTrue := false
	hasFalse := false

	for i := 0; i < 20; i++ {
		randBool := utils.RandomBoolean()
		if randBool {
			hasTrue = true
		} else {
			hasFalse = true
		}
		if hasFalse && hasTrue {
			break
		}
	}
	if !(hasFalse && hasTrue) {
		t.Error("rand bool not random")
	}
}

func TestRandom_Integer(t *testing.T) {
	int1 := utils.RandomInt(0, 10000)
	int2 := utils.RandomInt(0, 10000)

	if int1 == int2 {
		t.Error("random int not random")
	}
}

func checkCapitalCharacters(input string) bool {
	for _, char := range input {
		if !unicode.IsUpper(char) {
			return false
		}
	}
	return true
}

func checkLowerCharacters(input string) bool {
	for _, char := range input {
		if !unicode.IsLower(char) {
			return false
		}
	}
	return true
}
