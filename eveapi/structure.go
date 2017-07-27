package eveapi

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/antihax/goesi/esi"
)

type Structure struct {
	StructureID int64
	SystemID    int64
	TypeID      int64
	ProfileID   int64
	CurrentVuln VulnSchedule
	NextVuln    VulnSchedule
}

type VulnSchedule map[int][]int

func (s *VulnSchedule) Scan(src interface{}) error {
	if v, ok := src.(string); ok {
		return json.Unmarshal([]byte(v), &s)
	}
	return fmt.Errorf("invalid vulnerability schedule: %v", src)
}

func (s VulnSchedule) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s VulnSchedule) String() string {
	out := []string{}
	for day, hrs := range s {
		var d string
		switch day {
		case 0:
			d = "Sunday"
		case 1:
			d = "Monday"
		case 2:
			d = "Tuesday"
		case 3:
			d = "Wednesday"
		case 4:
			d = "Thursday"
		case 5:
			d = "Friday"
		default:
			d = "Saturday"
		}
		td := []string{}
		for _, hr := range hrs {
			td = append(td, fmt.Sprintf("%02d00", hr))
		}
		if len(td) != 0 {
			out = append(out, fmt.Sprintf("%s: %s", d, strings.Join(td, ", ")))
		}
	}
	return strings.Join(out, ", ")
}

func (api *EveAPI) GetCorporationStructures(ctx context.Context, corpID int) ([]*Structure, error) {
	res, _, err := api.client.ESI.CorporationApi.GetCorporationsCorporationIdStructures(ctx, int32(corpID), nil)
	if err != nil {
		return nil, err
	}
	structures := []*Structure{}
	for _, bp := range res {
		sched := map[int][]int{}
		for _, sch := range bp.CurrentVul {
			d, h := int(sch.Day), int(sch.Hour)
			if _, ok := sched[d]; !ok {
				sched[d] = []int{}
			}
			sched[d] = append(sched[d], h)
		}
		nsched := map[int][]int{}
		for _, sch := range bp.NextVul {
			d, h := int(sch.Day), int(sch.Hour)
			if _, ok := nsched[d]; !ok {
				nsched[d] = []int{}
			}
			nsched[d] = append(nsched[d], h)
		}
		structures = append(structures, &Structure{
			StructureID: bp.StructureId,
			SystemID:    int64(bp.SystemId),
			TypeID:      int64(bp.TypeId),
			ProfileID:   int64(bp.ProfileId),
			CurrentVuln: sched,
			NextVuln:    nsched,
		})
	}
	return structures, nil
}

func (api *EveAPI) UpdateCorporationStructureVulnSchedule(ctx context.Context, corpID int, structureID int, sched VulnSchedule) error {
	newSched := []esi.PutCorporationsCorporationIdStructuresStructureIdNewSchedule{}
	for day, hrs := range sched {
		for _, hr := range hrs {
			newSched = append(newSched, esi.PutCorporationsCorporationIdStructuresStructureIdNewSchedule{
				Day:  int32(day),
				Hour: int32(hr),
			})
		}
	}
	_, err := api.client.ESI.CorporationApi.PutCorporationsCorporationIdStructuresStructureId(ctx, int32(corpID), newSched, int64(structureID), nil)
	if err != nil {
		return err
	}
	return nil
}
