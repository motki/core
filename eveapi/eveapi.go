package eveapi

//import (
//	"github.com/antihax/goesi"
//)

type Config struct {
	ClientID  string `toml:"client_id"`
	SecretKey string `toml:"secret_key"`
	ReturnURL string `toml:"return_url"`
	Scopes    []string `toml:"scopes"`
}

func NewConfig(clientID, secretKey string) *Config {
	return &Config{
		ClientID:  clientID,
		SecretKey: secretKey,
		Scopes:    AllScopes,
	}
}
