package eveapi

import (
	"time"

	"golang.org/x/net/context"
)

// A Structure is a player-owned citadel.
type Structure struct {
	StructureID int64
	Name        string
	SystemID    int64
	TypeID      int64
}

// A CorporationStructure contains additional, sensitive information about a citadel.
type CorporationStructure struct {
	Structure
	ProfileID            int64
	Services             []string
	FuelExpires          time.Time
	StateStart           time.Time
	StateEnd             time.Time
	UnanchorsAt          time.Time
	State                string
	ReinforceWeekday     int32
	ReinforceHour        int32
	NextReinforceWeekday int32
	NextReinforceHour    int32
	NextReinforceTime    time.Time
}

func (api *EveAPI) GetCorporationStructures(ctx context.Context, corpID int) ([]*CorporationStructure, error) {
	res, _, err := api.client.ESI.CorporationApi.GetCorporationsCorporationIdStructures(ctx, int32(corpID), nil)
	if err != nil {
		return nil, err
	}
	var structures []*CorporationStructure
	for _, bp := range res {
		// ESI doesn't return the structure name in this API call for some reason.
		// Query the Universe ESI API for the structure's name.
		s, err := api.GetStructure(ctx, bp.StructureId)
		if err != nil {
			return nil, err
		}
		var srvs []string
		for _, r := range bp.Services {
			srvs = append(srvs, r.Name)
		}
		structures = append(structures, &CorporationStructure{
			Structure:            *s,
			ProfileID:            int64(bp.ProfileId),
			UnanchorsAt:          bp.UnanchorsAt,
			StateStart:           bp.StateTimerStart,
			StateEnd:             bp.StateTimerEnd,
			Services:             srvs,
			FuelExpires:          bp.FuelExpires,
			State:                bp.State,
			ReinforceWeekday:     bp.ReinforceWeekday,
			ReinforceHour:        bp.ReinforceHour,
			NextReinforceWeekday: bp.NextReinforceWeekday,
			NextReinforceHour:    bp.NextReinforceHour,
			NextReinforceTime:    bp.NextReinforceApply,
		})
	}
	return structures, nil
}

func (api *EveAPI) GetStructure(ctx context.Context, structureID int64) (*Structure, error) {
	res, _, err := api.client.ESI.UniverseApi.GetUniverseStructuresStructureId(ctx, int64(structureID), nil)
	if err != nil {
		return nil, err
	}
	return &Structure{
		StructureID: structureID,
		Name:        res.Name,
		SystemID:    int64(res.SolarSystemId),
		TypeID:      int64(res.TypeId),
	}, nil
}
