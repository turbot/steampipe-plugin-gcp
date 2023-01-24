package gcp

import (
	"context"
	"path"
	"regexp"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"google.golang.org/api/googleapi"
)

// function which returns an isIgnorableErrorPredicate for GCP API calls
func isIgnorableError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if gerr, ok := err.(*googleapi.Error); ok {
			if types.ToString(gerr.Code) == "403" {
				// return true, if service API is disabled
				regexExp := regexp.MustCompile(`googleapi: Error 403: [^\.]+ API has not been used in project [0-9]+ before or it is disabled\.`)
				return regexExp.MatchString(err.Error())
			}
			return helpers.StringSliceContains(notFoundErrors, types.ToString(gerr.Code))
		}
		return false
	}
}

// shouldIgnoreErrorPluginDefault:: Plugin level default function to ignore a set errors for hydrate functions based on "ignore_error_codes" config argument
func shouldIgnoreErrorPluginDefault() plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if !hasIgnoredErrorCodes(d.Connection) {
			return false
		}

		gcpConfig := GetConfig(d.Connection)
		if gerr, ok := err.(*googleapi.Error); ok {
			// Added to support regex in not found errors
			for _, pattern := range gcpConfig.IgnoreErrorCodes {
				if ok, _ := path.Match(pattern, types.ToString(gerr.Code)); ok {
					return true
				}
			}
		}
		return false
	}
}

func hasIgnoredErrorCodes(connection *plugin.Connection) bool {
	gcpConfig := GetConfig(connection)
	return len(gcpConfig.IgnoreErrorCodes) > 0
}
