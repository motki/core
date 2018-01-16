package server

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/core/proto"
)

func (srv *grpcServer) GetMarketPrice(ctx context.Context, req *proto.GetMarketPriceRequest) (resp *proto.GetMarketPriceResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetMarketPriceResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	if req.Token == nil {
		return nil, errors.New("token cannot be empty")
	}
	if len(req.TypeId) == 0 {
		return nil, errors.New("must pass at least one type ID")
	}
	var ids []int
	for _, id := range req.TypeId {
		ids = append(ids, int(id))
	}
	prices, err := srv.model.GetMarketPrices(ids[0], ids[1:]...)
	if err != nil {
		return nil, err
	}
	res := map[int64]*proto.MarketPrice{}
	for _, p := range prices {
		res[int64(p.TypeID)] = proto.MarketPriceToProto(p)
	}
	return &proto.GetMarketPriceResponse{Result: successResult, Prices: res}, nil
}
