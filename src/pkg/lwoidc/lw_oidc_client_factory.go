///////////////////////////////////////////////////////////////////////
// Copyright (C) 2016 VMware, Inc. All rights reserved.
// -- VMware Confidential
///////////////////////////////////////////////////////////////////////

package lwoidc

import (
	"github.com/vmware/harbor/src/pkg/oidc"
)

// ClientFactory is factory for oidc clients
type ClientFactory struct {
	options oidc.ConfigOptions
}

// NewClientFactory creates a ClientFactory
func NewClientFactory(
	options oidc.ConfigOptions) *ClientFactory {
	return &ClientFactory{
		options: options,
	}
}

// GetClient returns an oidc client
func (f *ClientFactory) GetClient() oidc.Client {
	return NewClient(f.options)
}

// RelClient release a client
func (f *ClientFactory) RelClient(c oidc.Client) {
}
