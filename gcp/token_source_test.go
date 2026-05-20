package gcp

import (
	"strings"
	"testing"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"golang.org/x/oauth2"
)

// TestConnectionConfigTokenSource_PicksUpInPlaceConfigMutation verifies that
// the source returns the CURRENT contents of Connection.Config on every
// Token call, not values captured at construction time. This is the core
// property that lets in-flight goroutines pick up rotated impersonate access
// tokens.
func TestConnectionConfigTokenSource_PicksUpInPlaceConfigMutation(t *testing.T) {
	tok1 := "ya29.ORIGINAL_IMPERSONATE_TOKEN"
	tok2 := "ya29.ROTATED_IMPERSONATE_TOKEN"

	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig(gcpConfig{
		ImpersonateAccessToken: &tok1,
	})
	src := &connectionConfigTokenSource{connection: conn}

	got1, err := src.Token()
	if err != nil {
		t.Fatalf("first Token failed: %v", err)
	}
	if got1.AccessToken != tok1 {
		t.Errorf("first Token returned wrong AccessToken: got %q, want %q", got1.AccessToken, tok1)
	}

	// Rotate the connection's config in place, mimicking what the plugin SDK
	// does in upsertConnectionData when UpdateConnectionConfigs delivers a
	// rotated impersonate_access_token.
	conn.SetConfig(gcpConfig{
		ImpersonateAccessToken: &tok2,
	})

	got2, err := src.Token()
	if err != nil {
		t.Fatalf("second Token failed: %v", err)
	}
	if got2.AccessToken != tok2 {
		t.Errorf("second Token returned wrong AccessToken after rotation: got %q, want %q", got2.AccessToken, tok2)
	}
}

func TestConnectionConfigTokenSource_ErrorsWhenImpersonateAccessTokenMissing(t *testing.T) {
	conn := &plugin.Connection{Name: "my-test-conn"}
	conn.SetConfig(gcpConfig{})
	src := &connectionConfigTokenSource{connection: conn}
	_, err := src.Token()
	if err == nil {
		t.Fatalf("Token should error when impersonate_access_token is missing, got nil")
	}
	if !strings.Contains(err.Error(), "impersonate_access_token") {
		t.Errorf("error should mention which field is missing, got: %v", err)
	}
	if !strings.Contains(err.Error(), conn.Name) {
		t.Errorf("error should include the connection name %q for diagnostics, got: %v", conn.Name, err)
	}
}

func TestConnectionConfigTokenSource_ErrorsWhenConnectionNil(t *testing.T) {
	src := &connectionConfigTokenSource{connection: nil}
	if _, err := src.Token(); err == nil {
		t.Errorf("Token should error when connection is nil, got nil")
	}
}

func TestConnectionConfigTokenSource_HandlesUnexpectedConfigType(t *testing.T) {
	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig("this is not a gcpConfig")
	src := &connectionConfigTokenSource{connection: conn}
	_, err := src.Token()
	if err == nil {
		t.Errorf("Token should error when config is not gcpConfig, got nil")
	}
	if !strings.Contains(err.Error(), "gcpConfig") {
		t.Errorf("error should mention expected type for diagnostics, got: %v", err)
	}
}

// TestConnectionConfigTokenSource_BearerTokenType guards against a regression
// where the returned token has an empty TokenType. The
// google.golang.org/api transport calls (*oauth2.Token).Type() which falls
// back to "Bearer" when TokenType is empty, but pinning it explicitly here
// keeps the wire-format contract stable across upstream changes.
func TestConnectionConfigTokenSource_BearerTokenType(t *testing.T) {
	tok := "ya29.SOMETOKEN"
	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig(gcpConfig{
		ImpersonateAccessToken: &tok,
	})
	src := &connectionConfigTokenSource{connection: conn}

	got, err := src.Token()
	if err != nil {
		t.Fatalf("Token failed: %v", err)
	}
	if got.TokenType != "Bearer" {
		t.Errorf("TokenType should be %q, got %q", "Bearer", got.TokenType)
	}
}

