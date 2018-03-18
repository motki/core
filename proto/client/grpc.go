package client

import (
	"crypto/tls"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"

	"github.com/motki/core/cache"
	"github.com/motki/core/log"
)

// Ensure that GRPCClient implements the Client interface.
var _ Client = &GRPCClient{}

// bootstrap communicates with a remote MOTKI installation using GRPC.
//
// A bootstrap connects to a process-local or remote GRPC server to facilitate
// remote procedure calls. When used with a remote GRPC server, the Client allows
// client applications to consume MOTKI and EVESDE data without storing anything
// on the local machine.
type bootstrap struct {
	serverAddr string
	token      string
	dialOpts   []grpc.DialOption
	logger     log.Logger
}

// GRPCClient is the defacto implementation of the Client interface.
type GRPCClient struct {
	*AssetClient
	*CharacterClient
	*EVEUniverseClient
	*InventoryClient
	*ItemTypeClient
	*LocationClient
	*MarketClient
	*ProductClient
	*StructureClient
	*UserClient

	// This type must be initialized using the package-level New function.

	noexport struct{} // Don't allow other packages to initialize this struct.
}

func newGRPCClient(m *bootstrap) *GRPCClient {
	return &GRPCClient{
		AssetClient:       &AssetClient{m},
		CharacterClient:   &CharacterClient{m},
		EVEUniverseClient: &EVEUniverseClient{m},
		InventoryClient:   &InventoryClient{m},
		ItemTypeClient:    &ItemTypeClient{m},
		LocationClient:    &LocationClient{m},
		MarketClient:      &MarketClient{m},
		ProductClient:     &ProductClient{m},
		StructureClient:   &StructureClient{m},
		UserClient:        &UserClient{m},
	}
}

// newRemoteGRPC creates a new GRPC client intended for use with a remote GRPC server.
func newRemoteGRPC(serverAddr string, l log.Logger, tlsConf *tls.Config) (*cachingGRPCClient, error) {
	m := &bootstrap{
		serverAddr: serverAddr,
		dialOpts:   []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(tlsConf))},
		logger:     l,
	}
	return &cachingGRPCClient{newGRPCClient(m), cache.New(cacheTTL)}, nil
}

// newLocalGRPC creates a new GRPC client for use with a process-local GRPC server.
//
// The bufconn.Listener passed in should be shared between both client and server. By default,
// this is handled by the model.LocalConfig type.
func newLocalGRPC(lis *bufconn.Listener, l log.Logger) (*cachingGRPCClient, error) {
	m := &bootstrap{
		logger: l,
		dialOpts: []grpc.DialOption{
			grpc.WithDialer(func(string, time.Duration) (net.Conn, error) {
				return lis.Dial()
			}),
			grpc.WithInsecure()}}
	cl := &cachingGRPCClient{newGRPCClient(m), cache.New(cacheTTL)}
	return cl, nil
}
