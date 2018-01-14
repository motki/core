package evedb

import (
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// An ItemType is a type of item in EVE.
type ItemType struct {
	ID          int
	Name        string
	Description string
}

type ItemTypeDetail struct {
	*ItemType

	GroupID   int
	GroupName string

	CategoryID   int
	CategoryName string

	Mass        decimal.Decimal
	Volume      decimal.Decimal
	Capacity    decimal.Decimal
	PortionSize int
	BasePrice   decimal.Decimal

	ParentTypeID      int
	BlueprintID       int
	DerivativeTypeIDs []int
}

const baseQueryItemType = `SELECT
  type."typeID"
, type."typeName"
, COALESCE(type."description", '')
FROM evesde."invTypes" type
`

// GetItemType fetches a specific ItemType from the database.
func (e *EveDB) GetItemType(typeID int) (*ItemType, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer e.pool.Release(c)
	r := c.QueryRow(
		baseQueryItemType+`WHERE type."typeID" = $1 AND type."published" = TRUE`, typeID)
	it := &ItemType{}
	err = r.Scan(&it.ID, &it.Name, &it.Description)
	if err != nil {
		return nil, err
	}
	return it, nil
}

// InterestingItemCategories contains a list of category IDs generally considered interesting.
//
// It doesn't contain the Blueprints category ID.
var InterestingItemCategories = []int{
	2, 4, 5, 6, 7, 8, 16, 17, 18,
	20, 22, 23, 24, 25, 30, 32, 34, 35, 39,
	40, 41, 42, 43, 46, 63, 65, 66, 87}

// InterestingItemCategoriesAndBlueprints contains a list of all published category IDs.
var InterestingItemCategoriesAndBlueprints = []int{
	2, 4, 5, 6, 7, 8, 9, 16, 17, 18,
	20, 22, 23, 24, 25, 30, 32, 34, 35, 39,
	40, 41, 42, 43, 46, 63, 65, 66, 87}

// QueryItemTypes returns a list of matching items given the query.
func (e *EveDB) QueryItemTypes(query string, catIDs ...int) ([]*ItemType, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer e.pool.Release(c)
	if len(catIDs) == 0 {
		// Default to Modules, Ships, Drones, and Charges
		catIDs = InterestingItemCategories
	}
	var cats []string
	for _, id := range catIDs {
		cats = append(cats, strconv.Itoa(id))
	}
	rs, err := c.Query(
		baseQueryItemType+
			`JOIN evesde."invGroups" grp
			    ON type."groupID" = grp."groupID" AND grp."categoryID" = ANY($1::INTEGER[])
			WHERE  type."published" = TRUE
			  AND type."typeName" ILIKE '%' || $2 || '%'
			LIMIT 20`, "{"+strings.Join(cats, ",")+"}", query)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*ItemType
	for rs.Next() {
		r := &ItemType{}
		err := rs.Scan(&r.ID, &r.Name, &r.Description)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

const baseQueryItemTypeDetail = `SELECT
  type."typeID"
, type."typeName"
, COALESCE(type."description", '')
, type."groupID"
, grp."groupName"
, grp."categoryID"
, cat."categoryName"
, COALESCE(type."mass", 0)
, COALESCE(type."volume", 0)
, COALESCE(type."capacity", 0)
, COALESCE(type."portionSize", 0)
, COALESCE(type."basePrice", 0)
, COALESCE(meta."parentTypeID", 0)
, COALESCE(bp."typeID", 0)
, COALESCE((SELECT STRING_AGG(metaDeriv."typeID"::VARCHAR, '|') FROM evesde."invMetaTypes" metaDeriv  WHERE metaDeriv."metaGroupID" = 2 AND metaDeriv."parentTypeID" = type."typeID" GROUP BY metaDeriv."parentTypeID"), '')
FROM evesde."invTypes" type
  JOIN evesde."invGroups" grp ON type."groupID" = grp."groupID"
  JOIN evesde."invCategories" cat ON grp."categoryID" = cat."categoryID"
  LEFT JOIN evesde."invMetaTypes" meta ON meta."metaGroupID" = 2 AND meta."typeID" = type."typeID"

  -- TODO: This is literally the worse join ever:
  LEFT JOIN evesde."invTypes" bp ON bp."typeName" = CONCAT(type."typeName", ' Blueprint')
` // TODO: Literally the worst literally ever

// GetItemTypeDetail fetches a specific ItemType with extra details from the database.
func (e *EveDB) GetItemTypeDetail(typeID int) (*ItemTypeDetail, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer e.pool.Release(c)
	var derivs string
	r := c.QueryRow(
		baseQueryItemTypeDetail+`WHERE type."typeID" = $1 AND type."published" = TRUE`, typeID)
	it := &ItemTypeDetail{ItemType: &ItemType{}}
	err = r.Scan(&it.ID, &it.Name, &it.Description, &it.GroupID, &it.GroupName, &it.CategoryID, &it.CategoryName, &it.Mass, &it.Volume, &it.Capacity, &it.PortionSize, &it.BasePrice, &it.ParentTypeID, &it.BlueprintID, &derivs)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(derivs, "|")
	var derivIDs []int
	for _, part := range parts {
		if v, err := strconv.Atoi(part); err == nil {
			derivIDs = append(derivIDs, v)
		}
	}
	it.DerivativeTypeIDs = derivIDs
	return it, nil
}

// QueryItemTypeDetails returns a list of matching items given the query.
func (e *EveDB) QueryItemTypeDetails(query string, catIDs ...int) ([]*ItemTypeDetail, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer e.pool.Release(c)
	if len(catIDs) == 0 {
		// Default to Modules, Ships, Drones, and Charges
		catIDs = InterestingItemCategories
	}
	var cats []string
	for _, id := range catIDs {
		cats = append(cats, strconv.Itoa(id))
	}
	rs, err := c.Query(
		baseQueryItemTypeDetail+
			`WHERE  type."published" = TRUE
			  AND type."typeName" ILIKE '%' || $2 || '%'
			  AND grp."categoryID" = ANY($1::INTEGER[])
			LIMIT 20`, "{"+strings.Join(cats, ",")+"}", query)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*ItemTypeDetail
	for rs.Next() {
		var derivs string
		it := &ItemTypeDetail{ItemType: &ItemType{}}
		err := rs.Scan(&it.ID, &it.Name, &it.Description, &it.GroupID, &it.GroupName, &it.CategoryID, &it.CategoryName, &it.Mass, &it.Volume, &it.Capacity, &it.PortionSize, &it.BasePrice, &it.ParentTypeID, &it.BlueprintID, &derivs)
		if err != nil {
			return nil, err
		}
		parts := strings.Split(derivs, "|")
		var derivIDs []int
		for _, part := range parts {
			if v, err := strconv.Atoi(part); err == nil {
				derivIDs = append(derivIDs, v)
			}
		}
		it.DerivativeTypeIDs = derivIDs
		res = append(res, it)
	}
	return res, nil
}

// A MaterialSheet describes what is necessary to build an item.
type MaterialSheet struct {
	*ItemType
	Materials   []*Material
	ProducesQty int
}

// A Material is a type and quantity of an item used for manufacturing.
type Material struct {
	*ItemType
	Quantity int
}

// GetBlueprint fetches a MaterialSheet from the database.
func (e *EveDB) GetBlueprint(typeID int) (*MaterialSheet, error) {
	it, err := e.GetItemTypeDetail(typeID)
	if err != nil {
		return nil, err
	}
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer e.pool.Release(c)
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
	var res []*Material
	for rs.Next() {
		r := &Material{ItemType: &ItemType{}}
		err := rs.Scan(&r.Name, &r.ID, &r.Quantity)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return &MaterialSheet{ItemType: it.ItemType, ProducesQty: it.PortionSize, Materials: res}, nil
}

// GetBlueprints is a utility function to retrieve multiple Blueprints.
func (e *EveDB) GetBlueprints(typeIDs ...int) ([]*MaterialSheet, error) {
	var res []*MaterialSheet
	for _, id := range typeIDs {
		bp, err := e.GetBlueprint(id)
		if err != nil {
			return nil, err
		}
		res = append(res, bp)
	}
	return res, nil
}
