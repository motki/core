package model

import (
	"database/sql"
	"encoding/json"

	"golang.org/x/net/context"

	"github.com/motki/core/eveapi"
)

type StructureManager struct {
	bootstrap

	corp *CorpManager
}

func newStructureManager(m bootstrap, corp *CorpManager) *StructureManager {
	return &StructureManager{m, corp}
}

func (m *StructureManager) GetCorporationStructures(ctx context.Context, corpID int) ([]*eveapi.CorporationStructure, error) {
	var err error
	if ctx, err = m.corp.authContext(ctx, corpID); err != nil {
		return nil, err
	}
	if jobs, err := m.getCorporationStructuresFromDB(corpID); err == nil && jobs != nil {
		return jobs, nil
	} else if err != nil {
		return nil, err
	}
	return m.getCorporationStructuresFromAPI(ctx, corpID)
}

func (m *StructureManager) GetStructure(ctx context.Context, structureID int) (*eveapi.Structure, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	r := c.QueryRow(
		`SELECT
			  c.structure_id
			, c.name
			, c.system_id
			, c.type_id
			FROM app.structures c
			WHERE c.structure_id = $1 
 			  AND c.fetched_at > (NOW() - INTERVAL '12 hours')`, structureID)
	s := &eveapi.Structure{}
	err = r.Scan(&s.StructureID, &s.Name, &s.SystemID, &s.TypeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return m.getStructureFromAPI(ctx, structureID)
		}
		return nil, err
	}
	return s, nil
}

func (m *StructureManager) getCorporationStructuresFromDB(corpID int) ([]*eveapi.CorporationStructure, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  c.structure_id
			, c.name
			, c.system_id
			, c.type_id
			, c.profile_id
			, c.fuel_expires
			, c.services
			, c.state_timer_start
			, c.state_timer_end
			, c.curr_state
			, c.unanchors_at
			, c.reinforce_weekday
			, c.reinforce_hour
			FROM app.structures c
			WHERE c.corporation_id = $1
			  AND c.fetched_at > (NOW() - INTERVAL '12 hours')`, corpID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*eveapi.CorporationStructure
	for rs.Next() {
		r := &eveapi.CorporationStructure{}
		var s []byte
		err := rs.Scan(
			&r.StructureID,
			&r.Name,
			&r.SystemID,
			&r.TypeID,
			&r.ProfileID,
			&r.FuelExpires,
			&s,
			&r.StateStart,
			&r.StateEnd,
			&r.State,
			&r.UnanchorsAt,
			&r.VulnWeekday,
			&r.VulnHour,
		)
		if err != nil {
			return nil, err
		}
		var srvs []string
		if err = json.Unmarshal(s, &srvs); err != nil {
			return nil, err
		}
		r.Services = srvs
		res = append(res, r)
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (m *StructureManager) getCorporationStructuresFromAPI(ctx context.Context, corpID int) ([]*eveapi.CorporationStructure, error) {
	strucs, err := m.eveapi.GetCorporationStructures(ctx, corpID)
	if err != nil {
		return nil, err
	}
	return m.apiCorporationStructuresToDB(corpID, strucs)
}

func (m *StructureManager) apiCorporationStructuresToDB(corpID int, strucs []*eveapi.CorporationStructure) ([]*eveapi.CorporationStructure, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	for _, struc := range strucs {
		b, err := json.Marshal(struc.Services)
		if err != nil {
			return nil, err
		}
		_, err = db.Exec(
			`INSERT INTO app.structures
					(structure_id, corporation_id, system_id, type_id, profile_id,
						fuel_expires, services, state_timer_start, state_timer_end,
						curr_state, unanchors_at, reinforce_weekday, reinforce_hour,
						name, fetched_at)
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, DEFAULT)
				ON CONFLICT ON CONSTRAINT "structures_pkey"
				  DO UPDATE SET corporation_id = EXCLUDED.corporation_id,
					name = EXCLUDED.name,
					profile_id = EXCLUDED.profile_id,
					unanchors_at = EXCLUDED.unanchors_at,
					state_timer_start = EXCLUDED.state_timer_start,
					state_timer_end = EXCLUDED.state_timer_end,
					services = EXCLUDED.services,
					fuel_expires = EXCLUDED.fuel_expires,
					curr_state = EXCLUDED.curr_state,
					reinforce_weekday = EXCLUDED.reinforce_weekday,
					reinforce_hour = EXCLUDED.reinforce_hour,
					fetched_at = DEFAULT`,
			struc.StructureID,
			corpID,
			struc.SystemID,
			struc.TypeID,
			struc.ProfileID,
			struc.FuelExpires,
			b,
			struc.StateStart,
			struc.StateEnd,
			struc.State,
			struc.UnanchorsAt,
			struc.VulnWeekday,
			struc.VulnHour,
			struc.Name,
		)
		if err != nil {
			return nil, err
		}
	}
	return strucs, nil
}

func (m *StructureManager) getStructureFromAPI(ctx context.Context, structureID int) (*eveapi.Structure, error) {
	s, err := m.eveapi.GetStructure(ctx, int64(structureID))
	if err != nil {
		return nil, err
	}
	return m.apiStructureToDB(s)
}

func (m *StructureManager) apiStructureToDB(struc *eveapi.Structure) (*eveapi.Structure, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	_, err = db.Exec(
		`INSERT INTO app.structures
					(structure_id, system_id, type_id, name, fetched_at)
					VALUES($1, $2, $3, $4, DEFAULT)
				ON CONFLICT ON CONSTRAINT "structures_pkey"
				  DO UPDATE SET name = EXCLUDED.name,
					fetched_at = DEFAULT`,
		struc.StructureID,
		struc.SystemID,
		struc.TypeID,
		struc.Name,
	)
	if err != nil {
		return nil, err
	}
	return struc, nil
}
