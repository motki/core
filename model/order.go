package model

import (
	"time"

	"github.com/jackc/pgx"
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

	// loner is used to denote an order that was fetched alone.
	// loner orders will not be considered when fetching all of a corporation's orders.
	loner bool
}

func (m *Manager) GetCorporationOrder(ctx context.Context, corpID, orderID int) (*MarketOrder, error) {
	var err error
	if ctx, err = m.corporationAuthContext(ctx, corpID); err != nil {
		return nil, err
	}
	order, err := m.getCorporationOrderFromDB(corpID, orderID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	if order != nil {
		return order, nil
	}
	return m.getCorporationOrderFromAPI(ctx, corpID, orderID)
}

func (m *Manager) getCorporationOrderFromDB(corpID, orderID int) (*MarketOrder, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	r := c.QueryRow(
		`SELECT
			  c.order_id
			, c.station_id
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
			  AND c.order_id = $2
			  AND c.fetched_at > (NOW() - INTERVAL '6 hours')`, corpID, orderID)
	o := &MarketOrder{}
	bid := 0
	err = r.Scan(
		&o.OrderID,
		&o.StationID,
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
		&bid,
		&o.Issued,
		&o.CharID,
	)
	if err != nil {
		return nil, err
	}
	if bid > 0 {
		o.Bid = true
	}
	return o, nil
}

func (m *Manager) getCorporationOrderFromAPI(ctx context.Context, corpID, orderID int) (*MarketOrder, error) {
	o, err := m.eveapi.GetCorporationOrder(ctx, corpID, orderID)
	if err != nil {
		return nil, err
	}
	res := &MarketOrder{
		OrderID:      o.OrderID,
		CharID:       o.CharID,
		StationID:    o.StationID,
		TypeID:       o.TypeID,
		VolEntered:   o.VolEntered,
		VolRemaining: o.VolRemaining,
		MinVolume:    o.MinVolume,
		OrderState:   o.OrderState,
		Range:        o.Range,
		AccountKey:   o.AccountKey,
		Duration:     o.Duration,
		Bid:          o.Bid,
		Escrow:       o.Escrow,
		Price:        o.Price,
		Issued:       o.Issued,
		loner:        true,
	}
	_, err = m.apiCorporationOrdersToDB(corpID, []*MarketOrder{res})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *Manager) GetCorporationOrders(ctx context.Context, corpID int) (orders []*MarketOrder, err error) {
	if ctx, err = m.corporationAuthContext(ctx, corpID); err != nil {
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

func (m *Manager) getCorporationOrdersFromDB(corporationID int) ([]*MarketOrder, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  c.order_id
			, c.station_id
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
			  AND c.loner = 0
			  AND c.fetched_at > (NOW() - INTERVAL '6 hours')`, corporationID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*MarketOrder
	for rs.Next() {
		o := &MarketOrder{}
		bid := 0
		err = rs.Scan(
			&o.OrderID,
			&o.StationID,
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
			&bid,
			&o.Issued,
			&o.CharID,
		)
		if err != nil {
			return nil, err
		}
		if bid > 0 {
			o.Bid = true
		}
		res = append(res, o)
	}
	return res, nil
}

func (m *Manager) getCorporationOrdersFromAPI(ctx context.Context, corpID int) ([]*MarketOrder, error) {
	orders, err := m.eveapi.GetCorporationOrders(ctx, corpID)
	if err != nil {
		return nil, err
	}
	var res []*MarketOrder
	for _, o := range orders {
		res = append(res, &MarketOrder{
			OrderID:      o.OrderID,
			CharID:       o.CharID,
			StationID:    o.StationID,
			TypeID:       o.TypeID,
			VolEntered:   o.VolEntered,
			VolRemaining: o.VolRemaining,
			MinVolume:    o.MinVolume,
			OrderState:   o.OrderState,
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

func (m *Manager) apiCorporationOrdersToDB(corpID int, orders []*MarketOrder) ([]*MarketOrder, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	for _, o := range orders {
		bid := 0
		if o.Bid {
			bid = 1
		}
		loner := 0
		if o.loner {
			loner = 1
		}
		_, err = db.Exec(
			`INSERT INTO app.market_orders
			 (
			     order_id
			   , station_id
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
			   , loner
			 )
			 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
			 ON CONFLICT
			   ON CONSTRAINT "market_orders_pkey"
			 DO UPDATE
			   SET volume_remaining = EXCLUDED.volume_remaining
			     , order_state = EXCLUDED.order_state
			     , range = EXCLUDED.range
			     , escrow = EXCLUDED.escrow
			     , price = EXCLUDED.price
			     , bid = EXCLUDED.bid
			     , issued = EXCLUDED.issued
			     , loner = EXCLUDED.loner
			     , fetched_at = DEFAULT`,
			o.OrderID,
			o.StationID,
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
			bid,
			o.Issued,
			corpID,
			o.CharID,
			loner,
		)
		if err != nil {
			return nil, err
		}
	}
	return orders, nil
}
