package gcp

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type gcpConfig struct {
	Project        *string `cty:"project"`
	CredentialFile *string `cty:"credential_file"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"project": {
		Type: schema.TypeString,
	},
	"credential_file": {
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
