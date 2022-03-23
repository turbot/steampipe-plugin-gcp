package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-gcp/gcp"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: gcp.Plugin})
}
