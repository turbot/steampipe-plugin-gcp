package gcp

import (
	"context"
	"path"
	"regexp"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"google.golang.org/api/googleapi"
)

// function which returns an isIgnorableErrorPredicate for GCP API calls
func isIgnorableError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if gerr, ok := err.(*googleapi.Error); ok {
			return helpers.StringSliceContains(notFoundErrors, types.ToString(gerr.Code))
		}
		return false
	}
}

// shouldIgnoreErrorPluginDefault:: Plugin level default function to ignore a set errors for hydrate functions based on "ignore_error_codes" and "ignore_error_messages" config argument
func shouldIgnoreErrorPluginDefault() plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		gcpConfig := GetConfig(d.Connection)

		logger := plugin.Logger(ctx)
		logger.Trace("ignore_error_predicate.shouldIgnoreErrorPluginDefault", "err.Error()", err.Error())

		// Add to support regex match as per error message
		for _, pattern := range gcpConfig.IgnoreErrorMessages {
			logger.Trace("ignore_error_predicate.shouldIgnoreErrorPluginDefault", "message_pattern", pattern)
			// TODO: maybe compile these at the connection level for performance reasons
			re := regexp.MustCompile(pattern)
			result := re.MatchString(err.Error())
			if result {
				logger.Debug("ignore_error_predicate.shouldIgnoreErrorPluginDefault", "ignore_error_message", err.Error())
				return true
			}
		}

		if gerr, ok := err.(*googleapi.Error); ok {
			// Added to support regex in not found errors
			for _, pattern := range gcpConfig.IgnoreErrorCodes {
				logger.Trace("ignore_error_predicate.shouldIgnoreErrorPluginDefault", "code_pattern", pattern)
				if ok, _ := path.Match(pattern, types.ToString(gerr.Code)); ok {
					logger.Debug("ignore_error_predicate.shouldIgnoreErrorPluginDefault", "ignore_error_code", err.Error())
					return true
				}
			}
		}
		return false
	}
}
