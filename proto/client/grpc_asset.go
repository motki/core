package client

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

// GetCorpBlueprints returns the current session's corporation's blueprints.
func (c *GRPCClient) GetCorpBlueprints() ([]*model.Blueprint, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewCorporationServiceClient(conn)
	res, err := service.GetCorpBlueprints(
		context.Background(),
		&proto.GetCorpBlueprintsRequest{
			Token: &proto.Token{Identifier: c.token},
		})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var bps []*model.Blueprint
	for _, bp := range res.Blueprint {
		bps = append(bps, proto.ProtoToBlueprint(bp))
	}
	return bps, nil
}
