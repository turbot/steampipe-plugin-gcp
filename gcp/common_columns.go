package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// column definitions for the common columns
var commonGCPRegionalColumns = []*plugin.Column{
	{
		Name:        "location",
		Description: "The GCP region in which the resource is located",
		Type:        proto.ColumnType_STRING,
		Hydrate:     getCommonColumns,
	},
	{
		Name:        "project",
		Description: "The Google Project in which the resource is located",
		Type:        proto.ColumnType_STRING,
		Hydrate:     getCommonColumns,
	},
}

var commonGCPColumns = []*plugin.Column{
	{
		Name:        "project",
		Description: "The Google Project in which the resource is located",
		Type:        proto.ColumnType_STRING,
		Hydrate:     getCommonColumns,
	},
}

// append the common gcp columns for REGIONAL resources onto the column list
func gcpRegionalColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonGCPRegionalColumns...)
}

// append the common gcp columns for GLOBAL resources onto the column list
func gcpColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonGCPColumns...)
}

// struct to store the common column data
type gcpCommonColumnData struct {
	Location, Project string
}

// get columns which are returned with all tables: location and project
func getCommonColumns(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cacheKey := "commonColumnData"
	var commonColumnData *gcpCommonColumnData
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		commonColumnData = cachedData.(*gcpCommonColumnData)
	} else {
		commonColumnData = &gcpCommonColumnData{
			Project: activeProject(),
			// Location: GetDefaultLocation(),
		}

		// save to extension cache
		d.ConnectionManager.Cache.Set(cacheKey, commonColumnData)
	}
	plugin.Logger(ctx).Trace("getCommonColumns__", "commonColumnData", commonColumnData)
	return commonColumnData, nil
}
