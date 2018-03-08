package server

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

func (srv *grpcServer) GetInventory(ctx context.Context, req *proto.GetInventoryRequest) (resp *proto.GetInventoryResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetInventoryResponse{
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
	corp, err := srv.model.GetCorporation(char.CorporationID)
	if err != nil {
		return nil, err
	}
	corpAuth, err := srv.model.GetCorporationAuthorization(char.CorporationID)
	if err != nil {
		return nil, err
	}
	items, err := srv.model.GetCorporationInventory(corpAuth.Context(), corp.CorporationID)
	if err != nil {
		return nil, err
	}
	its := make([]*proto.InventoryItem, len(items))
	for i, v := range items {
		its[i] = proto.InventoryItemToProto(v)
	}
	return &proto.GetInventoryResponse{
		Result: successResult,
		Item:   its,
	}, nil
}

func (srv *grpcServer) NewInventoryItem(ctx context.Context, req *proto.NewInventoryItemRequest) (resp *proto.InventoryItemResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.InventoryItemResponse{
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
	corp, err := srv.model.GetCorporation(char.CorporationID)
	if err != nil {
		return nil, err
	}
	corpAuth, err := srv.model.GetCorporationAuthorization(char.CorporationID)
	if err != nil {
		return nil, err
	}
	item, err := srv.model.NewInventoryItem(corpAuth.Context(), corp.CorporationID, int(req.TypeId), int(req.LocationId))
	if err != nil {
		return nil, err
	}
	return &proto.InventoryItemResponse{
		Result: successResult,
		Item:   proto.InventoryItemToProto(item),
	}, nil
}

func (srv *grpcServer) SaveInventoryItem(ctx context.Context, req *proto.SaveInventoryItemRequest) (resp *proto.InventoryItemResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.InventoryItemResponse{
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
	it := proto.ProtoToInventoryItem(req.Item)
	it.CorporationID = char.CorporationID
	if err := srv.model.SaveInventoryItem(corpAuth.Context(), it); err != nil {
		return nil, err
	}
	return &proto.InventoryItemResponse{
		Result: successResult,
		Item:   req.Item,
	}, nil
}
