///////////////////////////////////////////////////////////////////////
// Copyright (C) 2016 VMware, Inc. All rights reserved.
// -- VMware Confidential
///////////////////////////////////////////////////////////////////////

package lwoidc

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mendsley/gojwk"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/pkg/oidc"
)

// Client struct
type Client struct {
	config     *oidc.Config
	httpClient *http.Client
	endPoint   string
	domainName string
}

// GetHTTPClient returns the http Client used for making HTTP requests
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

func getVmdirUserPath(c *oidc.Config) string {
	return GetVmdirBasePath() + c.DomainName + baseVmdirPathUserSuffix
}

func getVmdirGroupPath(c *oidc.Config) string {
	return GetVmdirBasePath() + c.DomainName + baseVmdirPathGroupSuffix
}

func getIdmCertificatePath(c *oidc.Config) string {
	return GetIdmBasePath() + c.DomainName + urlPathForCertificate
}

// StaticCheckResponse checks responses from the lightwave server (Accessed from Mocked Tests)
func StaticCheckResponse(response *http.Response, inErr error) oidc.ErrorResponse {
	if inErr != nil {
		log.Info(fmt.Sprintf("StaticCheckResponse, error: %s", inErr.Error()))
		return NewErrorResponseWithStatusCode(inErr, response.StatusCode)
	}

	if response.StatusCode/100 == 2 {
		return nil
	}

	log.Error(fmt.Sprintf("StaticCheckResponse, status code %d", response.StatusCode))

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Info(fmt.Sprintf("StaticCheckResponse, cannot read body response, %s", err.Error()))
		return NewErrorResponseWithStatusCode(err, response.StatusCode)
	}

	var errorResponse ErrorResponse
	err = json.Unmarshal(respBody, &errorResponse)

	log.Info(fmt.Sprintf("StaticCheckResponse, error response body: %s", respBody))

	if err != nil {
		log.Info(fmt.Sprintf("StaticCheckResponse, cannot unmarshal body response, %s", err.Error()))
		return NewErrorResponseWithStatusCode(err, response.StatusCode)
	}

	errorResponse.SetStatusCode(response.StatusCode)
	log.Debug(fmt.Sprintf("StaticCheckResponse, errorResponse: %d, %s",
		errorResponse.StatusCode, errorResponse.GetFullMessage()))
	return &errorResponse
}

// NewClient returns a new lightwave OIDC Client
func NewClient(options oidc.ConfigOptions) *Client {
	config := InitOIDCConfig(options)
	InitVmdirBasePath(config)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.IgnoreCertificates,
			RootCAs:            config.RootCAs},
	}

	c := &Client{
		httpClient: &http.Client{Transport: tr},
		config:     config,
		endPoint:   strings.TrimRight(config.Endpoint, "/"),
		domainName: config.DomainName,
	}

	return c
}

// GetCerts gets certificates from lightwave
func (c *Client) GetCerts() oidc.Response {
	log.Debug(fmt.Sprintf("Entering GetCerts, url: %s",
		c.buildBaseURL(c.GetIdmCertificatePath())))
	// turn TLS verification off for
	url := c.buildBaseURL(c.GetIdmCertificatePath())

	// get the certs
	request, err := http.NewRequest(
		"GET",
		url,
		nil)
	if err != nil {
		log.Info(fmt.Sprintf("Exiting GetCerts, failed in new Request, %s",
			err.Error()))
		return &Response{Error: err}
	}

	response, error := c.httpClient.Do(request)
	if error != nil {
		return &Response{Error: error}
	}

	defer response.Body.Close()

	errorResponse := c.checkResponse(response, error)
	if errorResponse != nil {
		log.Info(fmt.Sprintf("Exiting GetCerts, failed in sending GET request, %s",
			errorResponse.GetError()))

		return &Response{HTTPResponse: response,
			ErrorResponse: errorResponse}
	}
	// read key content
	certificateByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Info(fmt.Sprintf("Exiting GetCerts, failed in reading key content from LW, %s",
			err.Error()))
		return &Response{HTTPResponse: response,
			Error: err}
	}

	return &Response{HTTPResponse: response, Data: certificateByte}
}

// GetTokenByPasswordGrant gets oidc access token via user:pass
func (c *Client) GetTokenByPasswordGrant(username string, password string) (oidc.TokenResponse, oidc.ErrorResponse) {
	log.Info("Entering GetTokenByPasswordGrant")
	body := fmt.Sprintf(passwordGrantFormatString, username, password, c.config.TokenScopes)
	return c.getToken(body)
}

// GetTokenByRefreshTokenGrant gets oidc access token via the frefresh token
func (c *Client) GetTokenByRefreshTokenGrant(refreshToken string) (oidc.TokenResponse, oidc.ErrorResponse) {
	log.Info("Entering GetTokenByRefreshTokenGrant")
	body := fmt.Sprintf(refreshTokenGrantFormatString, refreshToken)
	return c.getToken(body)
}

