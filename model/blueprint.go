package model

import (
	"golang.org/x/net/context"

	"github.com/motki/core/eveapi"
)

type BlueprintKind string

var (
	BlueprintOriginal BlueprintKind = "bpo"
	BlueprintCopy     BlueprintKind = "bpc"
)

type Blueprint struct {
	ItemID             int
	LocationID         int
	TypeID             int
	TypeName           string
	FlagID             int
	TimeEfficiency     int
	MaterialEfficiency int
	Kind               BlueprintKind
	Quantity           int

	// -1 = infinite runs (a BPO)
	Runs int
}

func blueprintFromEveAPI(bp *eveapi.Blueprint) *Blueprint {
	kind := BlueprintOriginal
	qty := bp.Quantity
	if qty == -2 {
		kind = BlueprintCopy
		qty = 1
	}
	return &Blueprint{
		ItemID:             int(bp.ItemID),
		LocationID:         int(bp.LocationID),
		TypeID:             int(bp.TypeID),
		TypeName:           bp.TypeName,
		FlagID:             int(bp.FlagID),
		TimeEfficiency:     int(bp.TimeEfficiency),
		MaterialEfficiency: int(bp.MaterialEfficiency),
		Kind:               kind,
		Quantity:           int(qty),
		Runs:               int(bp.Runs),
	}
}

type BlueprintManager struct {
	bootstrap

	corp *CorpManager
}

func newBlueprintManager(m bootstrap, corp *CorpManager) *BlueprintManager {
	return &BlueprintManager{m, corp}
}

func (m *BlueprintManager) GetCorporationBlueprints(ctx context.Context, corpID int) (jobs []*Blueprint, err error) {
	if ctx, err = m.corp.authContext(ctx, corpID); err != nil {
		return nil, err
	}
	jobs, err = m.getCorporationBlueprintsFromDB(corpID)
	if err != nil {
		return nil, err
	}
	if jobs != nil {
		return jobs, nil
	}
	return m.getCorporationBlueprintsFromAPI(ctx, corpID)
}

func (m *BlueprintManager) getCorporationBlueprintsFromDB(corpID int) ([]*Blueprint, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  c.item_id
			, c.location_id
			, c.type_id
			, c.type_name
			, c.quantity
			, c.kind
			, c.flag_id
			, c.time_efficiency
			, c.material_efficiency
			, c.runs
			FROM app.blueprints c
			WHERE c.corporation_id = $1
			  AND c.fetched_at > (NOW() - INTERVAL '12 hours')`, corpID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*Blueprint
	for rs.Next() {
		r := &Blueprint{}
		err := rs.Scan(
			&r.ItemID,
			&r.LocationID,
			&r.TypeID,
			&r.TypeName,
			&r.Quantity,
			&r.Kind,
			&r.FlagID,
			&r.TimeEfficiency,
			&r.MaterialEfficiency,
			&r.Runs,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (m *BlueprintManager) getCorporationBlueprintsFromAPI(ctx context.Context, corpID int) ([]*Blueprint, error) {
	bps, err := m.eveapi.GetCorporationBlueprints(ctx, corpID)
	if err != nil {
		return nil, err
	}
	var res []*Blueprint
	for _, bp := range bps {
		res = append(res, blueprintFromEveAPI(bp))
	}
	return m.apiCorporationBlueprintsToDB(corpID, res)
}

func (m *BlueprintManager) apiCorporationBlueprintsToDB(corpID int, bps []*Blueprint) ([]*Blueprint, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	for _, bp := range bps {
		_, err = db.Exec(
			`INSERT INTO app.blueprints
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, DEFAULT)`,
			corpID,
			0,
			bp.ItemID,
			bp.LocationID,
			bp.TypeID,
			bp.TypeName,
			bp.Quantity,
			bp.Kind,
			bp.FlagID,
			bp.TimeEfficiency,
			bp.MaterialEfficiency,
			bp.Runs,
		)
		if err != nil {
			return nil, err
		}
	}
	return bps, nil
}
