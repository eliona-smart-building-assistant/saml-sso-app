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
