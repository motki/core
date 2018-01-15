package client

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/core/evedb"
	"github.com/motki/core/proto"
)

// GetRegion returns information about the given region ID.
func (c *GRPCClient) GetRegion(regionID int) (*evedb.Region, error) {
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

// GetRegions returns a slice containing information about all regions in the EVE universe.
func (c *GRPCClient) GetRegions() ([]*evedb.Region, error) {
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

// GetSystem returns information about the given system ID.
func (c *GRPCClient) GetSystem(systemID int) (*evedb.System, error) {
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

// GetConstellation returns information about the given constellation ID.
func (c *GRPCClient) GetConstellation(constellationID int) (*evedb.Constellation, error) {
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

// GetRace returns information about the given race ID.
func (c *GRPCClient) GetRace(raceID int) (*evedb.Race, error) {
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

// GetRaces returns information about all races in the EVE universe.
func (c *GRPCClient) GetRaces() ([]*evedb.Race, error) {
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

// GetBloodline returns information about the given bloodline ID.
func (c *GRPCClient) GetBloodline(bloodlineID int) (*evedb.Bloodline, error) {
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

// GetAncestry returns information about the given ancestry ID.
func (c *GRPCClient) GetAncestry(ancestryID int) (*evedb.Ancestry, error) {
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
