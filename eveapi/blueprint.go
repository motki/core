package eveapi

import (
	"strconv"

	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"golang.org/x/net/context"
)

type Blueprint struct {
	ItemID             int64
	LocationID         int64
	TypeID             int64
	LocationFlag       string
	TimeEfficiency     int64
	MaterialEfficiency int64

	// -2 = BPC (and always qty 1), else BPO
	Quantity int64

	// -1 = infinite runs (a BPO)
	Runs int64
}

func (api *EveAPI) GetCorporationBlueprints(ctx context.Context, corpID int) ([]*Blueprint, error) {
	_, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var bps []*Blueprint
	for p, max := 1, 1; p <= max; p++ {
		res, resp, err := api.client.ESI.CorporationApi.GetCorporationsCorporationIdBlueprints(
			ctx,
			int32(corpID),
			&esi.GetCorporationsCorporationIdBlueprintsOpts{Page: optional.NewInt32(int32(p))})
		if err != nil {
			return nil, err
		}
		max, err = strconv.Atoi(resp.Header.Get("X-Pages"))
		if err != nil {
			api.logger.Debugf("error reading X-Pages header: ", err.Error())
		}
		for _, bp := range res {
			bps = append(bps, &Blueprint{
				ItemID:             bp.ItemId,
				LocationID:         bp.LocationId,
				LocationFlag:       bp.LocationFlag,
				TypeID:             int64(bp.TypeId),
				Quantity:           int64(bp.Quantity),
				TimeEfficiency:     int64(bp.TimeEfficiency),
				MaterialEfficiency: int64(bp.MaterialEfficiency),
				Runs:               int64(bp.Runs),
			})
		}
	}
	return bps, nil
}
