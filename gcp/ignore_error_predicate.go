package gcp

import (
	"context"
	"errors"
	"net/http"
	"path"
	"regexp"
	"slices"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"google.golang.org/api/googleapi"
)

// function which returns an isIgnorableErrorPredicate for GCP API calls
func isIgnorableError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if gerr, ok := err.(*googleapi.Error); ok {
			return slices.Contains(notFoundErrors, types.ToString(gerr.Code))
		}
		return false
	}
}

// isIgnorableErrorFunc is isIgnorableError as an ErrorPredicateWithContext, which is what
// IgnoreConfig.ShouldIgnoreErrorFunc takes; the context, query data and hydrate data are not needed
// to match an error code.
func isIgnorableErrorFunc(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	shouldIgnore := isIgnorableError(notFoundErrors)
	return func(_ context.Context, _ *plugin.QueryData, _ *plugin.HydrateData, err error) bool {
		return shouldIgnore(err)
	}
}

// shouldIgnoreErrorPluginDefault:: Plugin level default function to ignore a set errors for hydrate functions based on "ignore_error_codes" and "ignore_error_messages" config argument
func shouldIgnoreErrorPluginDefault() plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		gcpConfig := GetConfig(d.Connection)

		logger := plugin.Logger(ctx)

		// Add to support regex match as per error message
		for _, pattern := range gcpConfig.IgnoreErrorMessages {
			// Validate regex pattern
			re, er := regexp.Compile(pattern)
			if er != nil {
				panic(er.Error() + " the regex pattern configured in 'ignore_error_messages' is invalid. Edit your connection configuration file and then restart Steampipe")
			}
			result := re.MatchString(err.Error())
			if result {
				logger.Debug("ignore_error_predicate.shouldIgnoreErrorPluginDefault", "ignore_error_message", err.Error())
				return true
			}
		}

		if gerr, ok := err.(*googleapi.Error); ok {
			// Added to support regex in not found errors
			for _, pattern := range gcpConfig.IgnoreErrorCodes {
				if ok, _ := path.Match(pattern, types.ToString(gerr.Code)); ok {
					logger.Debug("ignore_error_predicate.shouldIgnoreErrorPluginDefault", "ignore_error_code", err.Error())
					return true
				}
			}
		}
		return false
	}
}

// Errors a GCP API returns for a condition that is expected to clear on its own: an exhausted
// per-minute quota, and the backend errors ("Internal error encountered", "The service is
// currently unavailable") any of them returns occasionally.
var retryableErrorCodes = []int{
	http.StatusTooManyRequests,
	http.StatusInternalServerError,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
}

// shouldRetryErrorPluginDefault:: Plugin level default function to retry the transient errors hydrate functions get, so that one failed call does not fail the whole query
func shouldRetryErrorPluginDefault() plugin.ErrorPredicateWithContext {
	// The SDK logs the retry decision and every attempt it then makes, so this only decides.
	return func(_ context.Context, _ *plugin.QueryData, _ *plugin.HydrateData, err error) bool {
		var gerr *googleapi.Error
		return errors.As(err, &gerr) && slices.Contains(retryableErrorCodes, gerr.Code)
	}
}

// quotaRetryConfig retries a call hydrated once per row against a per minute quota, where a
// failed call fails the whole query and replaying it costs one call per row again.
//
// It has to sit on each HydrateConfig: the SDK copies only the predicate of the plugin's
// DefaultRetryConfig, so a hydrate call left to it retries on the SDK's default backoff, about
// 23 seconds in all. Retries also skip the rate limiter, which paces only a call's first
// attempt, so the backoff alone has to carry them past an exhausted quota window: 2, 4, 8, 16
// then 30 seconds, one minute in all.
func quotaRetryConfig() *plugin.RetryConfig {
	return &plugin.RetryConfig{
		ShouldRetryErrorFunc: shouldRetryErrorPluginDefault(),
		MaxAttempts:          5,
		BackoffAlgorithm:     "Exponential",
		RetryInterval:        2000,
		CappedDuration:       30000,
	}
}
