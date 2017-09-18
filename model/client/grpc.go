package client

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/motki/motkid/log"
	"github.com/motki/motkid/model"
	"github.com/motki/motkid/model/proto"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"
)

var _ Client = &grpcClient{}

type grpcClient struct {
	serverAddr string
	token      string
	dialOpts   []grpc.DialOption
	logger     log.Logger
}

type dialerFunc func() (net.Conn, error)

func NewRemoteGRPC(serverAddr string, l log.Logger, tlsConf *tls.Config) (*grpcClient, error) {
	return &grpcClient{
		serverAddr: serverAddr,
		dialOpts:   []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(tlsConf))},
		logger:     l,
	}, nil
}

func NewLocalGRPC(lis *bufconn.Listener, l log.Logger) (*grpcClient, error) {
	cl := &grpcClient{logger: l}
	cl.dialOpts = append(cl.dialOpts, grpc.WithDialer(func(string, time.Duration) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
	return cl, nil
}

func (c *grpcClient) Authenticate(username, password string) (string, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	service := proto.NewAuthenticationServiceClient(conn)
	res, err := service.Authenticate(
		context.Background(),
		&proto.AuthenticateRequest{Username: username, Password: password})
	if err != nil {
		return "", err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return "", errors.New(res.Result.Description)
	}
	if res.Token == nil {
		return "", errors.New("expected token to be not empty, got nil")
	}
	c.token = res.Token.Identifier
	return res.Token.Identifier, nil
}

func (c *grpcClient) NewProduct(typeID int) (*model.Product, error) {
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

func (c *grpcClient) GetProduct(productID int) (*model.Product, error) {
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

func (c *grpcClient) GetProducts() ([]*model.Product, error) {
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

func (c *grpcClient) SaveProduct(product *model.Product) error {
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

func (c *grpcClient) UpdateProductPrices(product *model.Product) (*model.Product, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
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
			TypeId: []int64{int64(product.TypeID)},
		})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	price, ok := res.Prices[int64(product.TypeID)]
	if !ok {
		return nil, errors.Errorf("expected prices to contain price for typeID %d, got nil", product.TypeID)
	}
	product.MarketPrice = decimal.NewFromFloat(price.Average)
	return product, nil
}

func (c *grpcClient) CharacterForRole(role model.Role) (*model.Character, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInfoServiceClient(conn)
	var r proto.Role
	switch role {
	case model.RoleLogistics:
		r = proto.Role_LOGISTICS
	case model.RoleUser:
		r = proto.Role_USER
	default:
		//no op
	}
	res, err := service.GetCharacter(
		context.Background(),
		&proto.GetCharacterRequest{Token: &proto.Token{Identifier: c.token}, Role: r})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	if res.Character == nil {
		return nil, errors.New("expected grpc response to contain character, got nil")
	}
	return proto.ProtoToCharacter(res.Character), nil
}

func (c *grpcClient) GetCharacter(charID int) (*model.Character, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInfoServiceClient(conn)
	res, err := service.GetCharacter(
		context.Background(),
		&proto.GetCharacterRequest{Token: &proto.Token{Identifier: c.token}, CharacterId: int64(charID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	if res.Character == nil {
		return nil, errors.New("expected grpc response to contain character, got nil")
	}
	return proto.ProtoToCharacter(res.Character), nil
}

func (c *grpcClient) GetCorporation(corpID int) (*model.Corporation, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInfoServiceClient(conn)
	res, err := service.GetCorporation(
		context.Background(),
		&proto.GetCorporationRequest{Token: &proto.Token{Identifier: c.token}, CorporationId: int64(corpID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	corp := res.Corporation
	if corp == nil {
		return nil, errors.New("expected corporation in grpc response, got nil")
	}
	return proto.ProtoToCorporation(corp), nil
}

func (c *grpcClient) GetAlliance(allianceID int) (*model.Alliance, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewInfoServiceClient(conn)
	res, err := service.GetAlliance(
		context.Background(),
		&proto.GetAllianceRequest{Token: &proto.Token{Identifier: c.token}, AllianceId: int64(allianceID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	alliance := res.Alliance
	if alliance == nil {
		return nil, errors.New("expected alliance in grpc response, got nil")
	}
	return proto.ProtoToAlliance(alliance), nil
}
