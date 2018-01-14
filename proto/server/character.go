package server

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/motki/model"
	"github.com/motki/motki/proto"
)

func (srv *grpcServer) GetCharacter(ctx context.Context, req *proto.GetCharacterRequest) (resp *proto.CharacterResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.CharacterResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	var char *model.Character
	if charID := req.GetCharacterId(); charID != 0 {
		char, err = srv.model.GetCharacter(int(charID))
		if err != nil {
			return nil, err
		}
	} else {
		if req.Token == nil {
			return nil, errors.New("token cannot be empty")
		}
		var role model.Role = model.RoleLogistics
		switch req.GetRole() {
		case proto.Role_LOGISTICS:
			role = model.RoleLogistics
		case proto.Role_USER:
			role = model.RoleUser
		default:
			// no op
		}
		user, err := srv.model.GetUserBySessionKey(req.Token.Identifier)
		if err != nil {
			return nil, err
		}
		a, err := srv.model.GetAuthorization(user, role)
		if err != nil {
			return nil, err
		}
		char, err = srv.model.GetCharacter(a.CharacterID)
		if err != nil {
			return nil, err
		}
	}
	return &proto.CharacterResponse{
		Result:    successResult,
		Character: proto.CharacterToProto(char),
	}, nil
}

func (srv *grpcServer) GetCorporation(ctx context.Context, req *proto.GetCorporationRequest) (resp *proto.CorporationResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.CorporationResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	var corp *model.Corporation
	corp, err = srv.model.GetCorporation(int(req.GetCorporationId()))
	if err != nil {
		return nil, err
	}
	return &proto.CorporationResponse{
		Result:      successResult,
		Corporation: proto.CorporationToProto(corp),
	}, nil
}

func (srv *grpcServer) GetAlliance(ctx context.Context, req *proto.GetAllianceRequest) (resp *proto.AllianceResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.AllianceResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	var alliance *model.Alliance
	alliance, err = srv.model.GetAlliance(int(req.GetAllianceId()))
	if err != nil {
		return nil, err
	}
	return &proto.AllianceResponse{
		Result:   successResult,
		Alliance: proto.AllianceToProto(alliance),
	}, nil
}
