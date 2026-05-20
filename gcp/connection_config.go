package gcp

import (
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type gcpConfig struct {
	Project                   *string  `hcl:"project"`
	Credentials               *string  `hcl:"credentials"`
	ImpersonateAccessToken    *string  `hcl:"impersonate_access_token"`
	ImpersonateServiceAccount *string  `hcl:"impersonate_service_account"`
  	QuotaProject              *string  `hcl:"quota_project,optional"`
	IgnoreErrorMessages       []string `hcl:"ignore_error_messages,optional"`
	IgnoreErrorCodes          []string `hcl:"ignore_error_codes,optional"`
}

func ConfigInstance() interface{} {
	return &gcpConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) gcpConfig {
	if connection == nil {
		return gcpConfig{}
	}
	raw := connection.GetConfig()
	if raw == nil {
		return gcpConfig{}
	}
	config, _ := raw.(gcpConfig)
	return config
}


