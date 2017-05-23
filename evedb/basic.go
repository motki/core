package evedb

import "strings"

// A Race is a race in EVE.
type Race struct {
	ID               int
	Name             string
	Description      string
	IconFile         string
	ShortDescription string
}

type ItemType struct {
	ID int
	Name string
}

type Blueprint struct {
	*ItemType
	Materials []Material
}

type Material struct {
	ID int
	Name string
	Quantity int
}

func (r Race) Icon() string {
	// TODO: This is coupled to the web server, move it somewhere else
	if r.IconFile == "" {
		return "/images/Icons/items/7_64_15.png"
	}
	return strings.Replace(r.IconFile, "res:/ui/texture/icons", "/images/Icons/items", 1)
}

// GetRaces fetches all Races from the database.
func (e *EveDB) GetRaces() ([]*Race, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rs, err := c.Query(
		`SELECT
			  race."raceID"
			, race."raceName"
			, COALESCE(race."description", '')
			, COALESCE(icon."iconFile", '')
			, COALESCE(race."shortDescription", '')
			FROM evesde."chrRaces" race
			LEFT JOIN evesde."eveIcons" icon ON race."iconID" = icon."iconID"`)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	res := []*Race{}
	for rs.Next() {
		r := &Race{}
		err := rs.Scan(&r.ID, &r.Name, &r.Description, &r.IconFile, &r.ShortDescription)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

// GetItemType fetches a specific ItemType from the database.
func (e *EveDB) GetItemType(typeID int) (*ItemType, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	r := c.QueryRow(
		`SELECT
			  type."typeID" as ID,
			  type."typeName" as typeName
			FROM evesde."invTypes" type
			WHERE type."typeID" = $1 `, typeID)
	it := &ItemType{}
	err = r.Scan(&it.ID, &it.Name)
	if err != nil {
		return nil, err
	}
	return it, nil
}

// GetBlueprint fetches a Blueprint from the database.
func (e *EveDB) GetBlueprint(typeID int) (*Blueprint, error) {
	it, err := e.GetItemType(typeID)
	if err != nil {
		return nil, err
	}
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rs, err := c.Query(
		`SELECT
			  type."typeName" as typeName,
			  type."typeID" as materialID,
			  mats."quantity"
			FROM evesde."invTypeMaterials" mats
			INNER JOIN evesde."invTypes" type ON type."typeID" = mats."materialTypeID"
			WHERE mats."typeID" = $1 `, it.ID)

	if err != nil {
		return nil, err
	}
	defer rs.Close()
	res := []Material{}
	for rs.Next() {
		r := Material{}
		err := rs.Scan(&r.Name, &r.ID, &r.Quantity)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return &Blueprint{ItemType: it, Materials: res}, nil
}

// GetBlueprints is a utility function to retrieve multiple Blueprints.
func (e *EveDB) GetBlueprints(typeIDs ...int) ([]*Blueprint, error) {
	res := []*Blueprint{}
	for _, id := range typeIDs {
		bp, err := e.GetBlueprint(id)
		if err != nil {
			return nil, err
		}
		res = append(res, bp)
	}
	return res, nil
}