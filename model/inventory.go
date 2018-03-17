package model

import (
	"time"

	"github.com/jackc/pgx"
	"golang.org/x/net/context"
)

type InventoryItem struct {
	TypeID        int
	LocationID    int
	MinimumLevel  int
	CurrentLevel  int
	CorporationID int
	FetchedAt     time.Time
}

func (m *Manager) GetCorporationInventory(ctx context.Context, corpID int) (items []*InventoryItem, err error) {
	if ctx, err = m.corporationAuthContext(ctx, corpID); err != nil {
		return nil, err
	}
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  c.type_id
			, c.location_id
			, c.min_level
			, c.curr_level
			, c.fetched_at
			FROM app.inventory_items c
			WHERE c.corporation_id = $1`, corpID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	for rs.Next() {
		r := &InventoryItem{}
		err := rs.Scan(
			&r.TypeID,
			&r.LocationID,
			&r.MinimumLevel,
			&r.CurrentLevel,
			&r.FetchedAt,
		)

		if err != nil {
			return nil, err
		}
		r.CorporationID = corpID
		if r.FetchedAt.Before(time.Now().Add(-2 * time.Hour)) {
			err = m.updateInventoryItemLevel(ctx, r)
			if err != nil {
				return nil, err
			}
			err = m.SaveInventoryItem(ctx, r)
			if err != nil {
				return nil, err
			}
		}
		items = append(items, r)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return items, nil
}

func (m *Manager) updateInventoryItemLevel(ctx context.Context, item *InventoryItem) error {
	assets, err := m.GetCorporationAssetsByTypeAndLocationID(ctx, item.CorporationID, item.TypeID, item.LocationID)
	if err != nil {
		return err
	}
	item.CurrentLevel = 0
	oldest := time.Now()
	for _, a := range assets {
		item.CurrentLevel += a.Quantity
		if a.fetchedAt.Before(oldest) {
			oldest = a.fetchedAt
		}
	}
	item.FetchedAt = oldest
	return nil
}

func (m *Manager) NewInventoryItem(ctx context.Context, corpID, typeID, locationID int) (*InventoryItem, error) {
	var err error
	if ctx, err = m.corporationAuthContext(ctx, corpID); err != nil {
		return nil, err
	}
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	it := &InventoryItem{CorporationID: corpID}
	r := c.QueryRow(
		`SELECT
			  c.type_id
			, c.location_id
			, c.min_level
			, c.curr_level
			, c.fetched_at
			FROM app.inventory_items c
			WHERE c.corporation_id = $1 AND c.type_id = $2 AND c.location_id = $3`, corpID, typeID, locationID)
	if err := r.Scan(&it.TypeID, &it.LocationID, &it.MinimumLevel, &it.CurrentLevel, &it.FetchedAt); err != nil {
		if err == pgx.ErrNoRows {
			it = &InventoryItem{
				TypeID:     typeID,
				LocationID: locationID,

				CorporationID: corpID,
			}
		} else {
			return nil, err
		}
	}
	if err = m.updateInventoryItemLevel(ctx, it); err != nil {
		return nil, err
	}
	err = m.SaveInventoryItem(ctx, it)
	if err != nil {
		return nil, err
	}
	return it, nil
}

func (m *Manager) SaveInventoryItem(ctx context.Context, item *InventoryItem) error {
	var err error
	if ctx, err = m.corporationAuthContext(ctx, item.CorporationID); err != nil {
		return err
	}
	c, err := m.pool.Open()
	if err != nil {
		return err
	}
	defer m.pool.Release(c)
	_, err = c.Exec(`INSERT INTO app.inventory_items (
		type_id,
		location_id,
		curr_level,
		min_level,
		fetched_at,
		corporation_id)
	VALUES($1, $2, $3, $4, $5, $6)
	ON CONFLICT ON CONSTRAINT "inventory_items_pkey"
		 DO UPDATE SET curr_level = EXCLUDED.curr_level,
		     min_level = EXCLUDED.min_level,
		     fetched_at = EXCLUDED.fetched_at`,
		item.TypeID,
		item.LocationID,
		item.CurrentLevel,
		item.MinimumLevel,
		item.FetchedAt,
		item.CorporationID)
	return err
}
