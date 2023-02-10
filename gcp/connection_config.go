package gcp

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

type gcpConfig struct {
	Project                   *string  `cty:"project"`
	Credentials               *string  `cty:"credentials"`
	CredentialFile            *string  `cty:"credential_file"`
	ImpersonateServiceAccount *string  `cty:"impersonate_service_account"`
	IgnoreErrorCodes          []string `cty:"ignore_error_codes"`
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
	"ignore_error_codes": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
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
