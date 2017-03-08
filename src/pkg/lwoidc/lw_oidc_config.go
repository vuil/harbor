///////////////////////////////////////////////////////////////////////
// Copyright (C) 2016 VMware, Inc. All rights reserved.
// -- VMware Confidential
///////////////////////////////////////////////////////////////////////

package lwoidc

import (
	"github.com/vmware/harbor/src/pkg/oidc"
)

// Default Values
const configOIDCDomainNameDefValue string = "vsphere.local"
const configOIDCEndpointDefValue string = "https://localhost"
const configOIDCAdminUserDefValue string = "Administrator"
const configOIDCAdminPasswordDefValue string = "LW@pass1"
const configOIDCIgnoreCertificatesDefValue bool = true
const configOIDCScopesDefValue string = "openid offline_access id_groups at_groups rs_admin_server"
const configOIDCVmdirPathDefValue string = "/vmdir"
const configOIDCLWErrorsDefValue string = "/tmp/LWErrorMessage.json"

// InitOIDCConfig returns a oidc config
func InitOIDCConfig(options oidc.ConfigOptions) *oidc.Config {
	config := &oidc.Config{
		DomainName:         configOIDCDomainNameDefValue,
		Endpoint:           configOIDCEndpointDefValue,
		AdminUser:          configOIDCAdminUserDefValue,
		AdminPassword:      configOIDCAdminPasswordDefValue,
		IgnoreCertificates: configOIDCIgnoreCertificatesDefValue,
		TokenScopes:        configOIDCScopesDefValue,
		VmdirPath:          configOIDCVmdirPathDefValue,
		LWErrors:           configOIDCLWErrorsDefValue,
		RootCAs:            nil,
	}

	if options.DomainName != "" {
		config.DomainName = options.DomainName
	}
	if options.Endpoint != "" {
		config.Endpoint = options.Endpoint
	}
	if options.AdminUser != "" {
		config.AdminUser = options.AdminUser
	}
	if options.AdminPassword != "" {
		config.AdminPassword = options.AdminPassword
	}
	if options.TokenScopes != "" {
		config.TokenScopes = options.TokenScopes
	}
	if options.VmdirPath != "" {
		config.VmdirPath = options.VmdirPath
	}
	if options.LWErrors != "" {
		config.LWErrors = options.LWErrors
	}

	config.IgnoreCertificates = options.IgnoreCertificates
	config.RootCAs = options.RootCAs

	return config
}
