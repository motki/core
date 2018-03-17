package server

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

func (srv *grpcServer) GetStructure(ctx context.Context, req *proto.GetStructureRequest) (resp *proto.GetStructureResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetStructureResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.model.GetStructure(ctx, int(req.StructureId))
	if err != nil {
		return nil, err
	}
	return &proto.GetStructureResponse{
		Result:    successResult,
		Structure: proto.StructureToProto(res),
	}, nil
}

func (srv *grpcServer) GetCorpStructures(ctx context.Context, req *proto.GetCorpStructuresRequest) (resp *proto.GetCorpStructuresResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetCorpStructuresResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	if req.Token == nil {
		return nil, errors.New("token cannot be empty")
	}
	user, err := srv.model.GetUserBySessionKey(req.Token.Identifier)
	if err != nil {
		return nil, err
	}
	a, err := srv.model.GetAuthorization(user, model.RoleLogistics)
	if err != nil {
		return nil, err
	}
	res, err := srv.model.GetCorporationStructures(ctx, a.CorporationID)
	if err != nil {
		return nil, err
	}
	var strucs []*proto.CorporationStructure
	for _, s := range res {
		strucs = append(strucs, proto.CorpStructureToProto(s))
	}
	return &proto.GetCorpStructuresResponse{
		Result:     successResult,
		Structures: strucs,
	}, nil
}
