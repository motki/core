package client

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

// InventoryClient is the interface for corporation inventory management.
//
// Functionality provided by this client requires that the user's corporation
// is registered and opted-in to data collection.
type InventoryClient struct {
	// This type must be initialized using the package-level New function.

	*bootstrap
}

// GetInventory returns all inventory items for the current session's corporation.
//
// This method requires that the user's corporation has opted-in to data collection.
func (c *InventoryClient) GetInventory() ([]*model.InventoryItem, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInventoryServiceClient(conn)
	res, err := service.GetInventory(
		context.Background(),
		&proto.GetInventoryRequest{Token: &proto.Token{Identifier: c.token}})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var items []*model.InventoryItem
	for _, pr := range res.Item {
		items = append(items, proto.ProtoToInventoryItem(pr))
	}
	return items, nil
}

// NewInventoryItem creates a new inventory item for the given type ID and location ID.
//
// If an inventory item already exists for the given type and location ID, it will be returned.
//
// This method requires that the user's corporation has opted-in to data collection.
func (c *InventoryClient) NewInventoryItem(typeID, locationID int) (*model.InventoryItem, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInventoryServiceClient(conn)
	res, err := service.NewInventoryItem(
		context.Background(),
		&proto.NewInventoryItemRequest{
			Token:      &proto.Token{Identifier: c.token},
			TypeId:     int64(typeID),
			LocationId: int64(locationID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	if res.Item == nil {
		return nil, errors.New("expected grpc response to contain product, got nil")
	}
	return proto.ProtoToInventoryItem(res.Item), nil
}

// SaveInventoryItem attempts to save the given inventory item to the backend database.
//
// This method requires that the user's corporation has opted-in to data collection.
func (c *InventoryClient) SaveInventoryItem(item *model.InventoryItem) error {
	if c.token == "" {
		return ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	service := proto.NewInventoryServiceClient(conn)
	res, err := service.SaveInventoryItem(
		context.Background(),
		&proto.SaveInventoryItemRequest{
			Token: &proto.Token{Identifier: c.token},
			Item:  proto.InventoryItemToProto(item)})
	if err != nil {
		return err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return errors.New(res.Result.Description)
	}
	if res.Item == nil {
		return errors.New("expected grpc response to contain product, got nil")
	}
	return nil
}
