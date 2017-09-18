package model

import (
	"crypto/tls"

	"net"

	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc/test/bufconn"
)

type BackendKind string

var (
	BackendRemoteGRPC BackendKind = "remote_grpc"
	BackendLocalGRPC  BackendKind = "local_grpc"
)

type LocalConfig struct {
	Listener *bufconn.Listener
}

type Config struct {
	Kind       BackendKind  `toml:"kind"`
	RemoteGRPC RemoteConfig `toml:"remote_grpc"`
	LocalGRPC  LocalConfig  `toml:"local_grpc"`

	ServerEnabled bool         `toml:"enable_server"`
	ServerGRPC    ServerConfig `toml:"server_grpc"`
}

type RemoteConfig struct {
	ServerAddr         string `toml:"addr"`
	InsecureSkipVerify bool   `toml:"insecure_skip_verify_ssl"`
}

type ServerConfig struct {
	RemoteConfig
	AutoCert   bool     `toml:"autocert"`
	CertFile   string   `toml:"certfile"`
	CertKey    string   `toml:"keyfile"`
	ExtraHosts []string `toml:"extra_hosts"`

	listenHost string
	listenPort string
}

// TLSConfig attempts to load the configured certificate.
func (c RemoteConfig) TLSConfig() (*tls.Config, error) {
	return &tls.Config{NextProtos: []string{"h2"}, InsecureSkipVerify: c.InsecureSkipVerify}, nil
}

// TLSConfig attempts to load the configured certificate.
func (c ServerConfig) TLSConfig() (*tls.Config, error) {
	var err error
	c.listenHost, c.listenPort, err = net.SplitHostPort(c.ServerAddr)
	if err != nil {
		return nil, err
	}
	if c.AutoCert {
		hosts := append([]string{c.listenHost}, c.ExtraHosts...)
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(hosts...),
			Cache:      autocert.DirCache("certs"),
		}
		return &tls.Config{GetCertificate: m.GetCertificate}, nil
	}
	tc := &tls.Config{}
	tc.InsecureSkipVerify = c.InsecureSkipVerify
	tc.NextProtos = []string{"h2"}
	tc.Certificates = make([]tls.Certificate, 1)
	tc.Certificates[0], err = tls.LoadX509KeyPair(c.CertFile, c.CertKey)
	if err != nil {
		return nil, err
	}
	return tc, nil
}
