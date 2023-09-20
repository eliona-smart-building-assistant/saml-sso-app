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

package utils

import (
	"math/rand"
	"time"
)

func RandomBoolean() bool {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(2)
	return randomInt == 1
}

func RandomUrl() string {
	return "https://" + RandomCharacter(RandomInt(3, 13), false) + ".net"
}

func RandomCharacter(length int, capital bool) string {
	start := int('a')
	stop := int('z')

	bytes := make([]byte, length)

	if capital {
		start = int('A')
		stop = int('Z')
	}

	for i := 0; i < length; i++ {
		bytes[i] = byte(RandomInt(start, stop))
	}

	return string(bytes)
}

func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())

	return min + rand.Intn(max-min+1)
}
