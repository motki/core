package client

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

// NewProduct creates a new Production Chain for the given type ID.
//
// If a production chain already exists for the given type ID, it will be returned.
func (c *GRPCClient) NewProduct(typeID int) (*model.Product, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewProductServiceClient(conn)
	res, err := service.NewProduct(
		context.Background(),
		&proto.NewProductRequest{
			Token:  &proto.Token{Identifier: c.token},
			TypeId: int64(typeID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	if res.Product == nil {
		return nil, errors.New("expected grpc response to contain product, got nil")
	}
	return proto.ProtoToProduct(res.Product), nil
}

// GetProduct attempts to load an existing production chain using its unique product ID.
func (c *GRPCClient) GetProduct(productID int) (*model.Product, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewProductServiceClient(conn)
	res, err := service.GetProduct(
		context.Background(),
		&proto.GetProductRequest{
			Token: &proto.Token{Identifier: c.token},
			Id:    int32(productID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	if res.Product == nil {
		return nil, errors.New("expected grpc response to contain product, got nil")
	}
	return proto.ProtoToProduct(res.Product), nil
}

// GetProducts gets all production chains for the current session's corporation.
func (c *GRPCClient) GetProducts() ([]*model.Product, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewProductServiceClient(conn)
	res, err := service.GetProducts(
		context.Background(),
		&proto.GetProductsRequest{Token: &proto.Token{Identifier: c.token}})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var prods []*model.Product
	for _, pr := range res.Product {
		prods = append(prods, proto.ProtoToProduct(pr))
	}
	return prods, nil
}

// SaveProduct attempts to save the given production chain to the backend database.
func (c *GRPCClient) SaveProduct(product *model.Product) error {
	if c.token == "" {
		return ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	service := proto.NewProductServiceClient(conn)
	res, err := service.SaveProduct(
		context.Background(),
		&proto.SaveProductRequest{
			Token:   &proto.Token{Identifier: c.token},
			Product: proto.ProductToProto(product)})
	if err != nil {
		return err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return errors.New(res.Result.Description)
	}
	return nil
}

// UpdateProductPrices updates all items in a production chain with the latest market sell price.
func (c *GRPCClient) UpdateProductPrices(product *model.Product) (*model.Product, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewProductServiceClient(conn)
	res, err := service.UpdateProductPrices(
		context.Background(),
		&proto.UpdateProductPricesRequest{
			Token:   &proto.Token{Identifier: c.token},
			Product: proto.ProductToProto(product),
		})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	if res.Product == nil {
		return nil, errors.New("expected grpc response to contain product, got nil")
	}
	return proto.ProtoToProduct(res.Product), nil
}
