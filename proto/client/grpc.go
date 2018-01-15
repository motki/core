package client

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"

	"github.com/motki/core/log"
	"github.com/motki/core/proto"
	"github.com/motki/core/proto/server"
)

// Ensure that GRPCClient implements the Client interface.
var _ Client = &GRPCClient{}

// GRPCClient communicates with a remote MOTKI installation using GRPC.
//
// A GRPCClient connects to a process-local or remote GRPC server to facilitate
// remote procedure calls. When used with a remote GRPC server, the Client allows
// client applications to consume MOTKI and EVESDE data without storing anything
// on the local machine.
type GRPCClient struct {
	serverAddr string
	token      string
	dialOpts   []grpc.DialOption
	logger     log.Logger
}

// newRemoteGRPC creates a new GRPC client intended for use with a remote GRPC server.
func newRemoteGRPC(serverAddr string, l log.Logger, tlsConf *tls.Config) (*GRPCClient, error) {
	return &GRPCClient{
		serverAddr: serverAddr,
		dialOpts:   []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(tlsConf))},
		logger:     l,
	}, nil
}

// newLocalGRPC creates a new GRPC client for use with a process-local GRPC server.
//
// The bufconn.Listener passed in should be shared between both client and server. By default,
// this is handled by the model.LocalConfig type.
func newLocalGRPC(lis *bufconn.Listener, l log.Logger) (*GRPCClient, error) {
	cl := &GRPCClient{logger: l}
	cl.dialOpts = append(cl.dialOpts, grpc.WithDialer(func(string, time.Duration) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
	return cl, nil
}

// Authenticate attempts to authenticate with the GRPC server.
//
// If authentication is successful, a token is stored in the client. The token is passed
// along with subsequent operations and used to authorize the user's access for things
// such as corporation-related functionality.
func (c *GRPCClient) Authenticate(username, password string) error {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	service := proto.NewAuthenticationServiceClient(conn)
	res, err := service.Authenticate(
		context.Background(),
		&proto.AuthenticateRequest{Username: username, Password: password})
	if err != nil {
		return err
	}
	if res.Result.Status == proto.Status_FAILURE {
		if res.Result.Description == server.ErrBadCredentials.Error() {
			return ErrBadCredentials
		}
		return errors.New(res.Result.Description)
	}
	if res.Token == nil || res.Token.Identifier == "" {
		return errors.New("expected token to be not empty, got nothing")
	}
	c.token = res.Token.Identifier
	return nil
}
