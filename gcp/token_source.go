package gcp

import (
	"errors"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"golang.org/x/oauth2"
)

// connectionConfigTokenSource is an oauth2.TokenSource that reads
// impersonate_access_token from the steampipe connection config on every
// Token() call (subject to the google.golang.org/api transport's cached
// token provider, which we bound via a short Expiry below).
//
// Background
//
// The steampipe plugin SDK mutates the connection config in place when a new
// connection config arrives via UpdateConnectionConfigs (see
// steampipe-plugin-sdk/plugin/plugin_connection_config.go upsertConnectionData,
// which calls d.Connection.SetConfig(configStruct) under a write lock). The
// SDK comment at that site explicitly acknowledges that a query may already be
// executing with this Connection object.
//
// An oauth2.StaticTokenSource built once at client construction time captures
// the original token value. A goroutine holding a GCP API client backed by
// that static source keeps signing requests with the original token regardless
// of rotation. When the original token expires at Google, every subsequent
// request from that goroutine fails with a 401 / ExpiredToken-style error —
// even if a fresh valid token has been delivered to Connection.Config by the
// SDK.
//
// By re-reading the connection config on every Token() call, in-flight
// goroutines holding the same GCP API client pick up rotated tokens on the
// next signing operation (modulo the cached token provider TTL we set via
// Expiry below).
type connectionConfigTokenSource struct {
	connection *plugin.Connection
}

// tokenExpiresInterval is how long the returned *oauth2.Token is considered
// fresh by the google.golang.org/api transport's cached token provider.
// Overridable in tests.
var tokenExpiresInterval = 60 * time.Second

// Token implements oauth2.TokenSource. Called by the google.golang.org/api
// transport on every signed request, subject to the wrapping cached token
// provider.
func (s *connectionConfigTokenSource) Token() (tok *oauth2.Token, err error) {
	// Belt-and-suspenders: the SDK's Connection.GetConfig acquires an RLock so
	// torn interface reads cannot happen via that path. This recover still
	// converts any other panic inside Token into a clean error the
	// google.golang.org/api transport can surface, rather than propagating
	// through the transport middleware.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("connectionConfigTokenSource: panic during Token for connection %q: %v", s.connectionName(), r)
			tok = nil
		}
	}()

	if s.connection == nil {
		return nil, errors.New("connectionConfigTokenSource: connection is nil")
	}

	// Read the raw config through the SDK's Connection.GetConfig accessor
	// (which acquires the RLock) and type-assert directly here. Avoid the
	// local gcp/connection_config.go GetConfig helper for symmetry with the
	// AWS implementation — keeping the signing path independent of any
	// helper-level normalization makes the rotation contract easier to
	// reason about.
	raw := s.connection.GetConfig()
	cfg, ok := raw.(gcpConfig)
	if !ok {
		return nil, fmt.Errorf("connectionConfigTokenSource: connection %q config is %T, expected gcpConfig", s.connection.Name, raw)
	}

	if cfg.ImpersonateAccessToken == nil {
		return nil, fmt.Errorf("connectionConfigTokenSource: connection %q has no impersonate_access_token in config", s.connection.Name)
	}

	// The cached token provider that the google.golang.org/api transport
	// wraps the TokenSource in will NOT call Token again until the returned
	// token's Expiry has passed, so setting Expiry too far out defeats the
	// rotation-pickup goal: the cache would hold the original token in
	// memory long after it was rotated in Connection.Config.
	//
	// 60 seconds mirrors the AWS credentials provider's Expires interval
	// and keeps rotation latency bounded. Reading Connection.Config is an
	// in-memory type assertion + struct copy, so a short interval is
	// essentially free.
	return &oauth2.Token{
		AccessToken: *cfg.ImpersonateAccessToken,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(tokenExpiresInterval),
	}, nil
}

func (s *connectionConfigTokenSource) connectionName() string {
	if s.connection == nil {
		return "<nil>"
	}
	return s.connection.Name
}
