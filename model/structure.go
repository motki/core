package model

import (
	"database/sql"
	"encoding/json"

	"golang.org/x/net/context"

	"time"

	"fmt"

	"database/sql/driver"

	"github.com/motki/core/eveapi"
)

// A Structure is a player-owned citadel.
type Structure struct {
	StructureID int64  `json:"structure_id"`
	Name        string `json:"name"`
	SystemID    int64  `json:"system_id"`
	TypeID      int64  `json:"type_id"`
}

type ReinforceHour int

func (h ReinforceHour) String() string {
	return fmt.Sprintf("%d:00 UTC", h)
}

type ReinforceWindow struct {
	Weekday     time.Weekday  `json:"weekday"`
	Hour        ReinforceHour `json:"hour"`
	EffectiveAt time.Time     `json:"effective_at"`
}

func (h ReinforceWindow) String() string {
	return fmt.Sprintf("%s at %s", h.Weekday, h.Hour)
}

func (h ReinforceWindow) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *ReinforceWindow) Scan(src interface{}) error {
	s, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("invalid value for reinforcement window: %v", src)
	}
	return json.Unmarshal(s, &h)
}

// A CorporationStructure contains additional, sensitive information about a citadel.
type CorporationStructure struct {
	Structure
	ProfileID   int64     `json:"profile_id"`
	Services    []string  `json:"services"`
	FuelExpires time.Time `json:"fuel_expires"`
	StateStart  time.Time `json:"state_start"`
	StateEnd    time.Time `json:"state_end"`
	UnanchorsAt time.Time `json:"unanchors_at"`
	State       string    `json:"state"`

	CurrReinforceWindow ReinforceWindow `json:"curr_reinforce_window"`
	NextReinforceWindow ReinforceWindow `json:"next_reinforce_window"`
}

type StructureManager struct {
	bootstrap

	corp *CorpManager
}

func newStructureManager(m bootstrap, corp *CorpManager) *StructureManager {
	return &StructureManager{m, corp}
}

func (m *StructureManager) GetCorporationStructures(ctx context.Context, corpID int) ([]*CorporationStructure, error) {
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

func (m *StructureManager) GetStructure(ctx context.Context, structureID int) (*Structure, error) {
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
	s := &Structure{}
	err = r.Scan(&s.StructureID, &s.Name, &s.SystemID, &s.TypeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return m.getStructureFromAPI(ctx, structureID)
		}
		return nil, err
	}
	return s, nil
}

func (m *StructureManager) getCorporationStructuresFromDB(corpID int) ([]*CorporationStructure, error) {
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
			, c.curr_reinforce_window
			, c.next_reinforce_window
			FROM app.structures c
			WHERE c.corporation_id = $1
			  AND c.fetched_at > (NOW() - INTERVAL '12 hours')`, corpID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*CorporationStructure
	for rs.Next() {
		r := &CorporationStructure{
			CurrReinforceWindow: ReinforceWindow{EffectiveAt: time.Now()},
			NextReinforceWindow: ReinforceWindow{},
		}
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
			&r.CurrReinforceWindow,
			&r.NextReinforceWindow,
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

func (m *StructureManager) getCorporationStructuresFromAPI(ctx context.Context, corpID int) ([]*CorporationStructure, error) {
	strucs, err := m.eveapi.GetCorporationStructures(ctx, corpID)
	if err != nil {
		return nil, err
	}
	return m.apiCorporationStructuresToDB(corpID, strucs)
}

func (m *StructureManager) apiCorporationStructuresToDB(corpID int, strucs []*eveapi.CorporationStructure) ([]*CorporationStructure, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	res := make([]*CorporationStructure, len(strucs))
	for i, rs := range strucs {
		b, err := json.Marshal(rs.Services)
		if err != nil {
			return nil, err
		}
		s := &CorporationStructure{
			Structure:   (Structure)(rs.Structure),
			ProfileID:   rs.ProfileID,
			Services:    rs.Services,
			FuelExpires: rs.FuelExpires,
			StateStart:  rs.StateStart,
			StateEnd:    rs.StateEnd,
			UnanchorsAt: rs.UnanchorsAt,
			CurrReinforceWindow: ReinforceWindow{
				Weekday:     time.Weekday(rs.ReinforceWeekday),
				Hour:        ReinforceHour(rs.ReinforceHour),
				EffectiveAt: time.Now(),
			},
			NextReinforceWindow: ReinforceWindow{
				Weekday:     time.Weekday(rs.NextReinforceWeekday),
				Hour:        ReinforceHour(rs.NextReinforceHour),
				EffectiveAt: rs.NextReinforceTime,
			},
			State: rs.State,
		}
		_, err = db.Exec(
			`INSERT INTO app.structures
					(structure_id, corporation_id, system_id, type_id, profile_id,
						fuel_expires, services, state_timer_start, state_timer_end,
						curr_state, unanchors_at, name, curr_reinforce_window, 
						next_reinforce_window, fetched_at)
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
					curr_reinforce_window = EXCLUDED.curr_reinforce_window,
					next_reinforce_window = EXCLUDED.next_reinforce_window,
					fetched_at = DEFAULT`,
			s.StructureID,
			corpID,
			s.SystemID,
			s.TypeID,
			s.ProfileID,
			s.FuelExpires,
			b,
			s.StateStart,
			s.StateEnd,
			s.State,
			s.UnanchorsAt,
			s.Name,
			s.CurrReinforceWindow,
			s.NextReinforceWindow,
		)
		if err != nil {
			return nil, err
		}
		res[i] = s
	}
	return res, nil
}

func (m *StructureManager) getStructureFromAPI(ctx context.Context, structureID int) (*Structure, error) {
	s, err := m.eveapi.GetStructure(ctx, int64(structureID))
	if err != nil {
		return nil, err
	}
	return m.apiStructureToDB(s)
}

func (m *StructureManager) apiStructureToDB(struc *eveapi.Structure) (*Structure, error) {
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
	return &Structure{
		StructureID: struc.StructureID,
		SystemID:    struc.SystemID,
		TypeID:      struc.TypeID,
		Name:        struc.Name,
	}, nil
}
