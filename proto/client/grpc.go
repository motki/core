package client

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/motki/motki/log"
	"github.com/motki/motki/model"
	"github.com/motki/motki/proto"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/motki/motki/evedb"
	"github.com/motki/motki/proto/server"
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
		if res.Result.Description == server.ErrBadCredentials.Error() {
			return "", ErrBadCredentials
		}
		return "", errors.New(res.Result.Description)
	}
	if res.Token == nil || res.Token.Identifier == "" {
		return "", errors.New("expected token to be not empty, got nothing")
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
	idMap := map[int64]struct{}{}
	var addTypeID func(*model.Product)
	addTypeID = func(p *model.Product) {
		idMap[int64(p.TypeID)] = struct{}{}
		for _, sp := range p.Materials {
			addTypeID(sp)
		}
	}
	addTypeID(product)
	var typeIDs []int64
	for k := range idMap {
		typeIDs = append(typeIDs, k)
	}
	service := proto.NewMarketPriceServiceClient(conn)
	res, err := service.GetMarketPrice(
		context.Background(),
		&proto.GetMarketPriceRequest{
			Token:  &proto.Token{Identifier: c.token},
			TypeId: typeIDs,
		})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var updatePrice func(*model.Product) error
	updatePrice = func(p *model.Product) error {
		price, ok := res.Prices[int64(p.TypeID)]
		if !ok {
			return errors.Errorf("expected prices to contain price for typeID %d, got nil", product.TypeID)
		}
		p.MarketPrice = decimal.NewFromFloat(price.Average)
		for _, m := range p.Materials {
			if err := updatePrice(m); err != nil {
				return err
			}
		}
		return nil
	}
	err = updatePrice(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (c *grpcClient) GetMarketPrices(typeID int, typeIDs ...int) ([]*model.MarketPrice, error) {
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

func (c *grpcClient) GetMarketPrice(typeID int) (*model.MarketPrice, error) {
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

func (c *grpcClient) GetCorpBlueprints() ([]*model.Blueprint, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewCorporationServiceClient(conn)
	res, err := service.GetCorpBlueprints(
		context.Background(),
		&proto.GetCorpBlueprintsRequest{
			Token: &proto.Token{Identifier: c.token},
		})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var bps []*model.Blueprint
	for _, bp := range res.Blueprint {
		bps = append(bps, proto.ProtoToBlueprint(bp))
	}
	return bps, nil
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

func (c *grpcClient) GetRegion(regionID int) (*evedb.Region, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetRegion(
		context.Background(),
		&proto.GetRegionRequest{RegionId: int64(regionID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	pres := res.Region
	if pres == nil {
		return nil, errors.New("expected region in grpc response, got nil")
	}
	return proto.ProtoToRegion(pres), nil
}

func (c *grpcClient) GetRegions() ([]*evedb.Region, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetRegions(
		context.Background(),
		&proto.GetRegionsRequest{})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var results []*evedb.Region
	for _, pr := range res.Region {
		results = append(results, proto.ProtoToRegion(pr))
	}
	return results, nil
}

func (c *grpcClient) GetSystem(systemID int) (*evedb.System, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetSystem(
		context.Background(),
		&proto.GetSystemRequest{SystemId: int64(systemID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	pres := res.System
	if pres == nil {
		return nil, errors.New("expected system in grpc response, got nil")
	}
	return proto.ProtoToSystem(pres), nil
}

func (c *grpcClient) GetConstellation(constellationID int) (*evedb.Constellation, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetConstellation(
		context.Background(),
		&proto.GetConstellationRequest{ConstellationId: int64(constellationID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	pres := res.Constellation
	if pres == nil {
		return nil, errors.New("expected constellation in grpc response, got nil")
	}
	return proto.ProtoToConstellation(pres), nil
}

func (c *grpcClient) GetRace(raceID int) (*evedb.Race, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetRace(
		context.Background(),
		&proto.GetRaceRequest{RaceId: int64(raceID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	pres := res.Race
	if pres == nil {
		return nil, errors.New("expected race in grpc response, got nil")
	}
	return proto.ProtoToRace(pres), nil
}

func (c *grpcClient) GetRaces() ([]*evedb.Race, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetRaces(
		context.Background(),
		&proto.GetRacesRequest{})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var results []*evedb.Race
	for _, pr := range res.Race {
		results = append(results, proto.ProtoToRace(pr))
	}
	return results, nil
}

func (c *grpcClient) GetBloodline(bloodlineID int) (*evedb.Bloodline, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetBloodline(
		context.Background(),
		&proto.GetBloodlineRequest{BloodlineId: int64(bloodlineID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	pres := res.Bloodline
	if pres == nil {
		return nil, errors.New("expected bloodline in grpc response, got nil")
	}
	return proto.ProtoToBloodline(pres), nil
}

func (c *grpcClient) GetAncestry(ancestryID int) (*evedb.Ancestry, error) {
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewEveDBServiceClient(conn)
	res, err := service.GetAncestry(
		context.Background(),
		&proto.GetAncestryRequest{AncestryId: int64(ancestryID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	pres := res.Ancestry
	if pres == nil {
		return nil, errors.New("expected ancestry in grpc response, got nil")
	}
	return proto.ProtoToAncestry(pres), nil
}

func (c *grpcClient) GetItemType(typeID int) (*evedb.ItemType, error) {
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

func (c *grpcClient) GetItemTypeDetail(typeID int) (*evedb.ItemTypeDetail, error) {
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

func (c *grpcClient) QueryItemTypes(query string, catIDs ...int) ([]*evedb.ItemType, error) {
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

func (c *grpcClient) QueryItemTypeDetails(query string, catIDs ...int) ([]*evedb.ItemTypeDetail, error) {
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

func (c *grpcClient) GetMaterialSheet(typeID int) (*evedb.MaterialSheet, error) {
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
