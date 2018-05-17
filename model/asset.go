package model

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/core/eveapi"
	"github.com/motki/core/evedb"
)

type Asset struct {
	ItemID       int    `json:"item_id"`
	LocationID   int    `json:"location_id"`
	LocationType string `json:"location_type"`
	LocationFlag string `json:"location_flag"`
	TypeID       int    `json:"type_id"`
	Quantity     int    `json:"quantity"`
	Singleton    bool   `json:"singleton"`

	corpID    int
	fetchedAt time.Time
}

func assetFromEveAPI(a *eveapi.Asset) *Asset {
	return &Asset{
		ItemID:       a.ItemID,
		LocationID:   a.LocationID,
		TypeID:       a.TypeID,
		LocationFlag: a.LocationFlag,
		LocationType: a.LocationType,
		Quantity:     a.Quantity,
		Singleton:    a.Singleton,
	}
}

type AssetManager struct {
	bootstrap

	corp *CorpManager
}

func newAssetManager(m bootstrap, corp *CorpManager) *AssetManager {
	return &AssetManager{m, corp}
}

func (m *AssetManager) GetCorporationAssets(ctx context.Context, corpID int) (res []*Asset, err error) {
	if ctx, err = m.corp.authContext(ctx, corpID); err != nil {
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

// GetCorporationAssetsByTypeAndLocationID queries the database to find any assets
// with the given type and location.
//
// This method will not fetch assets from the API.
func (m *AssetManager) GetCorporationAssetsByTypeAndLocationID(ctx context.Context, corpID, typeID, locationID int) (res []*Asset, err error) {
	if _, err = m.corp.authContext(ctx, corpID); err != nil {
		return nil, err
	}
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  a.item_id
			, a.location_id
			, a.location_type
			, a.location_flag
			, a.type_id
			, a.quantity
			, a.singleton
			, a.corporation_id
			, a.fetched_at
			FROM app.assets a
			JOIN app.assets a2 ON a.location_id = a2.item_id 
				AND a2.valid = TRUE AND a2.corporation_id = a.corporation_id
			WHERE a.type_id = $2 
				AND a.corporation_id = $1 
				AND a.valid = true 
				AND (a.location_id = $3 OR a2.location_id = $3)`, corpID, typeID, locationID)
	if err != nil {
		return nil, err
	}
	for rs.Next() {
		r := &Asset{}
		err := rs.Scan(
			&r.ItemID,
			&r.LocationID,
			&r.LocationType,
			&r.LocationFlag,
			&r.TypeID,
			&r.Quantity,
			&r.Singleton,
			&r.corpID,
			&r.fetchedAt,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func (m *AssetManager) GetCorporationAsset(ctx context.Context, corpID int, itemID int) (res *Asset, err error) {
	if ctx, err = m.corp.authContext(ctx, corpID); err != nil {
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

func (m *AssetManager) GetAssetSystem(a *Asset) (*evedb.System, error) {
	switch {
	case a.LocationID < 60000000:
		// LocationID is a SystemID
		return m.evedb.GetSystem(a.LocationID)
	case a.LocationID < 66000000:
		// LocationID is a LocationID
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

func (m *AssetManager) getCorporationAssetFromDB(corpID int, itemID int) (*Asset, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs := c.QueryRow(
		`SELECT
			  a.item_id
			, a.location_id
			, a.location_type
			, a.location_flag
			, a.type_id
			, a.quantity
			, a.singleton
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
		&r.LocationType,
		&r.LocationFlag,
		&r.TypeID,
		&r.Quantity,
		&r.Singleton,
		&r.fetchedAt,
		&status,
		&valid,
		&r.corpID,
	)
	if err != nil {
		return nil, err
	}
	if !status {
		//	return nil, errors.New("stale database entry")
	}
	if !valid {
		//return nil, errors.New("invalid database entry")
	}
	return r, nil
}

func (m *AssetManager) getCorporationAssetsFromDB(corpID int) ([]*Asset, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  a.item_id
			, a.location_id
			, a.location_type
			, a.location_flag
			, a.type_id
			, a.quantity
			, a.singleton
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
			&r.LocationType,
			&r.LocationFlag,
			&r.TypeID,
			&r.Quantity,
			&r.Singleton,
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

func (m *AssetManager) getCorporationAssetsFromAPI(ctx context.Context, corpID int) ([]*Asset, error) {
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

func (m *AssetManager) apiCorporationAssetsToDB(corpID int, bps []*Asset) ([]*Asset, error) {
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
					(corporation_id, character_id, item_id, location_id, type_id, quantity, singleton, location_type, location_flag, fetched_at, valid)
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, DEFAULT, DEFAULT)`,
			corpID,
			0,
			bp.ItemID,
			bp.LocationID,
			bp.TypeID,
			bp.Quantity,
			bp.Singleton,
			bp.LocationType,
			bp.LocationFlag,
		)
		if err != nil {
			return nil, err
		}
	}
	return bps, nil
}
