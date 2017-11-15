package server

import (
	"github.com/antihax/goesi"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/motki/model"
	"github.com/motki/motki/proto"
)

func (srv *GRPCServer) getAuthorizedContext(tok *proto.Token, role model.Role) (context.Context, int, error) {
	if tok == nil || tok.Identifier == "" {
		return nil, 0, errors.New("token cannot be empty")
	}
	user, err := srv.model.GetUserBySessionKey(tok.Identifier)
	if err != nil {
		return nil, 0, err
	}
	a, err := srv.model.GetAuthorization(user, role)
	if err != nil {
		return nil, 0, err
	}
	source, err := srv.eveapi.TokenSource(a.Token)
	if err != nil {
		return nil, 0, err
	}
	info, err := srv.eveapi.Verify(source)
	if err != nil {
		return nil, 0, err
	}
	t, err := source.Token()
	if err != nil {
		return nil, 0, err
	}
	if err = srv.model.SaveAuthorization(user, role, int(info.CharacterID), t); err != nil {
		return nil, 0, err
	}
	return context.WithValue(context.Background(), goesi.ContextOAuth2, source), int(info.CharacterID), nil
}

func (srv *GRPCServer) GetCorpBlueprints(ctx context.Context, req *proto.GetCorpBlueprintsRequest) (resp *proto.GetCorpBlueprintsResponse, err error) {
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
	ctx, charID, err := srv.getAuthorizedContext(req.Token, model.RoleLogistics)
	if err != nil {
		return nil, err
	}
	char, err := srv.model.GetCharacter(charID)
	if err != nil {
		return nil, err
	}
	bps, err := srv.model.GetCorporationBlueprints(ctx, char.CorporationID)
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
