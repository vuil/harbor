///////////////////////////////////////////////////////////////////////
// Copyright (C) 2016 VMware, Inc. All rights reserved.
// -- VMware Confidential
///////////////////////////////////////////////////////////////////////

package tokenverifier

//go:generate mockgen -source=token.go  -destination ./token_mock.go -package tokenverifier

// TokenVerifier is the interface for verifying a lightwave token
type TokenVerifier interface {
	// Verify verifies a token and returns a JWTToken on success
	Verify(token string) (*JWTToken, error)
}

// JWTToken this is for LW Access Token
type JWTToken struct {
	TokenID    string   `json:"jti"`
	Subject    string   `json:"sub"`
	Audience   []string `json:"aud"`
	Groups     []string `json:"groups"`
	Issuer     string   `json:"iss"`
	IssuedAt   int64    `json:"iat"`
	ExpiresAt  int64    `json:"exp"`
	Scope      string   `json:"scope"`
	TokenType  string   `json:"token_type"`
	TokenClass string   `json:"token_class"`
	Tenant     string   `json:"tenant"`
	// It's possible to have more fields depending on how Lightwave defines the token.
	// This covers all the fields we currently have.
}

// JWTRefreshToken this is for LW Refresh Token, it doesn't have 'group' and 'audience' only contains one
type JWTRefreshToken struct {
	TokenID    string `json:"jti"`
	Subject    string `json:"sub"`
	Audience   string `json:"aud"`
	Issuer     string `json:"iss"`
	IssuedAt   int64  `json:"iat"`
	ExpiresAt  int64  `json:"exp"`
	Scope      string `json:"scope"`
	TokenType  string `json:"token_type"`
	TokenClass string `json:"token_class"`
	Tenant     string `json:"tenant"`
	// It's possible to have more fields depending on how Lightwave defines the token.
	// This covers all the fields we currently have.
}
