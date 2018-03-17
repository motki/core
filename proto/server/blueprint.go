package server

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

func (srv *grpcServer) GetCorpBlueprints(ctx context.Context, req *proto.GetCorpBlueprintsRequest) (resp *proto.GetCorpBlueprintsResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetCorpBlueprintsResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	if req.Token == nil {
		return nil, errors.New("token cannot be empty")
	}
	_, charID, err := srv.getAuthorizedContext(req.Token, model.RoleLogistics)
	if err != nil {
		return nil, err
	}
	char, err := srv.model.GetCharacter(charID)
	if err != nil {
		return nil, err
	}
	corpAuth, err := srv.model.GetCorporationAuthorization(char.CorporationID)
	if err != nil {
		return nil, err
	}
	bps, err := srv.model.GetCorporationBlueprints(corpAuth.Context(), char.CorporationID)
	if err != nil {
		return nil, err
	}
	var results []*proto.Blueprint
	for _, bp := range bps {
		results = append(results, proto.BlueprintToProto(bp))
	}
	return &proto.GetCorpBlueprintsResponse{
		Result:    successResult,
		Blueprint: results,
	}, nil
}
