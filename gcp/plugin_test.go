package gcp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
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

// Where one quota is a single budget for several methods, one limiter has to cover exactly those
// methods: a limiter per method would let each one fill the whole quota on its own, and a limiter
// reaching too wide would pace calls that are counted somewhere else.
func TestSharedQuotaRateLimiters(t *testing.T) {
	tests := []struct {
		limiter string
		// The calls a minute the API documents for the quota this limiter stands in for.
		quota float64
		// Tags of a hydrate call, against whether this limiter should pace it.
		calls map[string]bool
	}{
		{
			limiter: "gcp_billing",
			quota:   300,
			calls: map[string]bool{
				"billing/projects.getBillingInfo": true,
				"billing/accounts.get":            true,
				// The Cloud Billing Budget API is a separate service under its own quota.
				"billing/budgets.list":          false,
				"billing/budgets.get":           false,
				"resourcemanager/projects.list": false,
			},
		},
		{
			limiter: "gcp_cloudresourcemanager_projects_read",
			quota:   600,
			calls: map[string]bool{
				"resourcemanager/projects.getAncestry":               true,
				"resourcemanager/projects.getAccessApprovalSettings": true,
				// Counted against the per user quota the gcp_resourcemanager limiter covers.
				"resourcemanager/projects.list":         false,
				"resourcemanager/projects.getIamPolicy": false,
				"billing/projects.getBillingInfo":       false,
			},
		},
	}

	for _, test := range tests {
		definition := rateLimiterDefinition(t, test.limiter)

		if !reflect.DeepEqual(definition.Scope, []string{"connection", "service"}) {
			t.Errorf("%s scope is %v, so its methods do not share one bucket per connection", test.limiter, definition.Scope)
		}

		// A full bucket drains at once and refills over the same minute, so what has to stay
		// under the quota is the burst plus a minute of fill.
		if worstCase := float64(definition.BucketSize) + float64(definition.FillRate)*60; worstCase > test.quota {
			t.Errorf("%s: a burst of %d then %v per second puts up to %v calls in one minute, over the %v quota",
				test.limiter, definition.BucketSize, definition.FillRate, worstCase, test.quota)
		}

		for call, want := range test.calls {
			service, action, _ := strings.Cut(call, "/")
			scopeValues := map[string]string{"connection": "gcp", "service": service, "action": action}
			if matched := definition.SatisfiesFilters(scopeValues); matched != want {
				t.Errorf("%s matches %s: %t, expected %t", test.limiter, call, matched, want)
			}
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
