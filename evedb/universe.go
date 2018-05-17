package evedb

type System struct {
	SystemID        int     `json:"system_id"`
	Name            string  `json:"name"`
	RegionID        int     `json:"region_id"`
	ConstellationID int     `json:"constellation_id"`
	Security        float64 `json:"security"`
}

type Constellation struct {
	ConstellationID int    `json:"constellation_id"`
	Name            string `json:"name"`
	RegionID        int    `json:"region_id"`
}

type Region struct {
	RegionID int    `json:"region_id"`
	Name     string `json:"name"`
}

type Station struct {
	StationID       int    `json:"station_id"`
	StationTypeID   int    `json:"station_type_id"`
	CorporationID   int    `json:"corporation_id"`
	SystemID        int    `json:"system_id"`
	ConstellationID int    `json:"constellation_id"`
	RegionID        int    `json:"region_id"`
	Name            string `json:"name"`
}

func (e *EveDB) GetSystem(id int) (*System, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer e.pool.Release(c)
	row := c.QueryRow(
		`SELECT
			  s."solarSystemID"
			, s."constellationID"
			, s."regionID"
			, s."solarSystemName"
			, s."security"
			FROM evesde."mapSolarSystems" s
			WHERE s."solarSystemID" = $1`, id)
	r := System{}
	err = row.Scan(&r.SystemID, &r.ConstellationID, &r.RegionID, &r.Name, &r.Security)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (e *EveDB) GetConstellation(id int) (*Constellation, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer e.pool.Release(c)
	row := c.QueryRow(
		`SELECT
			  s."constellationID"
			, s."regionID"
			, s."constellationName"
			FROM evesde."mapConstellations" s
			WHERE s."constellationID" = $1`, id)
	r := Constellation{}
	err = row.Scan(&r.ConstellationID, &r.RegionID, &r.Name)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (e *EveDB) GetRegion(id int) (*Region, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer e.pool.Release(c)
	row := c.QueryRow(
		`SELECT
			  s."regionID"
			, s."regionName"
			FROM evesde."mapRegions" s
			WHERE s."regionID" = $1`, id)
	r := Region{}
	err = row.Scan(&r.RegionID, &r.Name)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (e *EveDB) GetAllRegions() ([]*Region, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer e.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  s."regionID"
			, s."regionName"
			FROM evesde."mapRegions" s`)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*Region
	for rs.Next() {
		r := Region{}
		err = rs.Scan(&r.RegionID, &r.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, &r)
	}
	return res, nil
}

func (e *EveDB) GetStation(stationID int) (*Station, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer e.pool.Release(c)
	r := c.QueryRow(
		`SELECT s."stationID"
     			, s."stationTypeID"
     			, s."stationName"
     			, s."solarSystemID"
     			, s."constellationID"
     			, s."regionID"
			, s."corporationID"
     			FROM evesde."staStations" s
     			  WHERE s."stationID" = $1`, stationID)
	s := &Station{}
	err = r.Scan(&s.StationID, &s.StationTypeID, &s.Name, &s.SystemID, &s.ConstellationID, &s.RegionID, &s.CorporationID)
	if err != nil {
		return nil, err
	}
	return s, nil
}
