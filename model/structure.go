package model

import (
	"golang.org/x/net/context"

	"github.com/motki/motki/eveapi"
)

func (m *Manager) GetCorporationStructures(ctx context.Context, corpID int) ([]*eveapi.Structure, error) {
	if jobs, err := m.getCorporationStructuresFromDB(corpID); err == nil && jobs != nil {
		return jobs, nil
	} else if err != nil {
		return nil, err
	}
	return m.getCorporationStructuresFromAPI(ctx, corpID)
}

func (m *Manager) getCorporationStructuresFromDB(corpID int) ([]*eveapi.Structure, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  c.structure_id
			, c.system_id
			, c.type_id
			, c.profile_id
			, c.curr_vuln
			, c.next_vuln
			FROM app.structures c
			WHERE c.corporation_id = $1
			  AND c.fetched_at > (NOW() - INTERVAL '12 hours')`, corpID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*eveapi.Structure
	for rs.Next() {
		r := &eveapi.Structure{}
		err := rs.Scan(
			&r.StructureID,
			&r.SystemID,
			&r.TypeID,
			&r.ProfileID,
			&r.CurrentVuln,
			&r.NextVuln,
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

func (m *Manager) getCorporationStructuresFromAPI(ctx context.Context, corpID int) ([]*eveapi.Structure, error) {
	strucs, err := m.eveapi.GetCorporationStructures(ctx, corpID)
	if err != nil {
		return nil, err
	}
	return m.apiCorporationStructuresToDB(corpID, strucs)
}

func (m *Manager) apiCorporationStructuresToDB(corpID int, strucs []*eveapi.Structure) ([]*eveapi.Structure, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	for _, struc := range strucs {
		_, err = db.Exec(
			`INSERT INTO app.structures
					VALUES($1, $2, $3, $4, $5, $6, $7, DEFAULT)`,
			corpID,
			struc.StructureID,
			struc.SystemID,
			struc.TypeID,
			struc.ProfileID,
			struc.CurrentVuln,
			struc.NextVuln,
		)
		if err != nil {
			return nil, err
		}
	}
	return strucs, nil
}
