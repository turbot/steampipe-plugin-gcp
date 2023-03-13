package main

import (
	"github.com/turbot/steampipe-plugin-gcp/gcp"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: gcp.Plugin})
}
