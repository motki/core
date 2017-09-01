package evedb

type System struct {
	SystemID        int
	Name            string
	RegionID        int
	ConstellationID int
	Security        float64
}

type Constellation struct {
	ConstellationID int
	Name            string
	RegionID        int
}

type Region struct {
	RegionID int
	Name     string
}

func (e *EveDB) GetSystem(id int) (*System, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
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
	defer c.Close()
	row := c.QueryRow(
		`SELECT
			, s."constellationID"
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
	defer c.Close()
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
	defer c.Close()
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