// GetJSONPublicKey get public key from LW server
func (c *Client) GetJSONPublicKey() oidc.Response {
	log.Debug(fmt.Sprintf("Entering GetJSONPublicKey, url: %s",
		c.buildBaseURL(urlPathForPublicKey+c.domainName)))

	url := c.buildBaseURL(urlPathForPublicKey + c.domainName)

	request, err := http.NewRequest(
		"GET",
		url,
		nil)
	if err != nil {
		log.Info(fmt.Sprintf("Exiting GetJSONPublicKey, failed in new Request, %s",
			err.Error()))
		return &Response{Error: err}
	}

	response, error := c.httpClient.Do(request)
	if error != nil {
		return &Response{Error: error}
	}

	defer response.Body.Close()

	errorResponse := c.checkResponse(response, error)

	if errorResponse != nil {
		log.Info(fmt.Sprintf("Exiting GetJSONPublicKey, failed in sending GET request, %s",
			errorResponse.GetError()))
		return &Response{HTTPResponse: response,
			ErrorResponse: errorResponse}
	}
	// read key content
	publicKeyByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Info(fmt.Sprintf("Exiting GetJSONPublicKey, failed in reading key content from LW, %s",
			err.Error()))
		return &Response{HTTPResponse: response,
			Error: err}
	}

	return &Response{HTTPResponse: response, Data: publicKeyByte}
}

