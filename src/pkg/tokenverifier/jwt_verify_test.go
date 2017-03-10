///////////////////////////////////////////////////////////////////////
// Copyright (C) 2016 VMware, Inc. All rights reserved.
// -- VMware Confidential
///////////////////////////////////////////////////////////////////////

package tokenverifier

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

var lwCertResponseBody = `[{"certificates":[{"encoded":"-----BEGIN CERTIFICATE-----\nMIIEFzCCAv+gAwIBAgIJAMZe2Jo/p3crMA0GCSqGSIb3DQEBCwUAMEUx\nHzAdBgNVBAMMFkNBLERDPXZzcGhlcmUsREM9bG9jYWwxCzAJBgNVBAYT\nAlVTMRUwEwYDVQQKDAwxMC4yMC43Mi4xOTIwHhcNMTcwMjIxMDAxMjE2\nWhcNMjcwMjE5MDAxMTQ5WjANMQswCQYDVQQGEwJVUzCCAiIwDQYJKoZI\nhvcNAQEBBQADggIPADCCAgoCggIBALOHyV8ZPyjDxlgZJ0T3N1kk7rN8\n6EW39nuOqUJ1YG4AaM8Z55MyHcm3vs1oKnEUYlBlCkAGQ0cs6KUW3zWB\nVywSdTyh3TfqEUUL4VoTuauU3RL5O52PCEi+GPEhsMtbRYyJTTQycd3M\nEXhAy1GBZYwI9qN3QzLCsSrLlguydUy/FXFSaIsVTUcvXFxqgvdz1rfr\nP8G786ex0srsmOFDMDJgESH5E+u+wMIvMWusAtAqTA0PqU7dXR5yWIsm\nRWXbMGia6bDn3Rtq49EQUEuQfkwObNsL/IQAZvqqHci5DI6tJEX3TGM9\ndGQ3hmcCsIk9sobwmhzpKTYC+HxDkK6CLmKBGOOuaSqPhqcSdNOq0cO1\nC7AMxqgl0Igml9maqsyn9kfhWpOhriVIEfLJq4v7PKKvIHehHscO0Ta7\nSdt8ZnBzmwWmGlVn0kE1vAX4kYqUjnxDQiJB4AW0YWP+ekyZXlbHefWb\n/UsBjp538ii1cByHCEqZajJR6dlMCZn8fR5ijWRv/xgJDmZOlBT8P9uT\n5DVPoX7EJZ/nyJr+5vuz1lxy6nJrG2W9Bn4VGQrcXzYkp+B18Iz/IqkP\nRHym2s34kLYv0xwRwAUt5dy6AdmJWv9S3zYRc3q5jRXWwJxl+MbjzA45\nIkC9lbVqj9RgI4r63uUHjQRjmJCHsSpR9VWnmdIFAgMBAAGjQjBAMB0G\nA1UdDgQWBBRD+4qYNY1d0dxC3rKQEi9U9RGGPzAfBgNVHSMEGDAWgBSB\nJBCR3WlBaF8cEIkMiU5Ovi0/RzANBgkqhkiG9w0BAQsFAAOCAQEAAAu7\nnk44sHwqyNY7PGvDWitqNVKUME2MNntn9yzqY5/4hwdatx8AqLlCSZVd\n5DZy8WHgo0K+LwRtpLiwuWelTjqfpkqnlNDVevrxXwXNALzNCPsbX7m7\nXyeZD0ClPB+lX+ub3XgNuZxswGrJUf0AAOE9iIaSVn4Jt+Efk2xWzMzs\n9l25qgHKtpKkX6DxdinttlVnQIlM1ZRdLUnycPPjouChcIfh8WXV0Q5u\nt0DXfZPBxnVeM9KWyjwo+kMZUCytxCm3HTzQ20Xk2NpiQVjZMc1M4wFt\nwb2uj+oznen42L4tXnyfCFpAupCzpq/iIfRv7vNsfXxBYWKFLGO+yuKl\nsw==\n-----END CERTIFICATE-----"},{"encoded":"-----BEGIN CERTIFICATE-----\nMIIDUjCCAjqgAwIBAgIJANahsMS/1Pi4MA0GCSqGSIb3DQEBCwUAMEUx\nHzAdBgNVBAMMFkNBLERDPXZzcGhlcmUsREM9bG9jYWwxCzAJBgNVBAYT\nAlVTMRUwEwYDVQQKDAwxMC4yMC43Mi4xOTIwHhcNMTcwMjIxMDAxMTQ5\nWhcNMjcwMjE5MDAxMTQ5WjBFMR8wHQYDVQQDDBZDQSxEQz12c3BoZXJl\nLERDPWxvY2FsMQswCQYDVQQGEwJVUzEVMBMGA1UECgwMMTAuMjAuNzIu\nMTkyMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyqmPi5pw\nR/b5gpPJxq3qdwNjVFPAk2qUsomDXbBcZ/x3ZjpEGxf/59XzOE41i7KY\n/l6l1oeVuKZA46Y/cuZ5ueH+H7lxT1aIYRqnc79nkFJI08GTJzwlZGM8\nrzGDPhwK6nQsEoa3eXslPzQknibvU3LgMbHVPhyYfpBLKRi9ZfVQVNna\ntFJOAq6jdek1SfPkZRXMrxl2AxJ+W40UwnWg3O3/uF+CAeEsFhC6gMZJ\nCdwLq5KgetDgz4y/S05MTNAHo4xlIr77BYOBnQ5NJ0A2qntJCl4jmIXR\nT4jrqyQL+TbbvlAuuR4B+X+XLxKkJzreR+6crZJ/i1nIabCYOTBVfwID\nAQABo0UwQzAdBgNVHQ4EFgQUgSQQkd1pQWhfHBCJDIlOTr4tP0cwDgYD\nVR0PAQH/BAQDAgEGMBIGA1UdEwEB/wQIMAYBAf8CAQAwDQYJKoZIhvcN\nAQELBQADggEBAIYJkAHwdwlClJUdb7IDPdN4QdFdkpLeTmpB6Znt+gFP\nb4ZOy93e7t+y9kniPxjCTXBNSIlpO3DPcMYRPyG7aw3gi5u1cDa7fW3F\n6CJ8Gwe1PNiVH7qY2mSnq7h8JN16LXkeyJyXVhuq9Q78+uU91GpljOjb\nqBZw0v9gKiX0nGeyJ2Hs2K0941idQ47Wan59RUmD0052cK+ujgVAKbdP\nl2WuJT4JsOtIZeTCcGEBzhCueGqrIRByUlMBfMwpJ0soFo8cIZuT5/wo\n1/XHqZocDZrgjMpzaVY4HE3pZ4YKzK6Oq3RYPdoDSjsbyppsziS06c3n\ncY39Jpk4y2a4TXZNOTs=\n-----END CERTIFICATE-----"}]}]`

