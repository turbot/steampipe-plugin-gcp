package gcp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/rate_limiter"
	"google.golang.org/api/googleapi"
)

// rateLimiterDefinition returns the named rate limiter of the plugin, initialised so that its
// "where" filter is parsed and can be evaluated.
func rateLimiterDefinition(t *testing.T, name string) *rate_limiter.Definition {
	t.Helper()

	for _, definition := range Plugin(context.Background()).RateLimiters {
		if definition.Name != name {
			continue
		}
		if errs := definition.Validate(); len(errs) > 0 {
			t.Fatalf("rate limiter %s is invalid: %v", name, errs)
		}
		if err := definition.Initialise(); err != nil {
			t.Fatalf("rate limiter %s has an unparseable where filter: %v", name, err)
		}
		return definition
	}

	t.Fatalf("plugin has no rate limiter named %s", name)
	return nil
}

// The Cloud Billing API counts every method against one per minute quota, so the limiter has to
// catch all of them, and only them: billing budgets are served by another API with its own quota.
func TestBillingRateLimiterTargetsTheCloudBillingMethods(t *testing.T) {
	definition := rateLimiterDefinition(t, "gcp_billing")

	if !reflect.DeepEqual(definition.Scope, []string{"connection", "service"}) {
		t.Fatalf("limiter scope is %v, so the methods do not share one bucket per connection", definition.Scope)
	}

	// A full bucket drains at once, so a minute can carry the burst on top of what refills in it,
	// and that total is what has to stay under the 300 calls a minute the API documents.
	if worstCase := float64(definition.BucketSize) + float64(definition.FillRate)*60; worstCase > 300 {
		t.Errorf("a burst of %d then %v per second puts up to %v calls in one minute, over the 300 quota", definition.BucketSize, definition.FillRate, worstCase)
	}

	tests := []struct {
		tags    map[string]string
		matched bool
	}{
		{map[string]string{"service": "billing", "action": "projects.getBillingInfo"}, true},
		{map[string]string{"service": "billing", "action": "accounts.get"}, true},
		{map[string]string{"service": "billing", "action": "budgets.list"}, false},
		{map[string]string{"service": "billing", "action": "budgets.get"}, false},
		{map[string]string{"service": "resourcemanager", "action": "projects.list"}, false},
	}
	for _, test := range tests {
		scopeValues := map[string]string{"connection": "gcp"}
		for key, value := range test.tags {
			scopeValues[key] = value
		}
		if matched := definition.SatisfiesFilters(scopeValues); matched != test.matched {
			t.Errorf("limiter matches %v: %t, expected %t", test.tags, matched, test.matched)
		}
	}
}

// A hydrate function is only paced by a rate limiter when its hydrate config carries the tags the
// limiter filters on, and nothing fails when they are missing: the calls just run uncapped.
func TestProjectBillingInfoHydrateIsTagged(t *testing.T) {
	for _, tableName := range []string{"gcp_project", "gcp_organization_project"} {
		table := Plugin(context.Background()).TableMap[tableName]
		if table == nil {
			t.Fatalf("plugin has no table named %s", tableName)
		}

		var tags map[string]string
		for _, hydrateConfig := range table.HydrateConfig {
			if sameHydrateFunc(hydrateConfig.Func, getProjectBillingInfo) {
				tags = hydrateConfig.Tags
				break
			}
		}

		if tags == nil {
			t.Errorf("%s has no hydrate config for getProjectBillingInfo, so its calls are not rate limited", tableName)
			continue
		}
		if tags["service"] != "billing" || tags["action"] != "projects.getBillingInfo" {
			t.Errorf("%s tags getProjectBillingInfo as %v, which the gcp_billing limiter does not match", tableName, tags)
		}
	}
}

func sameHydrateFunc(a plugin.HydrateFunc, b plugin.HydrateFunc) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}

// A hydrate error fails the whole query, so the transient ones have to be retried: an exhausted
// quota and a backend error both clear on their own, while a 403 or a 404 never will.
func TestShouldRetryErrorPluginDefault(t *testing.T) {
	shouldRetry := shouldRetryErrorPluginDefault()

	tests := []struct {
		err   error
		retry bool
	}{
		{&googleapi.Error{Code: 429, Message: "Quota exceeded for quota metric 'All requests'"}, true},
		{&googleapi.Error{Code: 500, Message: "Internal error encountered."}, true},
		{&googleapi.Error{Code: 503, Message: "The service is currently unavailable."}, true},
		// The google API client wraps its errors on some paths, so unwrap before giving up.
		{fmt.Errorf("getProjectBillingInfo: %w", &googleapi.Error{Code: 429}), true},
		{&googleapi.Error{Code: 403, Message: "Permission denied"}, false},
		{&googleapi.Error{Code: 404, Message: "Not found"}, false},
		{errors.New("context canceled"), false},
	}
	for _, test := range tests {
		if retry := shouldRetry(context.Background(), nil, nil, test.err); retry != test.retry {
			t.Errorf("retrying %v: %t, expected %t", test.err, retry, test.retry)
		}
	}
}
