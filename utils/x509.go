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

import (
	"bytes"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"math/rand"
	"net"
	"time"
)

const (
	CERT_ORGANISATION_O       = "IoTECH AG"
	CERT_COUNTY_C             = "CH"
	CERT_PROVINCE_ST          = "Zurich"
	CERT_LOCALITY_L           = "Winterthur"
	CERT_ORGANISATION_UNIT_OU = "eliona"
	CERT_COMMON_NAME_CN       = "eliona-saml-sp"
)

func GetCombinedX509Certificate(certificate string, privateKey string) (tls.Certificate, error) {
	return tls.X509KeyPair([]byte(certificate), []byte(privateKey))
}

func CreateSelfsignedX509Certificate(validityDays int,
	rsaKeyLen int, commonName *string, ipAddress *string) (string, string, error) {

	if rsaKeyLen < 2048 {
		return "", "", errors.New("rsa key length lower than 2048. to insecure.")
	}

	privateKey, err := rsa.GenerateKey(crand.Reader, rsaKeyLen)
	if err != nil {
		return "", "", err
	}

	cN := CERT_COMMON_NAME_CN
	if commonName != nil {
		cN = *commonName
	}

	var ip net.IP
	if ipAddress != nil {
		ip = net.ParseIP(*ipAddress)
	} else {
		ip = net.IPv4(127, 0, 0, 1)
	}

	certificate := &x509.Certificate{
		SerialNumber: big.NewInt(int64(rand.Int())),
		Subject: pkix.Name{
			Organization:       []string{CERT_ORGANISATION_O},
			Country:            []string{CERT_COUNTY_C},
			Province:           []string{CERT_PROVINCE_ST},
			Locality:           []string{CERT_LOCALITY_L},
			OrganizationalUnit: []string{CERT_ORGANISATION_UNIT_OU},
			CommonName:         cN,
		},
		IPAddresses: []net.IP{ip, net.IPv6loopback},

		NotBefore: time.Now(),

		NotAfter: time.Now().Local().Add(time.Duration(validityDays) *
			time.Hour * 24),

		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth},

		KeyUsage: x509.KeyUsageDigitalSignature,
		DNSNames: []string{cN}, // subject alter names (SAN)
	}

	// selfsign certificate
	signCert, err := x509.CreateCertificate(crand.Reader, certificate,
		certificate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	cert, err := PEMEncodeCertificate(signCert)
	if err == nil {
		key, err := PEMEncodeRsaKey(privateKey)
		return cert, key, err
	}

	return cert, "", err
}

func PEMEncodeCertificate(certificate []byte) (string, error) {
	pemForm := new(bytes.Buffer)
	err := pem.Encode(pemForm, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate,
	})
	return pemForm.String(), err
}

func PEMEncodeRsaKey(privateKey *rsa.PrivateKey) (string, error) {
	pemForm := new(bytes.Buffer)
	err := pem.Encode(pemForm, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return pemForm.String(), err
}
