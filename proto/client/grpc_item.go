package client

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/core/evedb"
	"github.com/motki/core/proto"
)

// GetItemType returns information about the given type ID.
func (c *GRPCClient) GetItemType(typeID int) (*evedb.ItemType, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetItemType(
		context.Background(),
		&proto.GetItemTypeRequest{TypeId: int64(typeID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	pres := res.Type
	if pres == nil {
		return nil, errors.New("expected item type in grpc response, got nil")
	}
	return proto.ProtoToItemType(pres), nil
}

// GetItemTypeDetail returns detailed information about the given type ID.
func (c *GRPCClient) GetItemTypeDetail(typeID int) (*evedb.ItemTypeDetail, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetItemTypeDetail(
		context.Background(),
		&proto.GetItemTypeDetailRequest{TypeId: int64(typeID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	pres := res.Type
	if pres == nil {
		return nil, errors.New("expected item type detail in grpc response, got nil")
	}
	return proto.ProtoToItemTypeDetail(pres), nil
}

// QueryItemTypes searches for matching item types given a search query and optional category IDs.
func (c *GRPCClient) QueryItemTypes(query string, catIDs ...int) ([]*evedb.ItemType, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	var cats []int64
	for _, cat := range catIDs {
		cats = append(cats, int64(cat))
	}
	res, err := service.QueryItemTypes(
		context.Background(),
		&proto.QueryItemTypesRequest{
			Query:      query,
			CategoryId: cats,
		})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var results []*evedb.ItemType
	for _, pr := range res.Types {
		results = append(results, proto.ProtoToItemType(pr))
	}
	return results, nil
}

// QueryItemTypeDetails searches for matching item types, returning detail type information for each match.
func (c *GRPCClient) QueryItemTypeDetails(query string, catIDs ...int) ([]*evedb.ItemTypeDetail, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	var cats []int64
	for _, cat := range catIDs {
		cats = append(cats, int64(cat))
	}
	res, err := service.QueryItemTypeDetails(
		context.Background(),
		&proto.QueryItemTypeDetailsRequest{
			Query:      query,
			CategoryId: cats,
		})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var results []*evedb.ItemTypeDetail
	for _, pr := range res.Types {
		results = append(results, proto.ProtoToItemTypeDetail(pr))
	}
	return results, nil
}

// GetMaterialSheet returns manufacturing information about the given type ID.
func (c *GRPCClient) GetMaterialSheet(typeID int) (*evedb.MaterialSheet, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetMaterialSheet(
		context.Background(),
		&proto.GetMaterialSheetRequest{TypeId: int64(typeID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	pres := res.MatSheet
	if pres == nil {
		return nil, errors.New("expected material sheet in grpc response, got nil")
	}
	return proto.ProtoToMatSheet(pres), nil
}
