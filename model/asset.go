package model

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/core/eveapi"
	"github.com/motki/core/evedb"
)

type Asset struct {
	ItemID      int
	LocationID  int
	TypeID      int
	Quantity    int
	FlagID      int
	Singleton   bool
	RawQuantity int

	corpID    int
	fetchedAt time.Time
}

func assetFromEveAPI(bp *eveapi.Asset) *Asset {
	return &Asset{
		ItemID:      bp.ItemID,
		LocationID:  bp.LocationID,
		TypeID:      bp.TypeID,
		FlagID:      bp.FlagID,
		Quantity:    bp.Quantity,
		Singleton:   bp.Singleton,
		RawQuantity: bp.RawQuantity,
	}
}

func (m *Manager) GetCorporationAssets(ctx context.Context, corpID int) (res []*Asset, err error) {
	if ctx, err = m.corporationAuthContext(ctx, corpID); err != nil {
		return nil, err
	}
	res, err = m.getCorporationAssetsFromDB(corpID)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	return m.getCorporationAssetsFromAPI(ctx, corpID)
}

func (m *Manager) GetCorporationAssetsByTypeAndLocationID(ctx context.Context, corpID, typeID, locationID int) (res []*Asset, err error) {
	if ctx, err = m.corporationAuthContext(ctx, corpID); err != nil {
		return nil, err
	}
	assets, err := m.GetCorporationAssets(ctx, corpID)
	if err != nil {
		return nil, err
	}
	for _, a := range assets {
		if a.TypeID == typeID && a.LocationID == locationID {
			res = append(res, a)
		}
	}
	return res, nil
}

func (m *Manager) GetCorporationAsset(ctx context.Context, corpID int, itemID int) (res *Asset, err error) {
	if ctx, err = m.corporationAuthContext(ctx, corpID); err != nil {
		return nil, err
	}
	res, err = m.getCorporationAssetFromDB(corpID, itemID)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	assets, err := m.getCorporationAssetsFromAPI(ctx, corpID)
	if err != nil {
		return nil, err
	}
	for _, a := range assets {
		if a.ItemID == itemID {
			return res, nil
		}
	}
	return nil, errors.New("unable to find asset")
}

func (m *Manager) GetAssetSystem(a *Asset) (*evedb.System, error) {
	switch {
	case a.LocationID < 60000000:
		// LocationID is a SystemID
		return m.evedb.GetSystem(a.LocationID)
	case a.LocationID < 66000000:
		// LocationID is a StationID
	case a.LocationID < 67000000:
		// LocationID is a conquerable station or outpost

	default:
		// LocationID is in a container (or citadel)
		ca, err := m.getCorporationAssetFromDB(a.corpID, a.LocationID)
		if err != nil {
			// TODO: it may be a public Citadel, so query the eveapi for the structure.
			return nil, err
		}
		return m.GetAssetSystem(ca)
	}
	return nil, errors.New("unable to find system for asset")
}

func (m *Manager) getCorporationAssetFromDB(corpID int, itemID int) (*Asset, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs := c.QueryRow(
		`SELECT
			  a.item_id
			, a.location_id
			, a.type_id
			, a.quantity
			, a.singleton
			, a.raw_quantity
			, a.flag_id
			, a.fetched_at
			, (a.fetched_at > (NOW() - INTERVAL '12 hours')) status
			, (a.valid = TRUE) validity
			, a.corporation_id
			FROM app.assets a
			WHERE a.corporation_id = $1
			  AND a.item_id = $2
			ORDER BY a.fetched_at ASC
			LIMIT 1`, corpID, itemID)
	if err != nil {
		return nil, err
	}
	r := &Asset{}
	var status bool
	var valid bool
	err = rs.Scan(
		&r.ItemID,
		&r.LocationID,
		&r.TypeID,
		&r.Quantity,
		&r.Singleton,
		&r.RawQuantity,
		&r.FlagID,
		&r.fetchedAt,
		&status,
		&valid,
		&r.corpID,
	)
	if err != nil {
		return nil, err
	}
	if !status {
		return nil, errors.New("stale database entry")
	}
	if !valid {
		return nil, errors.New("invalid database entry")
	}
	return r, nil
}

func (m *Manager) getCorporationAssetsFromDB(corpID int) ([]*Asset, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  a.item_id
			, a.location_id
			, a.type_id
			, a.quantity
			, a.singleton
			, a.raw_quantity
			, a.flag_id
			, a.corporation_id
			, a.fetched_at
			FROM app.assets a
			WHERE a.corporation_id = $1
			  AND a.valid = TRUE
			  AND a.fetched_at > (NOW() - INTERVAL '7 hours')`, corpID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*Asset
	for rs.Next() {
		r := &Asset{}
		err := rs.Scan(
			&r.ItemID,
			&r.LocationID,
			&r.TypeID,
			&r.Quantity,
			&r.Singleton,
			&r.RawQuantity,
			&r.FlagID,
			&r.corpID,
			&r.fetchedAt,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (m *Manager) getCorporationAssetsFromAPI(ctx context.Context, corpID int) ([]*Asset, error) {
	bps, err := m.eveapi.GetCorporationAssets(ctx, corpID)
	if err != nil {
		return nil, err
	}
	var res []*Asset
	for _, bp := range bps {
		res = append(res, assetFromEveAPI(bp))
	}
	return m.apiCorporationAssetsToDB(corpID, res)
}

func (m *Manager) apiCorporationAssetsToDB(corpID int, bps []*Asset) ([]*Asset, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	_, err = db.Exec(`UPDATE app.assets SET valid = FALSE WHERE corporation_id = $1`, corpID)
	if err != nil {
		return nil, err
	}
	for _, bp := range bps {
		_, err = db.Exec(
			`INSERT INTO app.assets
					(corporation_id, character_id, item_id, location_id, type_id, quantity, singleton, raw_quantity, flag_id, fetched_at, valid)
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, DEFAULT, DEFAULT)`,
			corpID,
			0,
			bp.ItemID,
			bp.LocationID,
			bp.TypeID,
			bp.Quantity,
			bp.Singleton,
			bp.RawQuantity,
			bp.FlagID,
		)
		if err != nil {
			return nil, err
		}
	}
	return bps, nil
}
