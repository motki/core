package model

import (
	"database/sql/driver"
	"time"

	"fmt"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
)

// One of: open, expired, cancelled
type OrderState string

const (
	OrderStateOpen      OrderState = "open"
	OrderStateExpired              = "expired"
	OrderStateCancelled            = "cancelled"
)

func (r OrderState) Value() (driver.Value, error) {
	return string(r), nil
}

func (r *OrderState) Scan(src interface{}) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("invalid %t for order state: %v", src, src)
	}
	switch OrderState(s) {
	case OrderStateOpen:
		*r = OrderStateOpen
	case OrderStateCancelled:
		*r = OrderStateCancelled
	case OrderStateExpired:
		*r = OrderStateExpired
	default:
		return fmt.Errorf("invalid value for order state: %v", s)
	}
	return nil
}

type MarketOrder struct {
	OrderID      int             `json:"order_id"`
	CharID       int             `json:"char_id"`
	LocationID   int             `json:"location_id"`
	TypeID       int             `json:"type_id"`
	VolEntered   int             `json:"vol_entered"`
	VolRemaining int             `json:"vol_remaining"`
	MinVolume    int             `json:"min_volume"`
	OrderState   OrderState      `json:"order_state"`
	Range        string          `json:"range"`
	AccountKey   int             `json:"account_key"`
	Duration     int             `json:"duration"`
	Escrow       decimal.Decimal `json:"escrow"`
	Price        decimal.Decimal `json:"price"`
	Bid          bool            `json:"bid"`
	Issued       time.Time       `json:"issued"`
}

func (m *MarketManager) GetCorporationOrder(ctx context.Context, corpID, orderID int) (*MarketOrder, error) {
	var err error
	if ctx, err = m.corp.authContext(ctx, corpID); err != nil {
		return nil, err
	}
	orders, err := m.GetCorporationOrders(ctx, corpID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	for _, o := range orders {
		if o.OrderID == orderID {
			return o, nil
		}
	}
	return nil, errors.New("order not found")
}

func (m *MarketManager) GetCorporationOrders(ctx context.Context, corpID int) (orders []*MarketOrder, err error) {
	if ctx, err = m.corp.authContext(ctx, corpID); err != nil {
		return nil, err
	}
	orders, err = m.getCorporationOrdersFromDB(corpID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	if len(orders) > 0 {
		return orders, nil
	}
	return m.getCorporationOrdersFromAPI(ctx, corpID)
}

func (m *MarketManager) getCorporationOrdersFromDB(corporationID int) ([]*MarketOrder, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  c.order_id
			, c.location_id
			, c.type_id
			, c.volume_entered
			, c.volume_remaining
			, c.min_volume
			, c.order_state
			, c.range
			, c.account_key
			, c.duration
			, c.escrow
			, c.price
			, c.bid
			, c.issued
			, c.character_id
			FROM app.market_orders c
			WHERE c.corporation_id = $1
			  AND c.fetched_at > (NOW() - INTERVAL '6 hours')`, corporationID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*MarketOrder
	for rs.Next() {
		o := &MarketOrder{}
		err = rs.Scan(
			&o.OrderID,
			&o.LocationID,
			&o.TypeID,
			&o.VolEntered,
			&o.VolRemaining,
			&o.MinVolume,
			&o.OrderState,
			&o.Range,
			&o.AccountKey,
			&o.Duration,
			&o.Escrow,
			&o.Price,
			&o.Bid,
			&o.Issued,
			&o.CharID,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func (m *MarketManager) getCorporationOrdersFromAPI(ctx context.Context, corpID int) ([]*MarketOrder, error) {
	orders, err := m.eveapi.GetCorporationOrders(ctx, corpID)
	if err != nil {
		return nil, err
	}
	var res []*MarketOrder
	for _, o := range orders {
		res = append(res, &MarketOrder{
			OrderID:      o.OrderID,
			CharID:       o.CharID,
			LocationID:   o.LocationID,
			TypeID:       o.TypeID,
			VolEntered:   o.VolEntered,
			VolRemaining: o.VolRemaining,
			MinVolume:    o.MinVolume,
			OrderState:   OrderState(o.OrderState),
			Range:        o.Range,
			AccountKey:   o.AccountKey,
			Duration:     o.Duration,
			Bid:          o.Bid,
			Escrow:       o.Escrow,
			Price:        o.Price,
			Issued:       o.Issued,
		})
	}
	return m.apiCorporationOrdersToDB(corpID, res)
}

func (m *MarketManager) apiCorporationOrdersToDB(corpID int, orders []*MarketOrder) ([]*MarketOrder, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	for _, o := range orders {
		_, err = db.Exec(
			`INSERT INTO app.market_orders
			 (
			     order_id
			   , location_id
			   , type_id
			   , volume_entered
			   , volume_remaining
			   , min_volume
			   , order_state
			   , range
			   , account_key
			   , duration
			   , escrow
			   , price
			   , bid
			   , issued
			   , corporation_id
			   , character_id
			 )
			 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
			 ON CONFLICT
			   ON CONSTRAINT "market_orders_pkey"
			 DO UPDATE
			   SET volume_remaining = EXCLUDED.volume_remaining
			     , order_state = EXCLUDED.order_state
			     , range = EXCLUDED.range
			     , escrow = EXCLUDED.escrow
			     , price = EXCLUDED.price
			     , fetched_at = DEFAULT`,
			o.OrderID,
			o.LocationID,
			o.TypeID,
			o.VolEntered,
			o.VolRemaining,
			o.MinVolume,
			o.OrderState,
			o.Range,
			o.AccountKey,
			o.Duration,
			o.Escrow,
			o.Price,
			o.Bid,
			o.Issued,
			corpID,
			o.CharID,
		)
		if err != nil {
			return nil, err
		}
	}
	return orders, nil
}
