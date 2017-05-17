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
	client *goesi.APIClient
	conf   Config
}

// New creates a new EveAPI with the given configuration.
func New(c Config) *EveAPI {
	t := httpcache.NewMemoryCacheTransport()
	t.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	hc := &http.Client{Transport: t}
	return &EveAPI{client: goesi.NewAPIClient(hc, c.UserAgent), conf: c}
}
