///////////////////////////////////////////////////////////////////////
// Copyright (C) 2016 VMware, Inc. All rights reserved.
// -- VMware Confidential
///////////////////////////////////////////////////////////////////////

package lwoidc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/pkg/oidc"
)

// LWErrors is the struct of errors
type LWErrors struct {
	CreateUserInvalidGrant   string `json:"create_user_invalid_grant"`
	LoginInvalidGrant        string `json:"login_invalid_grant"`
	AddUserToGroupNoUser     string `json:"add_user_to_group_no_user"`
	AddUserToGroupUserExists string `json:"add_user_to_group_user_exists"`
	CreateUserExists         string `json:"create_user_exists"`
	GetUserNotFound          string `json:"get_user_not_found"`
}

// InitLWErrors returns an initialized LWErrors
func InitLWErrors(lwErrorMsgFile string) (oidc.LWErrors, error) {
	log.Info(fmt.Sprintf("LWErrorMsgMap paths: %s", lwErrorMsgFile))
	file, err := ioutil.ReadFile(lwErrorMsgFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot read the LWErrorMessage file: %s", err))
		return nil, err
	}

	var LWErrorMsg LWErrors
	err = json.Unmarshal(file, &LWErrorMsg)

	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot unmarshal the LWErrorMessage file: %s", err))
		return nil, err
	}
	return &LWErrorMsg, nil
}

// GetCreateUserInvalidGrant /
func (l *LWErrors) GetCreateUserInvalidGrant() string {
	return l.CreateUserInvalidGrant
}

// GetLoginInvalidGrant /
func (l *LWErrors) GetLoginInvalidGrant() string {
	return l.LoginInvalidGrant
}

// GetCreateUserExists /
func (l *LWErrors) GetCreateUserExists() string {
	return l.CreateUserExists
}

// GetFuncUserNotFound /
func (l *LWErrors) GetFuncUserNotFound() string {
	return l.GetUserNotFound
}

// GetAddUserToGroupNoUser /
func (l *LWErrors) GetAddUserToGroupNoUser() string {
	return l.AddUserToGroupNoUser
}

// GetAddUserToGroupUserExists /
func (l *LWErrors) GetAddUserToGroupUserExists() string {
	return l.AddUserToGroupUserExists
}
