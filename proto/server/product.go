package server

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/motki/motki/model"
	"github.com/motki/motki/proto"
)

func productResponse(product *model.Product) *proto.ProductResponse {
	return &proto.ProductResponse{
		Result:  successResult,
		Product: proto.ProductToProto(product),
	}
}

func (srv *GRPCServer) GetProduct(ctx context.Context, req *proto.GetProductRequest) (resp *proto.ProductResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.ProductResponse{
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
	char, err := srv.model.GetCharacter(a.CharacterID)
	if err != nil {
		return nil, err
	}
	corp, err := srv.model.GetCorporation(char.CorporationID)
	if err != nil {
		return nil, err
	}
	prod, err := srv.model.GetProduct(corp.CorporationID, int(req.Id))
	if err != nil {
		return nil, err
	}
	return productResponse(prod), nil
}

func (srv *GRPCServer) GetProducts(ctx context.Context, req *proto.GetProductsRequest) (resp *proto.ProductsResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.ProductsResponse{
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
	char, err := srv.model.GetCharacter(a.CharacterID)
	if err != nil {
		return nil, err
	}
	corp, err := srv.model.GetCorporation(char.CorporationID)
	if err != nil {
		return nil, err
	}
	prods, err := srv.model.GetAllProducts(corp.CorporationID)
	if err != nil {
		return nil, err
	}
	resp = &proto.ProductsResponse{Result: successResult}
	for _, p := range prods {
		resp.Product = append(resp.Product, proto.ProductToProto(p))
	}
	return resp, nil
}

func (srv *GRPCServer) NewProduct(ctx context.Context, req *proto.NewProductRequest) (resp *proto.ProductResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.ProductResponse{
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
	char, err := srv.model.GetCharacter(a.CharacterID)
	if err != nil {
		return nil, err
	}
	corp, err := srv.model.GetCorporation(char.CorporationID)
	if err != nil {
		return nil, err
	}
	prod, err := srv.model.NewProduct(corp.CorporationID, int(req.TypeId))
	if err != nil {
		return nil, err
	}
	return productResponse(prod), nil
}

func (srv *GRPCServer) SaveProduct(ctx context.Context, req *proto.SaveProductRequest) (resp *proto.ProductResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.ProductResponse{
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
	char, err := srv.model.GetCharacter(a.CharacterID)
	if err != nil {
		return nil, err
	}
	corp, err := srv.model.GetCorporation(char.CorporationID)
	if err != nil {
		return nil, err
	}
	prod := proto.ProtoToProduct(req.Product)
	if prod.ProductID != 0 {
		_, err := srv.model.GetProduct(corp.CorporationID, prod.ProductID)
		if err != nil {
			return nil, err
		}
	}
	var setCorpID func(p *model.Product)
	setCorpID = func(p *model.Product) {
		p.CorporationID = corp.CorporationID
		for _, pr := range p.Materials {
			setCorpID(pr)
		}
	}
	setCorpID(prod)
	err = srv.model.SaveProduct(prod)
	if err != nil {
		return nil, err
	}
	return productResponse(prod), nil
}

func (srv *GRPCServer) GetMarketPrice(ctx context.Context, req *proto.GetMarketPriceRequest) (resp *proto.GetMarketPriceResponse, err error) {
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
		avg, _ := p.Avg.Float64()
		base, _ := p.Base.Float64()
		res[int64(p.TypeID)] = &proto.MarketPrice{
			TypeId:  int64(p.TypeID),
			Average: avg,
			Base:    base,
		}
	}
	return &proto.GetMarketPriceResponse{Result: successResult, Prices: res}, nil
}
