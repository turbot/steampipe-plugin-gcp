package gcp

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type gcpConfig struct {
	Project                   *string  `hcl:"project"`
	Credentials               *string  `hcl:"credentials"`
	CredentialFile            *string  `hcl:"credential_file"`
	ImpersonateServiceAccount *string  `hcl:"impersonate_service_account"`
	IgnoreErrorCodes          []string `hcl:"ignore_error_codes,optional"`
	QuotaProject              *string  `hcl:"quota_project,optional"`
}

func ConfigInstance() interface{} {
	return &gcpConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) gcpConfig {
	if connection == nil || connection.Config == nil {
		return gcpConfig{}
	}
	config, _ := connection.Config.(gcpConfig)
	return config
}
