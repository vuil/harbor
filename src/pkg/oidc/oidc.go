///////////////////////////////////////////////////////////////////////
// Copyright (C) 2016 VMware, Inc. All rights reserved.
// -- VMware Confidential
///////////////////////////////////////////////////////////////////////

package oidc

import (
	"crypto/x509"
	"net/http"
)

//go:generate mockgen -source=oidc.go  -destination ../oidc_test/gen_mock_oidc.go -package oidc_test -imports .=github.com/vmware/harbor/src/pkg/oidc

// Config oidc config struct
type Config struct {
	Endpoint           string
	DomainName         string
	AdminUser          string
	AdminPassword      string
	IgnoreCertificates bool
	TokenScopes        string
	VmdirPath          string
	LWErrors           string
	RootCAs            *x509.CertPool
}

// ConfigOptions options for configuring oidc config struct
type ConfigOptions struct {
	Endpoint           string
	DomainName         string
	AdminUser          string
	AdminPassword      string
	IgnoreCertificates bool
	TokenScopes        string
	VmdirPath          string
	LWErrors           string
	RootCAs            *x509.CertPool
}

// TokenResponse struct
type TokenResponse interface {
	GetExpiresIn() int
	GetAccessToken() string
	GetRefreshToken() string
	GetIDToken() string
	GetTokenType() string
}

// UserResponse struct
type UserResponse interface {
	GetLastName() string
	GetFirstName() string
	GetEmail() string
	GetDescription() string
	GetUpn() string
	GetDomain() string
	GetLocked() bool
	GetDisabled() bool
}

// UserListResponse struct
type UserListResponse interface {
	GetAllUsers() []UserResponse
	GetUserResponse(int) UserResponse
}

// ErrorResponse struct
type ErrorResponse interface {
	GetStatusCode() int
	SetStatusCode(int)
	GetDetails() string
	SetDetails(string)
	GetCause() string
	SetCause(string)
	GetErrorDescription() string
	SetErrorDescription(string)
	GetError() string
	SetError(string)
	GetFullMessage() string
}

// Response struct
type Response interface {
	GetError() error
	GetErrorResponse() ErrorResponse
	GetUserResponse() UserResponse
	GetUserListResponse() UserListResponse
	GetTokenResponse() TokenResponse
	GetHTTPResponse() *http.Response
	GetData() []byte
}

// Client interface
type Client interface {
	GetCerts() Response
	GetTokenByPasswordGrant(string, string) (TokenResponse, ErrorResponse)
	GetTokenByRefreshTokenGrant(string) (TokenResponse, ErrorResponse)
	GetJSONPublicKey() Response
	CreateUser(interface{}, TokenResponse) Response
	GetUser(string, string) Response
	UpdateUser(string, interface{}, TokenResponse) Response
	DeleteUser(string, TokenResponse) Response
	UpdatePassword(string, interface{}, TokenResponse) Response
	AddUserToGroup(string, string, TokenResponse) Response
	ListUser(string) Response
}

// ClientFactory interface
type ClientFactory interface {
	GetClient() Client
	RelClient(Client)
}

// LWErrors interface
type LWErrors interface {
	GetCreateUserInvalidGrant() string
	GetLoginInvalidGrant() string
	GetAddUserToGroupNoUser() string
	GetAddUserToGroupUserExists() string
	GetCreateUserExists() string
	GetFuncUserNotFound() string
}
