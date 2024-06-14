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

package utils

import "strings"

const (
	UTILS_OWN_URL_PLACEHOLDER = "{ownUrl}"
	UTILS_ERROR_PLACEHODER    = "{error}"
)

func SubstituteOwnUrlUrlString(url string, ownUrl string) (urlOrHtml string, isHtml bool) {
	isHtml = false

	if strings.Contains(url, "<html>") {
		isHtml = true
	}

	return strings.ReplaceAll(url, UTILS_OWN_URL_PLACEHOLDER, ownUrl), isHtml
}

func SubstituteError(content string, err []byte) (result string) {
	return strings.ReplaceAll(content, UTILS_ERROR_PLACEHODER, string(err))
}
