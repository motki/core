package model

import (
	"context"

	"github.com/motki/motkid/eveapi"
)

func (m *Manager) GetCorporationBlueprints(ctx context.Context, corpID int) (jobs []*eveapi.Blueprint, err error) {
	jobs, err = m.getCorporationBlueprintsFromDB(corpID)
	if err != nil {
		return nil, err
	}
	if jobs != nil {
		return jobs, nil
	}
	return m.getCorporationBlueprintsFromAPI(ctx, corpID)
}

func (m *Manager) getCorporationBlueprintsFromDB(corpID int) ([]*eveapi.Blueprint, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rs, err := c.Query(
		`SELECT
			  c.item_id
			, c.location_id
			, c.type_id
			, c.type_name
			, c.quantity
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
	res := []*eveapi.Blueprint{}
	for rs.Next() {
		r := &eveapi.Blueprint{}
		err := rs.Scan(
			&r.ItemID,
			&r.LocationID,
			&r.TypeID,
			&r.TypeName,
			&r.Quantity,
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

func (m *Manager) getCorporationBlueprintsFromAPI(ctx context.Context, corpID int) ([]*eveapi.Blueprint, error) {
	bps, err := m.eveapi.GetCorporationBlueprints(ctx, corpID)
	if err != nil {
		return nil, err
	}
	return m.apiCorporationBlueprintsToDB(corpID, bps)
}

func (m *Manager) apiCorporationBlueprintsToDB(corpID int, bps []*eveapi.Blueprint) ([]*eveapi.Blueprint, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	for _, bp := range bps {
		_, err = db.Exec(
			`INSERT INTO app.blueprints
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, DEFAULT)`,
			corpID,
			0,
			bp.ItemID,
			bp.LocationID,
			bp.TypeID,
			bp.TypeName,
			bp.Quantity,
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
