// Package eveapi manages fetching and posting data to the EVE Swagger API.
package eveapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/antihax/goesi"
	"github.com/gregjones/httpcache"
	"github.com/motki/motki/log"
)

var ErrNoToken = errors.New("unable to get token from context")

// Config represents the configuration for a EVE Swagger API client.
type Config struct {
	ClientID  string `toml:"client_id"`
	SecretKey string `toml:"secret_key"`
	ReturnURL string `toml:"return_url"`
	UserAgent string `toml:"user_agent"`
}

// EveAPI is the entry point for interacting with the EVE Swagger API.
type EveAPI struct {
	client  *goesi.APIClient
	ssoAuth *goesi.SSOAuthenticator
}

// New creates a new EveAPI with the given configuration.
func New(c Config, l log.Logger) *EveAPI {
	l.Debugf("eveapi: init with EVE Developer Portal application client ID: %s", c.ClientID)
	l.Debugf("eveapi: SSO return URL: %s", c.ReturnURL)
	l.Debugf("eveapi: API client user agent: %s", c.UserAgent)
	t := httpcache.NewMemoryCacheTransport()
	t.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	hc := &http.Client{Transport: t}
	return &EveAPI{
		client:  goesi.NewAPIClient(hc, c.UserAgent),
		ssoAuth: goesi.NewSSOAuthenticator(hc, c.ClientID, c.SecretKey, c.ReturnURL, AllScopes),
	}
}

func (api *EveAPI) AuthorizeURL(state string, scopes ...string) string {
	if len(scopes) == 0 {
		scopes = AllScopes
	}
	return api.ssoAuth.AuthorizeURL(state, true, scopes)
}

func (api *EveAPI) TokenExchange(code string) (*goesi.CRESTToken, error) {
	return api.ssoAuth.TokenExchange(code)
}

func (api *EveAPI) TokenSource(tok *goesi.CRESTToken) (goesi.CRESTTokenSource, error) {
	return api.ssoAuth.TokenSource(tok)
}

func (api *EveAPI) Verify(source goesi.CRESTTokenSource) (*goesi.VerifyResponse, error) {
	return api.ssoAuth.Verify(source)
}

func TokenFromContext(ctx context.Context) (goesi.CRESTTokenSource, error) {
	if v, ok := ctx.Value(goesi.ContextOAuth2).(goesi.CRESTTokenSource); ok {
		return v, nil
	}
	return nil, ErrNoToken
}
