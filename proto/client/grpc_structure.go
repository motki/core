package client

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/core/eveapi"
	"github.com/motki/core/proto"
)

func (c *GRPCClient) GetStructure(structureID int) (*eveapi.Structure, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInfoServiceClient(conn)
	res, err := service.GetStructure(
		context.Background(),
		&proto.GetStructureRequest{Token: &proto.Token{Identifier: c.token}, StructureId: int64(structureID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	return proto.ProtoToStructure(res.Structure), nil
}

func (c *GRPCClient) GetCorpStructures() ([]*eveapi.CorporationStructure, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewCorporationServiceClient(conn)
	res, err := service.GetCorpStructures(
		context.Background(),
		&proto.GetCorpStructuresRequest{Token: &proto.Token{Identifier: c.token}})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var strucs []*eveapi.CorporationStructure
	for _, k := range res.Structures {
		strucs = append(strucs, proto.ProtoToCorpStructure(k))
	}
	return strucs, nil
}
