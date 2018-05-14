package eveapi

import (
	"strconv"

	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"golang.org/x/net/context"
)

type Asset struct {
	ItemID       int
	LocationID   int
	LocationType string
	LocationFlag string
	TypeID       int
	Quantity     int
	Singleton    bool
}

func (api *EveAPI) GetCorporationAssets(ctx context.Context, corpID int) ([]*Asset, error) {
	_, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var assets []*Asset
	for max, p := 1, 1; p <= max; p++ {
		res, resp, err := api.client.ESI.AssetsApi.GetCorporationsCorporationIdAssets(
			ctx,
			int32(corpID),
			&esi.GetCorporationsCorporationIdAssetsOpts{Page: optional.NewInt32(int32(p))})
		if err != nil {
			return nil, err
		}
		max, err = strconv.Atoi(resp.Header.Get("X-Pages"))
		if err != nil {
			api.logger.Debugf("error reading X-Pages header: ", err.Error())
		}
		for _, a := range res {
			assets = append(assets, &Asset{
				ItemID:       int(a.ItemId),
				LocationID:   int(a.LocationId),
				LocationType: a.LocationType,
				LocationFlag: a.LocationFlag,
				TypeID:       int(a.TypeId),
				Quantity:     int(a.Quantity),
				Singleton:    a.IsSingleton,
			})
		}
	}

	return assets, nil
}
