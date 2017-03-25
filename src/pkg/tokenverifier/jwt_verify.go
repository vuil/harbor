///////////////////////////////////////////////////////////////////////
// Copyright (C) 2016 VMware, Inc. All rights reserved.
// -- VMware Confidential
///////////////////////////////////////////////////////////////////////

package tokenverifier

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// ErrTokenNotValid Return when token is not valid
	ErrTokenNotValid = errors.New("Error token not valid")
	// ErrCertificateNotFoundForKey Return when the certifcate is not in the response from LW
	ErrCertificateNotFoundForKey = errors.New("Certificate not found for key")
	// ErrTokenExpired token is expired
	ErrTokenExpired = errors.New("Token is expired")
)

// TokenVerifier object to call the Verify
type TokenVerifier struct {
	lightwavePublicCert *Certificates
}

// Cert for parsing json
type Cert struct {
	Cert string `json:"encoded"`
}

// Certs for parsing json
type Certs struct {
	Certs []Cert `json:"certificates"`
}

const (
	cacheHeader           string        = "Cache-Control"
	defaultCertsCacheTime time.Duration = 1 * time.Hour
)

// NewTokenVerifier creates a TokenVerifier with a lightwave certificate endpoint
func NewTokenVerifier(lwCertURL string) *TokenVerifier {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	tv := &TokenVerifier{
		lightwavePublicCert: &Certificates{URL: lwCertURL, Transport: tr},
	}
	return tv
}

// Verify verify the token and also return the parsed result
func (t *TokenVerifier) Verify(token string) (*JWTToken, error) {
	rawToken, err := t.lightwavePublicCert.validateToken(token)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	jwtToken := parseToken(rawToken)

	return jwtToken, err
}

// ParseTokenDetails parse the token. not validate it.
func ParseTokenDetails(token string) (*JWTToken, error) {
	jwtToken := &JWTToken{}
	// take out payload from token
	claimBytes, err := getPayloadFromToken(token)
	// unmarshal the jwtToken
	err = json.Unmarshal(claimBytes, &jwtToken)
	if err != nil {
		log.Println("unmarshal string ", err)
		return nil, err
	}

	return jwtToken, nil
}

// ParseRefreshTokenDetails parse the refresh token
func ParseRefreshTokenDetails(token string) (*JWTRefreshToken, error) {
	jwtRefreshToken := &JWTRefreshToken{}
	// take out payload from token
	claimBytes, err := getPayloadFromToken(token)
	// unmarshal the jwtRefreshToken
	err = json.Unmarshal(claimBytes, &jwtRefreshToken)
	if err != nil {
		log.Println("unmarshal string ", err)
		return nil, err
	}

	return jwtRefreshToken, nil
}

// take out payload bytes from token
func getPayloadFromToken(token string) ([]byte, error) {
	// signature is not json, second part, claims
	chunks := strings.Split(token, ".")
	if len(chunks) < 2 {
		return nil, ErrTokenNotValid
	}

	// only took out claims
	payload := chunks[1]

	claimBytes, err := jwt.DecodeSegment(payload)
	if err != nil {
		log.Println("decode string erorr ", err)
		return nil, err
	}

	return claimBytes, nil
}

// take out information from jwt.token.Claims to Type JWTTOKEN
func parseToken(rawToken *jwt.Token) *JWTToken {
	jwtToken := &JWTToken{}
	claims := rawToken.Claims.(jwt.MapClaims)

	if claims["jti"] != nil {
		jwtToken.TokenID = claims["jti"].(string)
	}
	if claims["sub"] != nil {
		jwtToken.Subject = claims["sub"].(string)
	}
	if claims["aud"] != nil {
		jwtToken.Audience = interfaceArrayToStringArray(claims["aud"].([]interface{}))
	}
	if claims["groups"] != nil {
		jwtToken.Groups = interfaceArrayToStringArray(claims["groups"].([]interface{}))
	}
	if claims["iss"] != nil {
		jwtToken.Issuer = claims["iss"].(string)
	}
	if claims["iat"] != nil {
		jwtToken.IssuedAt = int64(claims["iat"].(float64))
	}
	if claims["exp"] != nil {
		jwtToken.ExpiresAt = int64(claims["exp"].(float64))
	}
	if claims["scope"] != nil {
		jwtToken.Scope = claims["scope"].(string)
	}
	if claims["token_type"] != nil {
		jwtToken.TokenType = claims["token_type"].(string)
	}
	if claims["token_class"] != nil {
		jwtToken.TokenClass = claims["token_class"].(string)
	}
	if claims["tenant"] != nil {
		jwtToken.Tenant = claims["tenant"].(string)
	}

	return jwtToken
}

// transfer interface array to string array
func interfaceArrayToStringArray(interArray []interface{}) []string {
	var array []string
	for _, v := range interArray {
		array = append(array, v.(string))
	}
	return array
}
