// Package eveapi manages fetching and posting data to the EVE Swagger API.
package eveapi // import "github.com/motki/core/eveapi"

import (
	"net/http"

	"github.com/motki/core/log"

	"github.com/antihax/goesi"
	"github.com/gregjones/httpcache"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
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

	logger log.Logger
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

		logger: l,
	}
}

func (api *EveAPI) AuthorizeURL(state string, scopes ...string) string {
	if len(scopes) == 0 {
		scopes = AllScopes
	}
	return api.ssoAuth.AuthorizeURL(state, true, scopes)
}

func (api *EveAPI) TokenExchange(code string) (*oauth2.Token, error) {
	return api.ssoAuth.TokenExchange(code)
}

func (api *EveAPI) TokenSource(tok *oauth2.Token) (oauth2.TokenSource, error) {
	return api.ssoAuth.TokenSource(tok)
}

func (api *EveAPI) Verify(source oauth2.TokenSource) (*goesi.VerifyResponse, error) {
	return api.ssoAuth.Verify(source)
}

func TokenFromContext(ctx context.Context) (oauth2.TokenSource, error) {
	if v, ok := ctx.Value(goesi.ContextOAuth2).(oauth2.TokenSource); ok {
		return v, nil
	}
	return nil, ErrNoToken
}
