package server

import (
	"golang.org/x/net/context"

	"github.com/motki/core/model"
	"github.com/motki/core/proto"
	"github.com/pkg/errors"
)

func (srv *grpcServer) GetLocation(ctx context.Context, req *proto.GetLocationRequest) (resp *proto.LocationResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.LocationResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	if req.Token == nil {
		return nil, errors.New("token cannot be empty")
	}
	ctx, _, err = srv.getAuthorizedContext(req.Token, model.RoleLogistics)
	if err != nil {
		return nil, err
	}
	loc, err := srv.model.GetLocation(ctx, int(req.LocationId))
	if err != nil {
		return nil, err
	}
	return &proto.LocationResponse{
		Result:   successResult,
		Location: proto.LocationToProto(loc),
	}, nil
}

func (srv *grpcServer) QueryLocations(ctx context.Context, req *proto.QueryLocationsRequest) (resp *proto.LocationsResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.LocationsResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	if req.Token == nil {
		return nil, errors.New("token cannot be empty")
	}
	ctx, _, err = srv.getAuthorizedContext(req.Token, model.RoleLogistics)
	if err != nil {
		return nil, err
	}
	locs, err := srv.model.QueryLocations(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	var res []*proto.Location
	for _, l := range locs {
		res = append(res, proto.LocationToProto(l))
	}
	return &proto.LocationsResponse{
		Result:   successResult,
		Location: res,
	}, nil
}
