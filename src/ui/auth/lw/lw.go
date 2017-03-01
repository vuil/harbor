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

	"github.com/vmware/harbor/src/common/utils/log"

	"github.com/vmware/harbor/src/common/dao"
	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/ui/auth"
	"github.com/vmware/harbor/src/ui/config"

	"github.com/vmware/harbor/src/pkg/lwoidc"
	"github.com/vmware/harbor/src/pkg/oidc"
)

// Auth implements Authenticator interface to authenticate against LDAP
type Auth struct{}

// GetLWOptions returns the config options for connecting with lightwave
func GetLWOptions() oidc.ConfigOptions {
	lwSetting, err := config.LW()
	if err != nil {
		log.Error("Error loading Lightwave configuration")
		panic(err)
	}
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

// Login username : user name without the @<lw_domainName>
func Login(client oidc.Client, username string, password string) (oidc.TokenResponse, oidc.ErrorResponse) {
	log.Debug("lw login for ", username)
	lw, err := config.LW()
	if err != nil {
		log.Error("Error loading Lightwave configuration")
		panic(err)
	}
	userName := username + "@" + lw.DomainName
	token, errorResponse := client.GetTokenByPasswordGrant(userName, password)
	if errorResponse != nil {
		log.Info(fmt.Sprintf("GetTokenByPasswordGrant, User: %s  ,ErrorResponse: %d, %s", username,
			errorResponse.GetStatusCode(), errorResponse.GetFullMessage()))
	}
	return token, errorResponse
}

// Authenticate checks user's credential against the LW OIDC REST endpoint
// if the check is successful a dummy record will be inserted into DB,
// so that this user can be associated to other entities in the system.
func (l *Auth) Authenticate(m models.AuthModel) (*models.User, error) {
	user := m.Principal
	pass := m.Password
	u := models.User{}

	client := GetOIDCClient()

	token, errResponse := Login(client, user, pass)
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
	}
	return &u, nil
}

func init() {
	auth.Register("lw_auth", &Auth{})
}
