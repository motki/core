package client

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

// StructureClient retrieves information about player-owned Citadels in EVE.
type StructureClient struct {
	// This type must be initialized using the package-level New function.

	*bootstrap
}

// GetStructure returns public information about the given structure.
func (c *StructureClient) GetStructure(structureID int) (*model.Structure, error) {
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

// GetCorpStructures returns detailed information about corporation structures.
//
// This method requires that the user's corporation has opted-in to data collection.
func (c *StructureClient) GetCorpStructures() ([]*model.CorporationStructure, error) {
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
	var strucs []*model.CorporationStructure
	for _, k := range res.Structures {
		strucs = append(strucs, proto.ProtoToCorpStructure(k))
	}
	return strucs, nil
}
