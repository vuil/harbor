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

package lightwave

import (
	"fmt"

	"github.com/vmware/harbor/src/common"
	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/ui/config"
)

// GetSystemLwConf ...
func GetSystemLwConf() (models.LightwaveSetting, error) {
	var err error
	var lw models.LightwaveSetting
	var authMode string

	authMode, err = config.AuthMode()
	if err != nil {
		log.Errorf("can't load auth mode from system, error: %v", err)
		return lw, err
	}

	if authMode != common.LWAuth {
		return lw, fmt.Errorf("system auth_mode isn't lw_auth, please check configuration")
	}

	lwSetting, err := config.LW()

	lw.DomainName = lwSetting.DomainName
	lw.Endpoint = lwSetting.Endpoint
	lw.AdminUser = lwSetting.AdminUser
	lw.AdminPassword = lwSetting.AdminPassword
	lw.IgnoreCertificates = lwSetting.IgnoreCertificates
	lw.VmdirPath = lwSetting.VmdirPath
	lw.Scopes = lwSetting.Scopes

	return lw, nil
}