// GetUser gets user information from lightwave
func (c *Client) GetUser(userID string, authorization string) oidc.Response {
	log.Debug(fmt.Sprintf("Getting OIDC user: %# +vs", userID))

	request, err := http.NewRequest("GET", c.buildGetUserURL(userID), nil)

	if err != nil {
		return &Response{Error: err}
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", authorization)

	response, error := c.httpClient.Do(request)
	if error != nil {
		return &Response{Error: error}
	}

	defer response.Body.Close()

	errorResponse := c.checkResponse(response, error)

	if errorResponse != nil {
		return &Response{HTTPResponse: response, ErrorResponse: errorResponse}
	}

	lwUser := &UserResponse{}
	err = json.NewDecoder(response.Body).Decode(lwUser)

	if err != nil {
		return &Response{HTTPResponse: response, Error: err}
	}

	return &Response{HTTPResponse: response, UserResponse: lwUser}
}

// CreateUser creates a user in lightwave
func (c *Client) CreateUser(lwUser interface{}, token oidc.TokenResponse) oidc.Response {
	url := c.buildPostUserURL()

	log.Debug(fmt.Sprintf("CreateUser, URL: %s", url))

	response := c.performRequest("POST", url, lwUser, token)

	if response.GetError() != nil {
		e := response.GetError()
		log.Info(fmt.Sprintf("CreateUser error: %s", e))
	}

	if response.GetErrorResponse() != nil {
		e := response.GetErrorResponse()
		log.Info(fmt.Sprintf("CreateUser ErrorResponse: %d, %s",
			e.GetStatusCode(), e.GetFullMessage()))
	}

	return response
}

// UpdateUser updates information about a user in lightwave
func (c *Client) UpdateUser(userID string, lwUser interface{}, token oidc.TokenResponse) oidc.Response {
	url := c.buildPutUserURL(userID)

	log.Debug(fmt.Sprintf("UpdateUser, URL: %s", url))

	response := c.performRequest("PUT", url, lwUser, token)

	if response.GetError() != nil {
		e := response.GetError()
		log.Info(fmt.Sprintf("UpdateUser error: %s", e))
	}

	if response.GetErrorResponse() != nil {
		e := response.GetErrorResponse()
		log.Info(fmt.Sprintf("UpdateUser ErrorResponse: %d, %s",
			e.GetStatusCode(), e.GetFullMessage()))
	}

	return response
}

// DeleteUser delets a user from lightwave
func (c *Client) DeleteUser(userID string, token oidc.TokenResponse) oidc.Response {
	url := c.buildDeleteUserURL(userID)

	log.Debug(fmt.Sprintf("DeleteUser, URL: %s", url))

	response := c.performRequest("DELETE", url, nil, token)

	if response.GetError() != nil {
		e := response.GetError()
		log.Info(fmt.Sprintf("DeleteUser error: %s", e))
	}

	if response.GetErrorResponse() != nil {
		e := response.GetErrorResponse()
		log.Info(fmt.Sprintf("DeleteUser ErrorResponse: %d, %s",
			e.GetStatusCode(), e.GetFullMessage()))
	}

	return response
}

// UpdatePassword updates the password of a user
func (c *Client) UpdatePassword(userID string, lwPassword interface{}, token oidc.TokenResponse) oidc.Response {
	url := c.buildPutPwURL(userID)

	log.Debug(fmt.Sprintf("UpdatePassword, URL: %s", url))

	response := c.performRequest("PUT", url, lwPassword, token)

	if response.GetError() != nil {
		e := response.GetError()
		log.Info(fmt.Sprintf("UpdatePassword error: %s", e))
	}

	if response.GetErrorResponse() != nil {
		e := response.GetErrorResponse()
		log.Info(fmt.Sprintf("UpdatePassword ErrorResponse: %d, %s",
			e.GetStatusCode(), e.GetFullMessage()))
	}

	return response
}

// AddUserToGroup adds a user to a group
func (c *Client) AddUserToGroup(
	groupName string, userID string,
	token oidc.TokenResponse) oidc.Response {

	url := c.buildAddUserToGroupURL(c.domainName, groupName, userID)

	log.Debug(fmt.Sprintf("AddUserToGroup, URL: %s", url))

	response := c.performRequest("PUT", url, nil, token)

	if response.GetError() != nil {
		e := response.GetError()
		log.Info(fmt.Sprintf("AddUserToGroup error: %s", e))
	}

	if response.GetErrorResponse() != nil {
		e := response.GetErrorResponse()
		log.Info(fmt.Sprintf("AddUserToGroup OidcErrorResponse: %d, %s",
			e.GetStatusCode(), e.GetFullMessage()))
	}

	return response
}

// CreateGroup creates a new lightwave group
func (c *Client) CreateGroup(
	groupName string, description string,
	token oidc.TokenResponse) oidc.Response {

	lwGroup := &RestGroup{
		Name:   groupName,
		Domain: c.domainName,
		Details: &RestGroupDetails{
			Description: description,
		},
	}

	url := c.buildPostGroupURL()

	response := c.performRequest("POST", url, lwGroup, token)

	if response.GetError() != nil {
		e := response.GetError()
		log.Info(fmt.Sprintf("CreateGroup error: %s", e))
	}

	if response.GetErrorResponse() != nil {
		e := response.GetErrorResponse()
		log.Info(fmt.Sprintf("CreateGroup OidcErrorResponse: %d, %s",
			e.GetStatusCode(), e.GetFullMessage()))
	}

	return response
}

// ListUser enumerate all users
func (c *Client) ListUser(authorization string) oidc.Response {
	log.Debug(fmt.Sprintf("List all OIDC user"))

	request, err := http.NewRequest("GET", c.buildListAllUserURL(), nil)

	if err != nil {
		return &Response{Error: err}
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", authorization)

	response, error := c.httpClient.Do(request)
	if error != nil {
		return &Response{Error: error}
	}

	defer response.Body.Close()

	errorResponse := c.checkResponse(response, error)

	if errorResponse != nil {
		return &Response{HTTPResponse: response, ErrorResponse: errorResponse}
	}

	lwUser := &UserListResponse{}
	err = json.NewDecoder(response.Body).Decode(lwUser)

	if err != nil {
		return &Response{HTTPResponse: response, Error: err}
	}

	return &Response{HTTPResponse: response, UserListResponse: lwUser}
}

func (c *Client) getToken(body string) (oidc.TokenResponse, oidc.ErrorResponse) {
	log.Info(fmt.Sprintf("Entering GetToken, url: %s", c.buildBaseURL(urlPathForToken)))
	request, err := http.NewRequest("POST", c.buildBaseURL(urlPathForToken), strings.NewReader(body))
	if err != nil {
		return nil, NewErrorResponse(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(request)

	if err != nil {
		return nil, NewErrorResponse(err)
	}

	defer resp.Body.Close()

	errorResponse := c.checkResponse(resp, err)
	if errorResponse != nil {
		log.Info(fmt.Sprintf("Exiting GetToken, errorResponse: %d, %s",
			errorResponse.GetStatusCode(), errorResponse.GetFullMessage()))
		return nil, errorResponse
	}

	tokenResponse := &TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(tokenResponse)
	if err != nil {
		log.Info(fmt.Sprintf("Exiting GetToken, failed decoding: %s", err.Error()))
		return nil, NewErrorResponse(err)
	}

	return tokenResponse, nil
}

func (c *Client) parseBytePublicKey(publicKeyByte []byte) (*rsa.PublicKey, oidc.ErrorResponse) {
	log.Debug(fmt.Sprintf("Entering parseBytePublicKey, pubicKeyJson: %s", string(publicKeyByte)))

	// parse the json public key
	jwkJSON, err := gojwk.Unmarshal(publicKeyByte)
	if err != nil {
		log.Info(fmt.Sprintf("Exiting parseBytePublicKey, failed in unmarhsal the json key, %s",
			err.Error()))
		return nil, NewErrorResponse(err)
	}

	// store keys
	keys := make([]crypto.PublicKey, len(jwkJSON.Keys))
	// decode public key
	for ii, jwk := range jwkJSON.Keys {
		keys[ii], err = jwk.DecodePublicKey()
		if err != nil {
			log.Info(fmt.Sprintf("Exiting parseBytePublicKey, failed in decode public key, %s",
				err.Error()))
			return nil, NewErrorResponse(err)
		}
	}

	var publicKey *rsa.PublicKey
	var ok bool
	var rsaJSONRawKey = keys[0]

	// cast key to rsa key, use the rsa key by default
	if publicKey, ok = rsaJSONRawKey.(*rsa.PublicKey); !ok {
		log.Info(fmt.Sprintf("Exiting parseBytePublicKey, failed in castting key to rsa.PublicKey, %s",
			err.Error()))
		return nil, NewErrorResponse(err)
	}

	return publicKey, nil
}

func (c *Client) checkResponse(response *http.Response, inErr error) oidc.ErrorResponse {
	return StaticCheckResponse(response, inErr)
}

func (c *Client) performRequest(method string, url string, lwData interface{},
	token oidc.TokenResponse) *Response {
	log.Info(fmt.Sprintf("performRequest URL: %s", url))

	var reader io.Reader
	if lwData != nil {
		content, err := json.Marshal(lwData)
		if err != nil {
			return &Response{Error: err}
		}
		log.Info(fmt.Sprintf("performRequest with content"))
		reader = bytes.NewReader(content)
	} else {
		reader = nil
	}

	// new post request. add form
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return &Response{Error: err}
	}

	request.Header.Add("Content-Type", "application/json")
	adminToken := "Bearer " + token.GetAccessToken()
	request.Header.Add("Authorization", adminToken)

	resp, error := c.httpClient.Do(request)
	if error != nil {
		log.Error(fmt.Sprintf("performRequest, error: %v", error))
		return &Response{Error: error}
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return &Response{HTTPResponse: resp}
	}

	errorResponse := c.checkResponse(resp, error)

	if errorResponse != nil {
		return &Response{HTTPResponse: resp, ErrorResponse: errorResponse}
	}

	// decode user info from response
	userResponse := &UserResponse{}
	err = json.NewDecoder(resp.Body).Decode(userResponse)
	if err != nil {
		log.Error(fmt.Sprintf("performRequest, decode error: %v", err))
		return &Response{HTTPResponse: resp, Error: err}
	}

	return &Response{HTTPResponse: resp, UserResponse: userResponse}
}

func (c *Client) setTransport(tr http.RoundTripper) {
	c.httpClient.Transport = tr
}

// GetVmdirUserPath get the vmdir user path
func (c *Client) GetVmdirUserPath() string {
	return getVmdirUserPath(c.config)
}

// GetVmdirGroupPath gets the vmdir group path
func (c *Client) GetVmdirGroupPath() string {
	return getVmdirGroupPath(c.config)
}

// GetIdmCertificatePath gets the certificate path
func (c *Client) GetIdmCertificatePath() string {
	return getIdmCertificatePath(c.config)
}

func (c *Client) buildBaseURL(path string) string {
	return fmt.Sprintf("%s%s", c.endPoint, path)
}

func (c *Client) buildPostUserURL() string {
	return c.buildBaseURL(c.GetVmdirUserPath())
}

func (c *Client) buildPostGroupURL() string {
	return c.buildBaseURL(c.GetVmdirGroupPath())
}

func (c *Client) buildDeleteUserURL(userID string) string {
	return c.buildBaseUserURL(userID)
}

func (c *Client) buildPutUserURL(userID string) string {
	return c.buildBaseUserURL(userID)
}

func (c *Client) buildGetUserURL(userID string) string {
	return c.buildBaseUserURL(userID)
}

func (c *Client) buildBaseUserURL(userID string) string {
	return fmt.Sprintf("%s/%s@%s",
		c.buildBaseURL(c.GetVmdirUserPath()), userID, c.config.DomainName)
}

func (c *Client) buildPutPwURL(userID string) string {
	return fmt.Sprintf("%s/%s@%s/password",
		c.buildBaseURL(c.GetVmdirUserPath()), userID, c.config.DomainName)
}

func (c *Client) buildAddUserToGroupURL(tenant string, groupName string, userName string) string {
	return c.endPoint + c.GetVmdirGroupPath() +
		fmt.Sprintf("/%s@%s/members?type=user&members=%s@%s", groupName, tenant, userName, tenant)
}

func (c *Client) buildListAllUserURL() string {
	return c.endPoint + fmt.Sprintf("/idm/tenant/%s/search?domain=%s&type=USER", c.domainName, c.domainName)
}
