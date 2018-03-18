package client

import (
	"strconv"

	"time"

	"github.com/motki/core/cache"
	"github.com/motki/core/evedb"
	"github.com/motki/core/model"
	"github.com/pkg/errors"
)

// Cache time-to-live for static data.
const cacheTTL = 600 * time.Second

// cachingGRPCClient wraps a GRPC client and provides short-lived, in-memory
// caching for static data retrieved from a remote GRPC server.
type cachingGRPCClient struct {
	*GRPCClient

	cache *cache.Bucket
}

func cacheKey(prefix string, id int) string {
	return prefix + strconv.Itoa(id)
}

// GetRegion returns information about the given region ID.
func (c *cachingGRPCClient) GetRegion(regionID int) (*evedb.Region, error) {
	v, err := c.cache.Memoize(cacheKey("region:", regionID), func() (cache.Value, error) {
		return c.GRPCClient.GetRegion(regionID)
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.(*evedb.Region); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *evedb.Region from cache, got %T", v)
}

// GetRegions returns a slice containing information about all regions in the EVE universe.
func (c *cachingGRPCClient) GetRegions() ([]*evedb.Region, error) {
	v, err := c.cache.Memoize("regions", func() (cache.Value, error) {
		return c.GRPCClient.GetRegions()
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.([]*evedb.Region); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected []*evedb.Region from cache, got %T", v)
}

// GetSystem returns information about the given system ID.
func (c *cachingGRPCClient) GetSystem(systemID int) (*evedb.System, error) {
	v, err := c.cache.Memoize(cacheKey("system:", systemID), func() (cache.Value, error) {
		return c.GRPCClient.GetSystem(systemID)
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.(*evedb.System); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *evedb.System from cache, got %T", v)
}

// GetConstellation returns information about the given constellation ID.
func (c *cachingGRPCClient) GetConstellation(constellationID int) (*evedb.Constellation, error) {
	v, err := c.cache.Memoize(cacheKey("constellation:", constellationID), func() (cache.Value, error) {
		return c.GRPCClient.GetConstellation(constellationID)
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.(*evedb.Constellation); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *evedb.Constellation from cache, got %T", v)
}

// GetRace returns information about the given race ID.
func (c *cachingGRPCClient) GetRace(raceID int) (*evedb.Race, error) {
	v, err := c.cache.Memoize(cacheKey("race:", raceID), func() (cache.Value, error) {
		return c.GRPCClient.GetRace(raceID)
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.(*evedb.Race); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *evedb.Race from cache, got %T", v)
}

// GetRaces returns information about all races in the EVE universe.
func (c *cachingGRPCClient) GetRaces() ([]*evedb.Race, error) {
	v, err := c.cache.Memoize("races", func() (cache.Value, error) {
		return c.GRPCClient.GetRaces()
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.([]*evedb.Race); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected []*evedb.Race from cache, got %T", v)
}

// GetBloodline returns information about the given bloodline ID.
func (c *cachingGRPCClient) GetBloodline(bloodlineID int) (*evedb.Bloodline, error) {
	v, err := c.cache.Memoize(cacheKey("bloodline:", bloodlineID), func() (cache.Value, error) {
		return c.GRPCClient.GetBloodline(bloodlineID)
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.(*evedb.Bloodline); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *evedb.Bloodline from cache, got %T", v)
}

// GetAncestry returns information about the given ancestry ID.
func (c *cachingGRPCClient) GetAncestry(ancestryID int) (*evedb.Ancestry, error) {
	v, err := c.cache.Memoize(cacheKey("ancestry:", ancestryID), func() (cache.Value, error) {
		return c.GRPCClient.GetAncestry(ancestryID)
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.(*evedb.Ancestry); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *evedb.Ancestry from cache, got %T", v)
}

// GetItemType returns information about the given type ID.
func (c *cachingGRPCClient) GetItemType(typeID int) (*evedb.ItemType, error) {
	v, err := c.cache.Memoize(cacheKey("item:", typeID), func() (cache.Value, error) {
		return c.GRPCClient.GetItemType(typeID)
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.(*evedb.ItemType); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *evedb.ItemType from cache, got %T", v)
}

// GetItemTypeDetail returns detailed information about the given type ID.
func (c *cachingGRPCClient) GetItemTypeDetail(typeID int) (*evedb.ItemTypeDetail, error) {
	v, err := c.cache.Memoize(cacheKey("item_detail:", typeID), func() (cache.Value, error) {
		return c.GRPCClient.GetItemTypeDetail(typeID)
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.(*evedb.ItemTypeDetail); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *evedb.ItemTypeDetail from cache, got %T", v)
}

// GetMaterialSheet returns manufacturing information about the given type ID.
func (c *cachingGRPCClient) GetMaterialSheet(typeID int) (*evedb.MaterialSheet, error) {
	v, err := c.cache.Memoize(cacheKey("mat_sheet:", typeID), func() (cache.Value, error) {
		return c.GRPCClient.GetMaterialSheet(typeID)
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.(*evedb.MaterialSheet); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *evedb.MaterialSheet from cache, got %T", v)
}

// GetLocation returns location information for a denormalized location ID.
func (c *cachingGRPCClient) GetLocation(locationID int) (*model.Location, error) {
	v, err := c.cache.Memoize(cacheKey("location:", locationID), func() (cache.Value, error) {
		return c.GRPCClient.GetLocation(locationID)
	})
	if err != nil {
		return nil, err
	}
	if a, ok := v.(*model.Location); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *model.Location from cache, got %T", v)
}
