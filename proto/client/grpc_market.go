package client

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

// MarketClient retrieves market price information about types in the EVE universe.
type MarketClient struct {
	// This type must be initialized using the package-level New function.

	*bootstrap
}

// GetMarketPrices returns a slice of market prices for each of the given type IDs.
func (c *MarketClient) GetMarketPrices(typeID int, typeIDs ...int) ([]*model.MarketPrice, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewMarketPriceServiceClient(conn)
	ids := []int64{int64(typeID)}
	for _, id := range typeIDs {
		ids = append(ids, int64(id))
	}
	res, err := service.GetMarketPrice(
		context.Background(),
		&proto.GetMarketPriceRequest{
			Token:  &proto.Token{Identifier: c.token},
			TypeId: ids,
		})
	if err != nil {
		return nil, err
	}
	var results []*model.MarketPrice
	for _, p := range res.Prices {
		results = append(results, proto.ProtoToMarketPrice(p))
	}
	return results, nil
}

// GetMarketPrice returns the current market price for the given type ID.
func (c *MarketClient) GetMarketPrice(typeID int) (*model.MarketPrice, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewMarketPriceServiceClient(conn)
	res, err := service.GetMarketPrice(
		context.Background(),
		&proto.GetMarketPriceRequest{
			Token:  &proto.Token{Identifier: c.token},
			TypeId: []int64{int64(typeID)},
		})
	if err != nil {
		return nil, err
	}
	if p, ok := res.Prices[int64(typeID)]; ok {
		return proto.ProtoToMarketPrice(p), nil
	}
	return nil, errors.Errorf("expected grpc response to price for typeID %d, got none", typeID)
}
