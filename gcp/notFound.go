package gcp

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"google.golang.org/api/googleapi"
)

// function which returns an IsNotFoundErrorPredicate for GCP API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if gerr, ok := err.(*googleapi.Error); ok {
			return helpers.StringSliceContains(notFoundErrors, types.ToString(gerr.Code))
		}
		return false
	}
}
