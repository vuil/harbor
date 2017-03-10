/*
   Copyright (c) 2016 VMware, Inc. All Rights Reserved.
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package lw

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vmware/harbor/src/common/utils/log"

	"github.com/vmware/harbor/src/common/dao"
	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/ui/auth"
	"github.com/vmware/harbor/src/ui/config"

	"github.com/vmware/harbor/src/pkg/lwoidc"
	"github.com/vmware/harbor/src/pkg/oidc"
	"github.com/vmware/harbor/src/pkg/tokenverifier"
)

// Auth implements Authenticator interface to authenticate against LDAP
type Auth struct{}

const (
	// GroupPrefix is the prefix to all LW group names that harbor recognizes
	GroupPrefix string = "Harbor_Project"
)

// GroupRole a user's role in a group
type GroupRole struct {
	groupName string
	role      string
}

// GetLWOptions returns the config options for connecting with lightwave
func GetLWOptions() oidc.ConfigOptions {
	lwSetting := config.LW()
	return oidc.ConfigOptions{
		DomainName:         lwSetting.DomainName,
		Endpoint:           lwSetting.Endpoint,
		AdminUser:          lwSetting.AdminUser,
		AdminPassword:      lwSetting.AdminPassword,
		IgnoreCertificates: lwSetting.IgnoreCertificates,
		TokenScopes:        lwSetting.Scopes,
		VmdirPath:          lwSetting.VmdirPath,
		LWErrors:           "",
		RootCAs:            nil,
	}
}

// GetOIDCClient returns the oidc client used to connect to lightwave
func GetOIDCClient() oidc.Client {
	options := GetLWOptions()
	clientFactory := lwoidc.NewClientFactory(options)
	client := clientFactory.GetClient()
	return client
}

func getGroupRoles(jwtToken *tokenverifier.JWTToken) (result []GroupRole) {
	for _, lwGroupName := range jwtToken.Groups {
		parts := strings.Split(lwGroupName, "\\")
		groupName := parts[len(parts)-1]
		if strings.HasPrefix(groupName, GroupPrefix) {
			groupRoleString := groupName[len(GroupPrefix):]
			parts2 := strings.Split(groupRoleString, "_")
			result = append(result, GroupRole{parts2[0], parts2[1]})
		}
	}

	return
}

// getProjectRolesFromToken returns the map of a user's projects to user's role in them
func getUserProjectsFromToken(lwSetting config.LightwaveSetting, tokenVerifier tokenverifier.TokenVerifier, token oidc.TokenResponse) (groupRoles []GroupRole) {
	time.Sleep(1 * time.Second)
	jwtToken, error := tokenVerifier.Verify(token.GetAccessToken())
	log.Debugf("========== JWT Error = %+v\n Token = %+v\n", error, jwtToken)
	if error == nil {
		groupRoles = getGroupRoles(jwtToken)
		log.Debugf("========== GROUP roles = %+v\n\n", groupRoles)
	}

	return
}

// ensureUserProjects ensures that projects associated with a user's token is
// created. Creating of the projects is sometimes necessary if the user
// authenticated via Lightwave is already preassigned memberships to lightwave
// groups that map to Harbor projects
func ensureUserProjects(lwSetting config.LightwaveSetting, user models.User, token oidc.TokenResponse) {
	certURL := tokenverifier.GetLightwaveCertURL(lwSetting.Endpoint, lwSetting.DomainName)
	tv := tokenverifier.NewTokenVerifier(certURL)
	getUserProjectsFromToken(lwSetting, tv, token)

	// TODO create projects if necessary
}

// Login username : user name without the @<lw_domainName>
func Login(client oidc.Client, lwSetting config.LightwaveSetting, username string, password string) (oidc.TokenResponse, oidc.ErrorResponse) {
	log.Debug("lw login for ", username)

	domainName := lwSetting.DomainName

	userName := username + "@" + domainName
	token, errorResponse := client.GetTokenByPasswordGrant(userName, password)
	if errorResponse != nil {
		log.Info(fmt.Sprintf("GetTokenByPasswordGrant, User: %s, ErrorResponse: %d, %s", username,
			errorResponse.GetStatusCode(), errorResponse.GetFullMessage()))
	}
	return token, errorResponse
}

// Authenticate checks user's credentials against the LW OIDC REST endpoint
// if the check is successful a dummy record will be inserted into DB,
// so that this user can be associated to other entities in the system.
func (l *Auth) Authenticate(m models.AuthModel) (*models.User, error) {
	user := m.Principal
	pass := m.Password
	u := models.User{}

	client := GetOIDCClient()
	lwSetting := config.LW()

	token, errResponse := Login(client, lwSetting, user, pass)
	if errResponse != nil {
		return nil, errors.New(errResponse.GetFullMessage())
	}
	log.Debugf("Authenticated with token = %v\n", token)

	u.Username = m.Principal
	log.Debug("username:", u.Username)
	exist, err := dao.UserExists(u, "username")
	if err != nil {
		return nil, err
	}

	if exist {
		currentUser, err := dao.GetUser(u)
		if err != nil {
			return nil, err
		}
		u.UserID = currentUser.UserID
		ensureUserProjects(lwSetting, u, token)
	} else {
		u.Realname = m.Principal
		u.Password = "12345678Lw"
		u.Comment = "registered from Lightwave."
		if u.Email == "" {
			u.Email = u.Username + "@placeholder.com"
		}
		userID, err := dao.Register(u)
		if err != nil {
			return nil, err
		}
		u.UserID = int(userID)
		ensureUserProjects(lwSetting, u, token)
	}
	return &u, nil
}

func init() {
	auth.Register("lw_auth", &Auth{})
}
