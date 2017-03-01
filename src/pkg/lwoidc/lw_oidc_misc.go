///////////////////////////////////////////////////////////////////////
// Copyright (C) 2016 VMware, Inc. All rights reserved.
// -- VMware Confidential
///////////////////////////////////////////////////////////////////////

package lwoidc

import (
	"net/http"

	"fmt"

	"github.com/vmware/harbor/src/pkg/oidc"
)

// Certs
const certDownloadPath string = "/afd/vecs/ssl"

// URLS and paths
const baseVmdirTenant string = "tenant/"
const baseVmdirPathUserSuffix string = "/users"
const baseVmdirPathGroupSuffix string = "/groups"
const urlPathForPublicKey string = "/openidconnect/jwks/"
const urlPathForToken string = "/openidconnect/token"
const urlPathForCertificate string = "/certificates/?scope=TENANT"

// Token request helpers
const passwordGrantFormatString = "grant_type=password&username=%s&password=%s&scope=%s"
const refreshTokenGrantFormatString = "grant_type=refresh_token&refresh_token=%s"
const idmBasePath = "/idm/tenant/"

var vmdirBasePath string

// Cert download helper

// LwCert struct
type LwCert struct {
	Value string `json:"encoded"`
}

// RestUserDetails struct
type RestUserDetails struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Description string `json:"description"`
}

// RestPasswordDetails struct
type RestPasswordDetails struct {
	Password *string `json:"password,omitempty"`
	Lifetime int64   `json:"lifetime"`
	LastSet  int64   `json:"lastSet"`
}

// RestUser struct
type RestUser struct {
	Name            string               `json:"name"`
	Domain          string               `json:"domain"`
	Locked          bool                 `json:"locked"`
	Disabled        bool                 `json:"disabled"`
	Details         *RestUserDetails     `json:"details,omitempty"`
	PasswordDetails *RestPasswordDetails `json:"passwordDetails,omitempty"`
}

// RestPassword struct
type RestPassword struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword,omitempty"`
}

// TokenResponse LW OIDC Token Response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
}

// UserResponseDetails response deails
type UserResponseDetails struct {
	LastName    string `json:"lastName"`
	FirstName   string `json:"firstName"`
	Email       string `json:"email"`
	Description string `json:"description"`
	Upn         string `json:"upn"`
}

// UserResponse is struct for LW OIDC User Response
type UserResponse struct {
	Domain   string              `json:"domain"`
	Locked   bool                `json:"locked"`
	Disabled bool                `json:"disabled"`
	Details  UserResponseDetails `json:"details,omitempty"`
}

// UserListResponse is struct of list of users
type UserListResponse struct {
	Users []*UserResponse `json:"users"`
}

// ErrorResponse LW OIDC Response
type ErrorResponse struct {
	StatusCode       int
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Details          string `json:"details"`
	Cause            string `json:"cause"`
}

// Response is returned by most of the OIDC calls
type Response struct {
	Error            error
	ErrorResponse    oidc.ErrorResponse
	UserResponse     oidc.UserResponse
	UserListResponse oidc.UserListResponse
	TokenResponse    oidc.TokenResponse
	HTTPResponse     *http.Response
	Data             []byte
}

// InitVmdirBasePath (Base dir support for older versions of LW)
func InitVmdirBasePath(oidcConfig *oidc.Config) {
	vmdirBasePath = fmt.Sprintf("%s/%s", oidcConfig.VmdirPath, baseVmdirTenant)
}

// GetVmdirBasePath returns the basepath for vmdir
func GetVmdirBasePath() string {
	return vmdirBasePath
}

// GetIdmBasePath returns the basepath for idm
func GetIdmBasePath() string {
	return idmBasePath
}

// Lw OIDC Token Response

// GetExpiresIn returns ExpiresIn
func (r *TokenResponse) GetExpiresIn() int {
	return r.ExpiresIn
}

// GetAccessToken returns AccessToken
func (r *TokenResponse) GetAccessToken() string {
	return r.AccessToken
}

// GetRefreshToken returns RefreshToken
func (r *TokenResponse) GetRefreshToken() string {
	return r.RefreshToken
}

// GetIDToken returns IDToken
func (r *TokenResponse) GetIDToken() string {
	return r.IDToken
}

// GetTokenType returns TokenType
func (r *TokenResponse) GetTokenType() string {
	return r.TokenType
}

// Lw OIDC User Response

// GetLastName returns LastName
func (r *UserResponse) GetLastName() string {
	return r.Details.LastName
}

