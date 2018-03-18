package client

import (
	"github.com/motki/core/proto"
	"github.com/motki/core/proto/server"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// UserClient handles authenticating a user session.
type UserClient struct {
	// This type must be initialized using the package-level New function.

	*bootstrap
}

// Authenticate attempts to authenticate with the GRPC server.
//
// If authentication is successful, a token is stored in the client. The token is passed
// along with subsequent operations and used to authorize the user's access for things
// such as corporation-related functionality.
func (c *UserClient) Authenticate(username, password string) error {
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

// Authenticated returns true if the current session is authenticated.
func (c *UserClient) Authenticated() bool {
	return c.token != ""
}
