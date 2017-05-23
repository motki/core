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

type Blueprint struct {
	ID int
	Materials []Material
}

type Material struct {
	ID int
	Name string
	MaterialID int
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

// GetRaces fetches all Races from the database.
func (e *EveDB) GetBlueprint(/* typeID int */) (*Blueprint, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rs, err := c.Query(
		`SELECT
			  mats."typeID" as typeID,
			  type."typeName" as typeName,
			  type."typeID" as materialID,
			  mats."quantity"
			FROM evesde."invTypeMaterials" mats
			INNER JOIN evesde."invTypes" type ON type."typeID" = mats."materialTypeID"
			WHERE mats."typeID" IN(16242) `) // TODO: Pass in the typeID

	if err != nil {
		return nil, err
	}
	defer rs.Close()
	res := []Material{}
	for rs.Next() {
		r := Material{}
		err := rs.Scan(&r.ID, &r.Name, &r.MaterialID, &r.Quantity)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return &Blueprint{Materials: res}, nil
}
