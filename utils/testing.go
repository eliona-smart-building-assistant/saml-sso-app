package utils

import (
	"errors"
	"net/url"
	"saml-sso/apiserver"

	"github.com/eliona-smart-building-assistant/go-utils/log"
)

func CreateRandomApiBasicConfig() apiserver.BasicConfiguration {
	var basicCnf apiserver.BasicConfiguration = apiserver.BasicConfiguration{
		Enable: RandomBoolean(),
		OwnUrl: RandomUrl(),
	}

	cert, key, err := CreateSelfsignedX509Certificate(RandomInt(10, 1000), 2048, nil, nil)
	if err != nil {
		log.Error("utils/testing", "create test cert: %v", err)
	}

	basicCnf.ServiceProviderCertificate = cert
	basicCnf.ServiceProviderPrivateKey = key

	if RandomBoolean() {
		randUrl := RandomUrl()
		basicCnf.IdpMetadataUrl = &randUrl
	} else {
		xml := "" // ToDo
		basicCnf.IdpMetadataXml = &xml
	}
	return basicCnf
}

func CreateRandomApiAdvancedConfig() apiserver.AdvancedConfiguration {
	return apiserver.AdvancedConfiguration{
		Enable:                   RandomBoolean(),
		AllowInitializationByIdp: RandomBoolean(),
		SignedRequest:            RandomBoolean(),
		ForceAuthn:               RandomBoolean(),
		EntityId:                 RandomUrl() + "/" + RandomCharacter(5, false),
		CookieSecure:             RandomBoolean(),
		LoginFailedUrl:           RandomUrl() + "/" + RandomCharacter(RandomInt(2, 10), false),
	}
}

func CreateRandomApiAttributeMap() apiserver.AttributeMap {
	mapping := apiserver.AttributeMap{
		Enable: RandomBoolean(),
		Email:  RandomCharacter(RandomInt(1, 20), false) + "@" + RandomCharacter(RandomInt(5, 20), false) + ".net",
	}
	if RandomBoolean() {
		firstName := RandomCharacter(RandomInt(1, 21), false)
		lastName := RandomCharacter(RandomInt(1, 32), true)
		phone := RandomCharacter(RandomInt(5, 12), false)
		mapping.FirstName = &firstName
		mapping.LastName = &lastName
		mapping.Phone = &phone
	}
	return mapping
}

func CreateRandomApiPermissions() apiserver.Permissions {
	perm := apiserver.Permissions{
		Enable:                  RandomBoolean(),
		DefaultSystemRole:       "superadmin",
		DefaultProjRole:         "admin",
		SystemRoleSamlAttribute: nil, // DoTo
		SystemRoleMap:           nil,
		ProjRoleSamlAttribute:   nil,
		ProjRoleMap:             nil,
	}

	if RandomBoolean() {
		perm.DefaultProjRole = "operator"
		perm.DefaultSystemRole = "regular"
	}

	return perm
}

func ValidateUrl(in string) error {

	_, err := url.ParseRequestURI(in)
	if err != nil {
		return errors.New("parse url failed: " + err.Error())
	}

	u, err := url.Parse(in)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return errors.New("uncomplete url. missing scheme or host: " + err.Error())
	}

	return nil
}
