package eveapi

import (
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
)

type MarketOrder struct {
	OrderID      int
	CharID       int
	StationID    int
	TypeID       int
	VolEntered   int
	VolRemaining int
	MinVolume    int
	OrderState   int
	Range        int
	AccountKey   int
	Duration     int
	Escrow       decimal.Decimal
	Price        decimal.Decimal
	Bid          bool
	Issued       time.Time
}

func (api *EveAPI) GetCorporationOrders(ctx context.Context, corpID int) (orders []*MarketOrder, err error) {
	tok, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := api.client.EVEAPI.CorporationMarketOrdersXML(tok, int64(corpID))
	if err != nil {
		return nil, err
	}
	for _, j := range res.Entries {
		order := &MarketOrder{
			OrderID:      int(j.OrderID),
			CharID:       int(j.CharID),
			StationID:    int(j.StationID),
			TypeID:       int(j.TypeID),
			VolEntered:   int(j.VolEntered),
			VolRemaining: int(j.VolRemaining),
			MinVolume:    int(j.MinVolume),
			OrderState:   int(j.OrderState),
			Range:        int(j.Range),
			AccountKey:   int(j.AccountKey),
			Duration:     int(j.Duration),
			Escrow:       decimal.NewFromFloat(j.Escrow),
			Price:        decimal.NewFromFloat(j.Price),
			Bid:          j.Bid,
			Issued:       j.Issued.Time,
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (api *EveAPI) GetCorporationOrder(ctx context.Context, corpID, orderID int) (*MarketOrder, error) {
	tok, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := api.client.EVEAPI.CorporationMarketOrderXML(tok, int64(corpID), int64(orderID))
	if err != nil {
		return nil, err
	}
	for _, j := range res.Entries {
		return &MarketOrder{
			OrderID:      int(j.OrderID),
			CharID:       int(j.CharID),
			StationID:    int(j.StationID),
			TypeID:       int(j.TypeID),
			VolEntered:   int(j.VolEntered),
			VolRemaining: int(j.VolRemaining),
			MinVolume:    int(j.MinVolume),
			OrderState:   int(j.OrderState),
			Range:        int(j.Range),
			AccountKey:   int(j.AccountKey),
			Duration:     int(j.Duration),
			Escrow:       decimal.NewFromFloat(j.Escrow),
			Price:        decimal.NewFromFloat(j.Price),
			Bid:          j.Bid,
			Issued:       j.Issued.Time,
		}, nil
	}
	return nil, errors.New("not found")
}