// TestConnectionConfigTokenSource_ExpiryIsShort guards against a regression
// where Expiry drifts back out to the original token TTL — which would
// defeat the rotation-pickup goal because the google.golang.org/api
// transport's cached token provider would not call Token again until then.
// The AWS analogue validated 60 seconds; we allow a generous bound here to
// permit small adjustments without breaking the test, but anything more
// than a few minutes would silently leak stale creds for that long.
func TestConnectionConfigTokenSource_ExpiryIsShort(t *testing.T) {
	tok := "ya29.SOMETOKEN"
	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig(gcpConfig{
		ImpersonateAccessToken: &tok,
	})
	src := &connectionConfigTokenSource{connection: conn}

	before := time.Now()
	got, err := src.Token()
	if err != nil {
		t.Fatalf("Token failed: %v", err)
	}

	if got.Expiry.IsZero() {
		t.Fatalf("Token.Expiry must be non-zero; a zero Expiry causes oauth2.ReuseTokenSource " +
			"semantics to reuse the same token forever, defeating the rotation-pickup goal")
	}

	maxAllowed := before.Add(5 * time.Minute)
	if got.Expiry.After(maxAllowed) {
		t.Errorf("Token.Expiry should be within ~5min of now to bound rotation latency; got %s (now=%s, max=%s). "+
			"A long Expiry defeats the rotation-pickup goal because the google.golang.org/api transport's "+
			"cached token provider won't call Token again until then.",
			got.Expiry.UTC().Format(time.RFC3339Nano),
			before.UTC().Format(time.RFC3339Nano),
			maxAllowed.UTC().Format(time.RFC3339Nano))
	}

	if !got.Valid() {
		t.Errorf("Token.Valid() should be true immediately after construction; got false")
	}
}

// TestConnectionConfigTokenSource_ReuseLayerPicksUpRotation wraps the source
// in oauth2.ReuseTokenSource — the same Reuse-style caching the
// google.golang.org/api transport applies internally — and verifies that a
// rotation in Connection.Config is visible to a Token call after the cached
// token's Expiry has passed. The unit tests above bypass the cache and only
// prove the source returns the right values when invoked directly.
//
// We shrink tokenExpiresInterval to 50ms so the test runs in tens of
// milliseconds instead of waiting out the production 60s window.
func TestConnectionConfigTokenSource_ReuseLayerPicksUpRotation(t *testing.T) {
	orig := tokenExpiresInterval
	tokenExpiresInterval = 50 * time.Millisecond
	defer func() { tokenExpiresInterval = orig }()

	tok1 := "ya29.ORIGINAL_IMPERSONATE_TOKEN"
	tok2 := "ya29.ROTATED_IMPERSONATE_TOKEN"

	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig(gcpConfig{
		ImpersonateAccessToken: &tok1,
	})
	src := &connectionConfigTokenSource{connection: conn}
	reuse := oauth2.ReuseTokenSource(nil, src)

	got1, err := reuse.Token()
	if err != nil {
		t.Fatalf("first reuse.Token failed: %v", err)
	}
	if got1.AccessToken != tok1 {
		t.Errorf("first reuse.Token: got %q, want %q", got1.AccessToken, tok1)
	}

	// Rotate the connection config in place — what the SDK does on
	// UpdateConnectionConfigs.
	conn.SetConfig(gcpConfig{
		ImpersonateAccessToken: &tok2,
	})

	// Sleep past the reuse layer's Expiry window so the next Token call goes
	// back to the underlying source rather than serving the cached value.
	time.Sleep(100 * time.Millisecond)

	got2, err := reuse.Token()
	if err != nil {
		t.Fatalf("second reuse.Token failed: %v", err)
	}
	if got2.AccessToken != tok2 {
		t.Errorf("after rotation, reuse layer returned stale token: got %q, want %q. "+
			"This is the ExpiredToken-under-rotation bug — the cached token provider held "+
			"the original token past its replacement in Connection.Config.",
			got2.AccessToken, tok2)
	}
}
