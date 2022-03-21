package gcp

import (
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/schema"
)

type gcpConfig struct {
	Project                   *string `cty:"project"`
	Credentials               *string `cty:"credentials"`
	CredentialFile            *string `cty:"credential_file"`
	ImpersonateServiceAccount *string `cty:"impersonate_service_account"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"project": {
		Type: schema.TypeString,
	},
	"credentials": {
		Type: schema.TypeString,
	},
	"credential_file": {
		Type: schema.TypeString,
	},
	"impersonate_service_account": {
		Type: schema.TypeString,
	},
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
