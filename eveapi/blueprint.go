package eveapi

import (
	"context"
)

type Blueprint struct {
	ItemID             int64
	LocationID         int64
	TypeID             int64
	TypeName           string
	Quantity           int64
	FlagID             int64
	TimeEfficiency     int64
	MaterialEfficiency int64
	Runs               int64
}

func (api *EveAPI) GetCorporationBlueprints(ctx context.Context, corpID int) ([]*Blueprint, error) {
	tok, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := api.client.EVEAPI.CorporationBlueprintsXML(tok, int64(corpID))
	if err != nil {
		return nil, err
	}
	bps := []*Blueprint{}
	for _, bp := range res.Entries {
		bps = append(bps, &Blueprint{
			ItemID:             bp.ItemID,
			LocationID:         bp.LocationID,
			TypeID:             bp.TypeID,
			TypeName:           bp.TypeName,
			Quantity:           bp.Quantity,
			FlagID:             bp.FlagID,
			TimeEfficiency:     bp.TimeEfficiency,
			MaterialEfficiency: bp.MaterialEfficiency,
			Runs:               bp.Runs,
		})
	}
	return bps, nil
}
