package evedb

import (
	"strconv"
	"strings"
)

// An ItemType is a type of item in EVE.
type ItemType struct {
	ID          int
	Name        string
	Description string
}

// A Blueprint describes what is necessary to build an item.
type Blueprint struct {
	*ItemType
	Materials []Material
}

// A Material is a type and quantity of an item used for manufacturing.
type Material struct {
	*ItemType
	Quantity int
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
			  type."typeID"
			, type."typeName"
			, type."description"
			FROM evesde."invTypes" type
			WHERE type."typeID" = $1`, typeID)
	it := &ItemType{}
	err = r.Scan(&it.ID, &it.Name, &it.Description)
	if err != nil {
		return nil, err
	}
	return it, nil
}

// QueryItemTypes returns a list of matching items given the query.
func (e *EveDB) QueryItemTypes(query string, catIDs ...int) ([]*ItemType, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	if len(catIDs) == 0 {
		// Default to Modules, Ships, Drones, and Charges
		catIDs = []int{6, 7, 8, 18}
	}
	cats := []string{}
	for _, id := range catIDs {
		cats = append(cats, strconv.Itoa(id))
	}
	rs, err := c.Query(
		`SELECT
			  type."typeID"
			, type."typeName"
			FROM evesde."invTypes" type
			  JOIN evesde."invGroups" grp
			    ON type."groupID" = grp."groupID" AND grp."categoryID" = ANY($1::INTEGER[])
			WHERE  type."published" = TRUE
			  AND type."typeName" ILIKE '%' || $2 || '%'
			LIMIT 20`, "{"+strings.Join(cats, ",")+"}", query)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	res := []*ItemType{}
	for rs.Next() {
		r := &ItemType{}
		err := rs.Scan(&r.ID, &r.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
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
			  typ."typeName" as typeName
			, typ."typeID" as materialID
			, mats."quantity"
			FROM evesde."invTypeMaterials" mats
			INNER JOIN evesde."invTypes" typ ON typ."typeID" = mats."materialTypeID"
			WHERE mats."typeID" = $1`, it.ID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	res := []Material{}
	for rs.Next() {
		r := Material{ItemType: &ItemType{}}
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
