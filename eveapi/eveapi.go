// Package eveapi manages fetching and posting data to the EVE Swagger API.
package eveapi

import (
	"net/http"

	"github.com/antihax/goesi"
	"github.com/gregjones/httpcache"
)

type Config struct {
	ClientID  string `toml:"client_id"`
	SecretKey string `toml:"secret_key"`
	ReturnURL string `toml:"return_url"`
	UserAgent string `toml:"user_agent"`
}

type EveAPI struct {
	client *goesi.APIClient
	conf   Config
}

func New(c Config) *EveAPI {
	t := httpcache.NewMemoryCacheTransport()
	t.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	hc := &http.Client{Transport: t}
	return &EveAPI{client: goesi.NewAPIClient(hc, c.UserAgent), conf: c}
}
