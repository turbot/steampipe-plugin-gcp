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

// A limiter stands in for one API quota, so it has to pace exactly the calls that quota counts: a
// limiter reaching too wide would pace calls that are counted somewhere else. Where the quota is
// one budget for several methods, the scope leaves out the action so they share one bucket: a
// limiter per method would let each one fill the whole quota on its own.
func TestQuotaRateLimiters(t *testing.T) {
	tests := []struct {
		limiter string
		// The calls a minute allowed by the quota this limiter stands in for.
		quota float64
		scope []string
		// Tags of a hydrate call, against whether this limiter should pace it.
		calls map[string]bool
	}{
		{
			limiter: "gcp_billing",
			// cloudbilling.googleapis.com/standard_quota, counted across the service.
			quota: 300,
			scope: []string{"connection", "service"},
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
			limiter: "gcp_cloudresourcemanager_projects_get_ancestry",
			quota:   600,
			scope:   []string{"connection", "service", "action"},
			calls: map[string]bool{
				"resourcemanager/projects.getAncestry": true,
				// Paced by gcp_resourcemanager. projects.list has a quota of its own, while
				// projects.getIamPolicy is counted against the same read quota as getAncestry.
				"resourcemanager/projects.list":         false,
				"resourcemanager/projects.getIamPolicy": false,
				// The Access Approval API, under its own quota.
				"accessapproval/projects.getAccessApprovalSettings": false,
				"billing/projects.getBillingInfo":                   false,
			},
		},
		{
			limiter: "gcp_accessapproval",
			// accessapproval.googleapis.com/requests, counted across the service.
			quota: 300,
			scope: []string{"connection", "service"},
			calls: map[string]bool{
				"accessapproval/projects.getAccessApprovalSettings": true,
				"resourcemanager/projects.getAncestry":              false,
				"billing/projects.getBillingInfo":                   false,
			},
		},
	}

	for _, test := range tests {
		definition := rateLimiterDefinition(t, test.limiter)

		if !reflect.DeepEqual(definition.Scope, test.scope) {
			t.Errorf("%s scope is %v, expected %v", test.limiter, definition.Scope, test.scope)
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

// The project hydrates run once per project against a per minute quota. Each one is only paced
// when its hydrate config carries the tags its limiter filters on, and only retries past an
// exhausted quota window when that config carries its own retry settings, since the SDK takes
// nothing but the predicate from the plugin's default. Neither fails when missing: the calls just
// run uncapped, or retry for a fraction of the window.
func TestProjectHydratesArePacedAndRetried(t *testing.T) {
	hydrates := []struct {
		name    string
		hydrate plugin.HydrateFunc
		limiter string
	}{
		{"getProjectAccessApprovalSettings", getProjectAccessApprovalSettings, "gcp_accessapproval"},
		{"getProjectAncestors", getProjectAncestors, "gcp_cloudresourcemanager_projects_get_ancestry"},
		{"getProjectBillingInfo", getProjectBillingInfo, "gcp_billing"},
	}

	for _, tableName := range []string{"gcp_project", "gcp_organization_project"} {
		table := Plugin(context.Background()).TableMap[tableName]
		if table == nil {
			t.Fatalf("plugin has no table named %s", tableName)
		}

		for _, h := range hydrates {
			var config *plugin.HydrateConfig
			for i := range table.HydrateConfig {
				if sameHydrateFunc(table.HydrateConfig[i].Func, h.hydrate) {
					config = &table.HydrateConfig[i]
					break
				}
			}
			if config == nil {
				t.Errorf("%s has no hydrate config for %s, so its calls are neither paced nor retried", tableName, h.name)
				continue
			}

			scopeValues := map[string]string{"connection": "gcp", "service": config.Tags["service"], "action": config.Tags["action"]}
			if !rateLimiterDefinition(t, h.limiter).SatisfiesFilters(scopeValues) {
				t.Errorf("%s tags %s as %v, which the %s limiter does not match", tableName, h.name, config.Tags, h.limiter)
			}

			retry := config.RetryConfig
			switch {
			case retry == nil:
				t.Errorf("%s has no retry config for %s", tableName, h.name)
			case retry.ShouldRetryErrorFunc == nil:
				// The SDK would fill in the plugin's predicate, and keep these settings.
				t.Errorf("%s retry config for %s has no predicate of its own", tableName, h.name)
			case retry.MaxAttempts <= 0 || retry.BackoffAlgorithm != "Exponential":
				t.Errorf("%s retries %s %d times with a %q backoff, expected a bounded exponential one",
					tableName, h.name, retry.MaxAttempts, retry.BackoffAlgorithm)
			default:
				// The SDK waits RetryInterval doubled on each retry, capped at CappedDuration if set.
				var waitedMs int64
				for i := int64(0); i < retry.MaxAttempts; i++ {
					wait := retry.RetryInterval << i
					if retry.CappedDuration > 0 {
						wait = min(wait, retry.CappedDuration)
					}
					waitedMs += wait
				}
				if waitedMs < 60000 {
					t.Errorf("%s retries %s over %d ms, short of the one minute a quota window lasts", tableName, h.name, waitedMs)
				}
			}
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
