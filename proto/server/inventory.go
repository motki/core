package server

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/motki/motki/model"
	"github.com/motki/motki/proto"
)

func inventoryResponse(items []*model.InventoryItem) *proto.GetCorpInventoryResponse {
	its := make([]*proto.InventoryItem, len(items))
	for i, v := range items {
		its[i] = proto.InventoryItemToProto(v)
	}
	return &proto.GetCorpInventoryResponse{
		Result: successResult,
		Item:   its,
	}
}

func (srv *GRPCServer) GetCorpInventory(ctx context.Context, req *proto.GetCorpInventoryRequest) (resp *proto.GetCorpInventoryResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetCorpInventoryResponse{
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
	prod, err := srv.model.GetCorporationInventory(corpAuth.Context(), corp.CorporationID)
	if err != nil {
		return nil, err
	}
	return inventoryResponse(prod), nil
}
