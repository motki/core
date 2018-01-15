package server

import (
	"golang.org/x/net/context"

	"github.com/motki/core/proto"
)

func (srv *grpcServer) GetRegion(ctx context.Context, req *proto.GetRegionRequest) (resp *proto.GetRegionResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetRegionResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetRegion(int(req.RegionId))
	if err != nil {
		return nil, err
	}
	return &proto.GetRegionResponse{
		Result: successResult,
		Region: proto.RegionToProto(res),
	}, nil
}

func (srv *grpcServer) GetRegions(ctx context.Context, req *proto.GetRegionsRequest) (resp *proto.GetRegionsResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetRegionsResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetAllRegions()
	if err != nil {
		return nil, err
	}
	var results []*proto.Region
	for _, r := range res {
		results = append(results, proto.RegionToProto(r))
	}
	return &proto.GetRegionsResponse{
		Result: successResult,
		Region: results,
	}, nil
}

func (srv *grpcServer) GetConstellation(ctx context.Context, req *proto.GetConstellationRequest) (resp *proto.GetConstellationResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetConstellationResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetConstellation(int(req.ConstellationId))
	if err != nil {
		return nil, err
	}
	return &proto.GetConstellationResponse{
		Result:        successResult,
		Constellation: proto.ConstellationToProto(res),
	}, nil
}

func (srv *grpcServer) GetSystem(ctx context.Context, req *proto.GetSystemRequest) (resp *proto.GetSystemResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetSystemResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetSystem(int(req.SystemId))
	if err != nil {
		return nil, err
	}
	return &proto.GetSystemResponse{
		Result: successResult,
		System: proto.SystemToProto(res),
	}, nil
}

func (srv *grpcServer) GetRace(ctx context.Context, req *proto.GetRaceRequest) (resp *proto.GetRaceResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetRaceResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetRace(int(req.RaceId))
	if err != nil {
		return nil, err
	}
	return &proto.GetRaceResponse{
		Result: successResult,
		Race:   proto.RaceToProto(res),
	}, nil
}

func (srv *grpcServer) GetRaces(ctx context.Context, req *proto.GetRacesRequest) (resp *proto.GetRacesResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetRacesResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetRaces()
	if err != nil {
		return nil, err
	}
	var races []*proto.Race
	for _, r := range res {
		races = append(races, proto.RaceToProto(r))
	}
	return &proto.GetRacesResponse{
		Result: successResult,
		Race:   races,
	}, nil
}

func (srv *grpcServer) GetBloodline(ctx context.Context, req *proto.GetBloodlineRequest) (resp *proto.GetBloodlineResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetBloodlineResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetBloodline(int(req.BloodlineId))
	if err != nil {
		return nil, err
	}
	return &proto.GetBloodlineResponse{
		Result:    successResult,
		Bloodline: proto.BloodlineToProto(res),
	}, nil
}

func (srv *grpcServer) GetAncestry(ctx context.Context, req *proto.GetAncestryRequest) (resp *proto.GetAncestryResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetAncestryResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetAncestry(int(req.AncestryId))
	if err != nil {
		return nil, err
	}
	return &proto.GetAncestryResponse{
		Result:   successResult,
		Ancestry: proto.AncestryToProto(res),
	}, nil
}

func (srv *grpcServer) GetItemType(ctx context.Context, req *proto.GetItemTypeRequest) (resp *proto.GetItemTypeResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetItemTypeResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetItemType(int(req.TypeId))
	if err != nil {
		return nil, err
	}
	return &proto.GetItemTypeResponse{
		Result: successResult,
		Type:   proto.ItemTypeToProto(res),
	}, nil
}

func (srv *grpcServer) GetItemTypeDetail(ctx context.Context, req *proto.GetItemTypeDetailRequest) (resp *proto.GetItemTypeDetailResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetItemTypeDetailResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetItemTypeDetail(int(req.TypeId))
	if err != nil {
		return nil, err
	}
	return &proto.GetItemTypeDetailResponse{
		Result: successResult,
		Type:   proto.ItemTypeDetailToProto(res),
	}, nil
}

func (srv *grpcServer) QueryItemTypes(ctx context.Context, req *proto.QueryItemTypesRequest) (resp *proto.QueryItemTypesResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.QueryItemTypesResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	var cats []int
	for _, c := range req.CategoryId {
		cats = append(cats, int(c))
	}
	res, err := srv.evedb.QueryItemTypes(req.Query, cats...)
	if err != nil {
		return nil, err
	}
	var results []*proto.ItemType
	for _, r := range res {
		results = append(results, proto.ItemTypeToProto(r))
	}
	return &proto.QueryItemTypesResponse{
		Result: successResult,
		Types:  results,
	}, nil
}

func (srv *grpcServer) QueryItemTypeDetails(ctx context.Context, req *proto.QueryItemTypeDetailsRequest) (resp *proto.QueryItemTypeDetailsResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.QueryItemTypeDetailsResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	var cats []int
	for _, c := range req.CategoryId {
		cats = append(cats, int(c))
	}
	res, err := srv.evedb.QueryItemTypeDetails(req.Query, cats...)
	if err != nil {
		return nil, err
	}
	var results []*proto.ItemTypeDetail
	for _, r := range res {
		results = append(results, proto.ItemTypeDetailToProto(r))
	}
	return &proto.QueryItemTypeDetailsResponse{
		Result: successResult,
		Types:  results,
	}, nil
}

func (srv *grpcServer) GetMaterialSheet(ctx context.Context, req *proto.GetMaterialSheetRequest) (resp *proto.GetMaterialSheetResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.GetMaterialSheetResponse{
				Result: errorResult(err),
			}
			err = nil
		}
	}()
	res, err := srv.evedb.GetBlueprint(int(req.TypeId))
	if err != nil {
		return nil, err
	}
	return &proto.GetMaterialSheetResponse{
		Result:   successResult,
		MatSheet: proto.MatSheetToProto(res),
	}, nil
}
