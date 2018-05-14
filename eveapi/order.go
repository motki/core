package eveapi

import (
	"strconv"
	"time"

	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
)

type MarketOrder struct {
	OrderID      int
	CharID       int // TODO: Doesn't exist in the ESI response
	LocationID   int
	TypeID       int
	VolEntered   int
	VolRemaining int
	MinVolume    int
	OrderState   string
	Range        string
	AccountKey   int
	Duration     int
	Escrow       decimal.Decimal
	Price        decimal.Decimal
	Bid          bool
	Issued       time.Time
}

func (api *EveAPI) GetCorporationOrders(ctx context.Context, corpID int) (orders []*MarketOrder, err error) {
	_, err = TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	for max, p := 1, 1; p <= max; p++ {
		res, resp, err := api.client.ESI.MarketApi.GetCorporationsCorporationIdOrders(ctx, int32(corpID), &esi.GetCorporationsCorporationIdOrdersOpts{Page: optional.NewInt32(int32(p))})
		if err != nil {
			return nil, err
		}
		max, err = strconv.Atoi(resp.Header.Get("X-Pages"))
		if err != nil {
			api.logger.Debugf("error reading X-Pages header: ", err.Error())
		}
		for _, j := range res {
			order := &MarketOrder{
				OrderID: int(j.OrderId),
				//CharID:       int(j.CharId),
				LocationID:   int(j.LocationId),
				TypeID:       int(j.TypeId),
				VolEntered:   int(j.VolumeTotal),
				VolRemaining: int(j.VolumeRemain),
				MinVolume:    int(j.MinVolume),
				OrderState:   "open",
				Range:        j.Range_,
				AccountKey:   int(j.WalletDivision),
				Duration:     int(j.Duration),
				Escrow:       decimal.NewFromFloat(j.Escrow),
				Price:        decimal.NewFromFloat(j.Price),
				Bid:          j.IsBuyOrder,
				Issued:       j.Issued,
			}
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (api *EveAPI) GetCorporationOrdersHistory(ctx context.Context, corpID int) (orders []*MarketOrder, err error) {
	_, err = TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	for max, p := 1, 1; p <= max; p++ {
		res, resp, err := api.client.ESI.MarketApi.GetCorporationsCorporationIdOrdersHistory(
			ctx,
			int32(corpID),
			&esi.GetCorporationsCorporationIdOrdersHistoryOpts{Page: optional.NewInt32(int32(p))})
		if err != nil {
			return nil, err
		}
		max, err = strconv.Atoi(resp.Header.Get("X-Pages"))
		if err != nil {
			api.logger.Debugf("error reading X-Pages header: ", err.Error())
		}
		for _, j := range res {
			order := &MarketOrder{
				OrderID: int(j.OrderId),
				//CharID:       int(j.CharId),
				LocationID:   int(j.LocationId),
				TypeID:       int(j.TypeId),
				VolEntered:   int(j.VolumeTotal),
				VolRemaining: int(j.VolumeRemain),
				MinVolume:    int(j.MinVolume),
				OrderState:   j.State,
				Range:        j.Range_,
				AccountKey:   int(j.WalletDivision),
				Duration:     int(j.Duration),
				Escrow:       decimal.NewFromFloat(j.Escrow),
				Price:        decimal.NewFromFloat(j.Price),
				Bid:          j.IsBuyOrder,
				Issued:       j.Issued,
			}
			orders = append(orders, order)
		}
	}

	return orders, nil
}
