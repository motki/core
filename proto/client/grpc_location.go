package client

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
)

// LocationClient provides location information using denormalized location IDs.
type LocationClient struct {
	// This type must be initialized using the package-level New function.

	*bootstrap
}

// GetLocation returns the given Location using the denormalized location ID.
//
// If the user's corporation has opted-in, asset and structure information is used to
// enhance the results.
func (c *LocationClient) GetLocation(locationID int) (*model.Location, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewLocationServiceClient(conn)
	res, err := service.GetLocation(
		context.Background(),
		&proto.GetLocationRequest{Token: &proto.Token{Identifier: c.token}, LocationId: int64(locationID)})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	return proto.ProtoToLocation(res.Location), nil
}

// QueryLocation return locations that match the input query.
//
// If the user's corporation has opted-in, asset and structure information is used to
// enhance the results.
func (c *LocationClient) QueryLocations(query string) ([]*model.Location, error) {
	if c.token == "" {
		return nil, ErrNotAuthenticated
	}
	conn, err := grpc.Dial(c.serverAddr, c.dialOpts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	service := proto.NewLocationServiceClient(conn)
	res, err := service.QueryLocations(
		context.Background(),
		&proto.QueryLocationsRequest{Token: &proto.Token{Identifier: c.token}, Query: query})
	if err != nil {
		return nil, err
	}
	if res.Result.Status == proto.Status_FAILURE {
		return nil, errors.New(res.Result.Description)
	}
	var locs []*model.Location
	for _, k := range res.Location {
		locs = append(locs, proto.ProtoToLocation(k))
	}
	return locs, nil
}
