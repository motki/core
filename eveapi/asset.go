package eveapi

import (
	"golang.org/x/net/context"
)

type Asset struct {
	ItemID      int
	LocationID  int
	TypeID      int
	Quantity    int
	FlagID      int
	Singleton   bool
	RawQuantity int
}

func (api *EveAPI) GetCorporationAssets(ctx context.Context, corpID int) ([]*Asset, error) {
	tok, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := api.client.EVEAPI.CorporationAssetsXML(tok, int64(corpID))
	if err != nil {
		return nil, err
	}
	var assets []*Asset
	for _, a := range res.Entries {
		assets = append(assets, &Asset{
			ItemID:      int(a.ItemID),
			LocationID:  int(a.LocationID),
			TypeID:      int(a.TypeID),
			Quantity:    int(a.Quantity),
			FlagID:      int(a.FlagID),
			Singleton:   a.Singleton,
			RawQuantity: int(a.RawQuantity),
		})
	}
	return assets, nil
}
