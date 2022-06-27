package gcp

import (
	"context"
	"regexp"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
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

			for _, pattern := range notFoundErrors {
				if strings.Contains(err.Error(), pattern) {
					return true
				}
			}

		}
		return false
	}
}

// isNotFoundError:: function which returns an ErrorPredicate for Azure API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		gcpConfig := GetConfig(d.Connection)

		// If the get or list hydrate functions have an overriding IgnoreConfig
		// defined using the isNotFoundError function, then it should
		// also check for errors in the "ignore_error_codes" config argument
		allErrors := append(notFoundErrors, gcpConfig.IgnoreErrorCodes...)
		// Added to support regex in not found errors
		for _, pattern := range allErrors {
			if strings.Contains(err.Error(), pattern) {
				return true
			}
		}
		return false
	}
}