func getLWCertBody() string {
	return fmt.Sprintf("%s", lwCertResponseBody)
}

func TestJwttokenVerificationWithInvalidTokenFirstPartFail(t *testing.T) {
	fmt.Println("Testing: invalid token first part, valid cert")

	token, err := getTokenFromFile("tokens/invalid_token_first")
	assert.Nil(t, err)

	tokenVerifier := NewTokenVerifier("https://lwserver/idm/tenant/vsphere.local/certificates/?scope=TENANT")
	_, err = tokenVerifier.Verify(token)
	assert.NotNil(t, err)
	assert.NotEqual(t, err, ErrTokenExpired)
	fmt.Printf("Error: %# +v\n", pretty.Formatter(err))
}

func TestJwttokenVerificationWithInvalidTokenSecondPartFail(t *testing.T) {
	fmt.Println("Testing: invalid token second part, valid cert")

	token, err := getTokenFromFile("tokens/invalid_token_second")
	assert.Nil(t, err)

	tokenVerifier := NewTokenVerifier("https://lwserver/idm/tenant/vsphere.local/certificates/?scope=TENANT")
	_, err = tokenVerifier.Verify(token)
	assert.NotNil(t, err)
	assert.NotEqual(t, err, ErrTokenExpired)
	fmt.Printf("Error: %# +v\n", pretty.Formatter(err))
}

func TestJwttokenVerificationWithInvalidTokenThirdPartFail(t *testing.T) {
	fmt.Println("Testing: invalid token third part, valid cert")

	token, err := getTokenFromFile("tokens/invalid_token_third")
	assert.Nil(t, err)
	tokenVerifier := NewTokenVerifier("https://lwserver/idm/tenant/vsphere.local/certificates/?scope=TENANT")
	_, err = tokenVerifier.Verify(token)
	assert.NotNil(t, err)
	assert.NotEqual(t, err, ErrTokenExpired)
	fmt.Printf("Error: %# +v\n", pretty.Formatter(err))
}

func TestJwttokenVerificationWithCertExpiredTokenFailed(t *testing.T) {
	defer gock.Off()
	gock.New("https://lwserver").
		Get("idm/tenant/vsphere.local/certificates").
		MatchParams(map[string]string{"scope": "TENANT"}).
		Reply(200).
		BodyString(getLWCertBody())

	fmt.Println("Testing: expired token")
	token, err := getTokenFromFile("tokens/valid_expired_token")
	assert.Nil(t, err)

	tokenVerifier := &TokenVerifier{
		lightwavePublicCert: &Certificates{
			URL:       "https://lwserver/idm/tenant/vsphere.local/certificates/?scope=TENANT",
			Transport: nil},
	}

	jwTToken, err := tokenVerifier.Verify(token)

	fmt.Printf("Error: %# +v\n", pretty.Formatter(err))
	fmt.Printf("Token: %# +v\n", pretty.Formatter(jwTToken))
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrTokenExpired)
}

func TestParseTokenInformation(t *testing.T) {
	fmt.Println("Testing: parse token's detail")
	token, err := getTokenFromFile("tokens/valid_expired_token")
	assert.Nil(t, err)

	jwTToken, err := ParseTokenDetails(token)
	assert.Nil(t, err)
	fmt.Printf("Token: %# +v\n", pretty.Formatter(jwTToken))
	assert.Equal(t, jwTToken.Subject, "Administrator@vsphere.local")
}

func getTokenFromFile(path string) (string, error) {
	var tokenByte []byte
	var err error

	// read key content
	tokenByte, err = ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(tokenByte), nil
}
