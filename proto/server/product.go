package server

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

func productResponse(product *model.Product) *proto.ProductResponse {
	return &proto.ProductResponse{
		Result:  successResult,
		Product: proto.ProductToProto(product),
	}
}

func (srv *grpcServer) GetProduct(ctx context.Context, req *proto.GetProductRequest) (resp *proto.ProductResponse, err error) {
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
	_, err = srv.model.GetCorporationAuthorization(char.CorporationID)
	if err != nil {
		return nil, err
	}
	prod, err := srv.model.GetProduct(corp.CorporationID, int(req.Id))
	if err != nil {
		return nil, err
	}
	return productResponse(prod), nil
}

func (srv *grpcServer) GetProducts(ctx context.Context, req *proto.GetProductsRequest) (resp *proto.ProductsResponse, err error) {
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

func (srv *grpcServer) NewProduct(ctx context.Context, req *proto.NewProductRequest) (resp *proto.ProductResponse, err error) {
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

func (srv *grpcServer) SaveProduct(ctx context.Context, req *proto.SaveProductRequest) (resp *proto.ProductResponse, err error) {
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
		// TODO: always get the existing product to make sure it belongs to the corp. inefficient.
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

func (srv *grpcServer) UpdateProductPrices(ctx context.Context, req *proto.UpdateProductPricesRequest) (resp *proto.ProductResponse, err error) {
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
	_, err = srv.model.GetAuthorization(user, model.RoleLogistics)
	if err != nil {
		return nil, err
	}
	prod := proto.ProtoToProduct(req.Product)
	err = srv.model.UpdateProductMarketPrices(prod, prod.MarketRegionID)
	if err != nil {
		return nil, err
	}
	return productResponse(prod), nil
}
