package utils_test

import (
	"saml-sso/utils"
	"testing"
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
