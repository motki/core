// Package eveapi manages fetching and posting data to the EVE Swagger API.
package eveapi

import (
	"net/http"

	"github.com/antihax/goesi"
	"github.com/gregjones/httpcache"
)

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
func New(c Config) *EveAPI {
	t := httpcache.NewMemoryCacheTransport()
	t.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	hc := &http.Client{Transport: t}
	return &EveAPI{
		client:  goesi.NewAPIClient(hc, c.UserAgent),
		ssoAuth: goesi.NewSSOAuthenticator(hc, c.ClientID, c.SecretKey, c.ReturnURL, AllScopes),
	}
}

func (api *EveAPI) AuthorizeURL(state string) string {
	return api.ssoAuth.AuthorizeURL(state, true, AllScopes)
}

func (api *EveAPI) TokenExchange(code string) (*goesi.CRESTToken, error) {
	return api.ssoAuth.TokenExchange(code)
}

func (api *EveAPI) Verify(tok *goesi.CRESTToken) (*goesi.VerifyResponse, error) {
	source, err := api.ssoAuth.TokenSource(tok)
	if err != nil {
		return nil, err
	}
	return api.ssoAuth.Verify(source)
}
