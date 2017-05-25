package evedb

// An ItemType is a type of item in EVE.
type ItemType struct {
	ID   int
	Name string
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
			  type."typeID" as ID
			, type."typeName" as typeName
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
