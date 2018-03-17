package model

import (
	"github.com/motki/core/eveapi"
	"github.com/motki/core/evedb"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Location describes a station, structure, or solar system in the EVE universe.
//
// This is a basic abstraction over the loosely defined "location ID" found in various
// API responses. A Location may represent an NPC station, a player-owned Citadel, or
// simply a solar system.
//
// Any Location will contain, at a minimum, System, Constellation, and Region
// info. Station and Structure may be nil if the location is only as specific as a
// solar system. Otherwise, Station OR Structure will be populated, but never both.
type Location struct {
	// The original LocationID this location represents.
	LocationID int

	// The solar system this location is found in.
	System *evedb.System

	// The constellation this location is found in.
	Constellation *evedb.Constellation

	// The region this location is found in.
	Region *evedb.Region

	// The NPC station at this location. May be nil.
	Station *evedb.Station

	// The player-owned structure at this location. May be nil.
	Structure *eveapi.Structure

	// Prevent other packages from creating this type.
	noexport struct{}
}

func (l Location) String() string {
	if l.IsCitadel() {
		return l.Structure.Name
	}
	if l.IsStation() {
		return l.Station.Name
	}
	return l.System.Name
}

// IsStation returns true if the location is a NPC station.
//
// Use this method to determine if a location represents a NPC station.
func (l Location) IsStation() bool {
	return l.Station != nil
}

// IsCitadel returns true if the location is a player-controlled citadel.
//
// Use this method to determine if a location represents a player-controlled structure.
func (l Location) IsCitadel() bool {
	return l.Structure != nil
}

// IsSystem returns true if the location does not contain station or citadel information.
//
// Use this method to determine if a location is neither a station nor structure. In other
// words, this method returns true if the location is strictly just a solar system.
func (l Location) IsSystem() bool {
	return l.Station == nil && l.Structure == nil
}

// ParentID returns the StationID, StructureID, or SystemID that this location exists in.
func (l Location) ParentID() int {
	if l.IsStation() {
		return l.Station.StationID
	} else if l.IsCitadel() {
		return int(l.Structure.StructureID)
	}
	return l.System.SystemID
}

// GetLocation attempts to resolve the given location.
func (m *Manager) GetLocation(ctx context.Context, locationID int) (*Location, error) {
	// Magic numbers here are sourced from:
	// - http://eveonline-third-party-documentation.readthedocs.io/en/latest/xmlapi/character/char_assetlist.html
	// - https://oldforums.eveonline.com/?a=topic&threadID=667487
	const offsetOfficeIDToStationID = 6000001
	const legacyOutpostStart = 60014861
	const legacyOutpostEndInclusive = 60014928
	loc := &Location{}
	var corpID int
	if c, ok := authContextFromContext(ctx); ok {
		corpID = c.CorporationID()
	}
	var err error
	switch {
	case locationID < 60000000:
		// locationID is a SystemID.
		loc.System, err = m.evedb.GetSystem(locationID)
		if err != nil {
			return nil, err
		}

	case locationID < 61000000:
		// locationID is a Station or legacy outpost.
		if locationID >= legacyOutpostStart && locationID <= legacyOutpostEndInclusive {
			// Conquerable outpost pre-dating player outposts. Not yet supported.
			return nil, errors.Errorf("unable to determine details for locationID %d, conquerable outposts are not supported", locationID)
		}
		// Not a legacy outpost, must be a station.
		loc.Station, err = m.evedb.GetStation(locationID)

	case locationID < 66000000:
		// locationID is a conquerable outpost. Not yet supported.
		return nil, errors.Errorf("unable to determine details for locationID %d, conquerable outposts are not supported", locationID)

	case locationID < 67000000:
		// locationID is a rented office.
		loc.Station, err = m.evedb.GetStation(locationID - offsetOfficeIDToStationID)

	default:
		// locationID might be a citadel.
		s, err := m.GetStructure(ctx, locationID)
		if err == nil {
			loc.Structure = s
			break
		}
		// locationID is in a container somewhere.
		if corpID != 0 {
			// Corporation is opted-in, we can query for asset information.
			if ca, err := m.GetCorporationAsset(ctx, corpID, locationID); err == nil {
				// Found an asset with the given locationID, call GetLocation on the asset's location.
				if loc, err = m.GetLocation(ctx, ca.LocationID); err != nil {
					return nil, errors.Wrapf(err, "unable to determine details for locationID %d", locationID)
				}
				break
			} else if err != ErrCorpNotRegistered {
				return nil, errors.Wrapf(err, "unable to determine details for locationID %d", locationID)
			}
		}

	}

	// Second round of resolution; ensure that loc.System gets populated.
	switch {
	case loc.System != nil:
		// do nothing

	case loc.Structure != nil:
		if loc.System, err = m.evedb.GetSystem(int(loc.Structure.SystemID)); err != nil {
			return nil, err
		}

	case loc.Station != nil:
		if loc.System, err = m.evedb.GetSystem(loc.Station.SystemID); err != nil {
			return nil, err
		}

	default:
		return nil, errors.Errorf("unable to determine details for locationID %d", locationID)
	}

	// Finally, populate the Constellation and Region information.
	if loc.Constellation, err = m.evedb.GetConstellation(loc.System.ConstellationID); err != nil {
		return nil, err
	}
	if loc.Region, err = m.evedb.GetRegion(loc.System.RegionID); err != nil {
		return nil, err
	}
	loc.LocationID = locationID
	return loc, nil
}

func (m *Manager) QueryLocations(ctx context.Context, query string) ([]*Location, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	var corpID int
	if c, ok := authContextFromContext(ctx); ok {
		corpID = c.CorporationID()
	}
	r, err := c.Query(`SELECT s.structure_id
FROM app.structures s
WHERE s.name ILIKE '%' || $1 || '%'
  AND s.corporation_id = $2

UNION ALL

SELECT s."stationID"
FROM evesde."staStations" s
WHERE s."stationName" ILIKE '%' || $1 || '%'

UNION ALL

SELECT s."solarSystemID"
FROM evesde."mapSolarSystems" s
WHERE s."solarSystemName" ILIKE '%' || $1 || '%'`, query, corpID)
	if err != nil {
		return nil, err
	}
	var res []*Location
	for r.Next() {
		var i int
		err = r.Scan(&i)
		if err != nil {
			return nil, err
		}
		// TODO: expensive call in a loop
		loc, err := m.GetLocation(ctx, i)
		if err != nil {
			return nil, err
		}
		res = append(res, loc)
	}
	return res, nil
}