// GetFirstName returns FirstName
func (r *UserResponse) GetFirstName() string {
	return r.Details.FirstName
}

// GetEmail returns Email
func (r *UserResponse) GetEmail() string {
	return r.Details.Email
}

// GetDescription returns Description
func (r *UserResponse) GetDescription() string {
	return r.Details.Description
}

// GetUpn returns Upn
func (r *UserResponse) GetUpn() string {
	return r.Details.Upn
}

// GetDomain returns Domain
func (r *UserResponse) GetDomain() string {
	return r.Domain
}

// GetLocked returns Locked
func (r *UserResponse) GetLocked() bool {
	return r.Locked
}

// GetDisabled returns Disabled
func (r *UserResponse) GetDisabled() bool {
	return r.Disabled
}

// Lw OIDC UserListResponse

// GetAllUsers returns an array of user responses
func (r *UserListResponse) GetAllUsers() []oidc.UserResponse {
	array := []oidc.UserResponse{}
	for i := range r.Users {
		array = append(array, r.GetUserResponse(i))
	}
	return array
}

// GetUserResponse returns a user response
func (r *UserListResponse) GetUserResponse(index int) oidc.UserResponse {
	if index < 0 || index >= len(r.Users) {
		return nil
	}
	return r.Users[index]
}

// LW OIDC Error Response

// NewErrorResponse returns an ErrorResponse
func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Error:      err.Error(),
	}
}

// NewErrorResponseWithStatusCode returns an ErrorResponse
func NewErrorResponseWithStatusCode(err error, statusCode int) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: statusCode,
		Error:      err.Error(),
	}
}

// GetStatusCode gets the error response's status code
func (e *ErrorResponse) GetStatusCode() int {
	return e.StatusCode
}

// SetStatusCode sets the error response's status code
func (e *ErrorResponse) SetStatusCode(statusCode int) {
	e.StatusCode = statusCode
}

// GetError gets the error response's error string
func (e *ErrorResponse) GetError() string {
	return e.Error
}

// SetError sets the error response's error string
func (e *ErrorResponse) SetError(error string) {
	e.Error = error
}

// GetErrorDescription gets the error response's error description
func (e *ErrorResponse) GetErrorDescription() string {
	return e.ErrorDescription
}

// SetErrorDescription sets the error response's error description
func (e *ErrorResponse) SetErrorDescription(description string) {
	e.ErrorDescription = description
}

// GetDetails gets the error response's details
func (e *ErrorResponse) GetDetails() string {
	return e.Details
}

// SetDetails sets the error response's details
func (e *ErrorResponse) SetDetails(details string) {
	e.Details = details
}

// GetCause gets the error response's cause
func (e *ErrorResponse) GetCause() string {
	return e.Cause
}

// SetCause sets the error response's cause
func (e *ErrorResponse) SetCause(cause string) {
	e.Cause = cause
}

// GetFullMessage returns the error response's full message
func (e *ErrorResponse) GetFullMessage() string {
	fullMessage := ""
	var addSep = false

	if e.Error != "" {
		fullMessage += e.Error
		addSep = true
	}

	if e.ErrorDescription != "" {
		if addSep {
			fullMessage += "; "
		}
		fullMessage += e.ErrorDescription
		addSep = true
	}

	if e.Details != "" {
		if addSep {
			fullMessage += "; "
		}
		fullMessage += e.Details
		addSep = true
	}

	if e.Cause != "" {
		if addSep {
			fullMessage += "; "
		}
		fullMessage += e.Cause
		addSep = true
	}

	return fullMessage
}

// OIDC response

// GetError returns the response's error
func (r *Response) GetError() error {
	return r.Error
}

// GetErrorResponse returns the response's error response
func (r *Response) GetErrorResponse() oidc.ErrorResponse {
	return r.ErrorResponse
}

// GetUserResponse returns the response's user respose
func (r *Response) GetUserResponse() oidc.UserResponse {
	return r.UserResponse
}

// GetUserListResponse returns the response's user list response
func (r *Response) GetUserListResponse() oidc.UserListResponse {
	return r.UserListResponse
}

// GetTokenResponse returns the response's token response
func (r *Response) GetTokenResponse() oidc.TokenResponse {
	return r.TokenResponse
}

// GetHTTPResponse returns the response's data
func (r *Response) GetHTTPResponse() *http.Response {
	return r.HTTPResponse
}

// GetData returns the response's data
func (r *Response) GetData() []byte {
	return r.Data
}
