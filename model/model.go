// Package model encapsulates the persistence layer of the MOTKI application.
package model // import "github.com/motki/core/model"

import (
	"github.com/motki/core/db"
	"github.com/motki/core/eveapi"
	"github.com/motki/core/evedb"
	"github.com/motki/core/evemarketer"
	"github.com/motki/core/log"
)

// A bootstrap contains the core, shared dependencies.
type bootstrap struct {
	pool   *db.ConnPool
	evedb  *evedb.EveDB
	eveapi *eveapi.EveAPI
	ec     *evemarketer.EveMarketer
}

// A Manager handles loading and saving of data.
//
// Most data is stored in the configured database and only
// fetched from external sites when necessary.
type Manager struct {
	*AssetManager
	*BlueprintManager
	*CharacterManager
	*CorpManager
	*IndustryManager
	*InventoryManager
	*LocationManager
	*MailManager
	*MarketManager
	*ProductManager
	*StructureManager
	*UserManager

	noexport struct{}
}

// NewManager creates a new Manager, ready for use.
func NewManager(pool *db.ConnPool, evedb *evedb.EveDB, api *eveapi.EveAPI, ec *evemarketer.EveMarketer) *Manager {
	m := bootstrap{pool: pool, evedb: evedb, eveapi: api, ec: ec}

	char := newCharacterManager(m)
	user := newUserManager(m, char)
	corp := newCorpManager(m, user)
	asset := newAssetManager(m, corp)
	market := newMarketManager(m, corp)
	structure := newStructureManager(m, corp)

	return &Manager{
		AssetManager:     asset,
		BlueprintManager: newBlueprintManager(m, corp),
		CharacterManager: char,
		CorpManager:      corp,
		IndustryManager:  newIndustryManager(m, corp),
		InventoryManager: newInventoryManager(m, corp, asset),
		LocationManager:  newLocationManager(m, asset, structure),
		MailManager:      newMailManager(m),
		MarketManager:    market,
		ProductManager:   newProductManager(m, corp, market),
		StructureManager: structure,
		UserManager:      user,
	}
}

// UpdateCorporationData fetches updated data for all opted-in corporations.
//
// The function returned by this method is intended to be invoke in regular intervals.
func (m *Manager) UpdateCorporationDataFunc(logger log.Logger) func() error {
	return func() error {
		corps, err := m.GetCorporationsOptedIn()
		if err != nil {
			return err
		}
		if len(corps) == 0 {
			logger.Debugf("no corporations opted in, not updating corp data")
			return nil
		}
		for _, corpID := range corps {
			logger.Debugf("updating data for corp %d", corpID)
			a, err := m.GetCorporationAuthorization(corpID)
			if err != nil {
				logger.Errorf("error getting corp auth: %s", err.Error())
				continue
			}

			ctx := a.Context()
			if _, err := m.FetchCorporationDetail(ctx); err != nil {
				logger.Errorf("error fetching corp details: %s", err.Error())
			}
			if res, err := m.GetCorporationAssets(ctx, a.CorporationID); err != nil {
				logger.Errorf("error fetching corp assets: %s", err.Error())
			} else {
				logger.Debugf("fetched %d assets for corporation %d", len(res), a.CorporationID)
			}

			if res, err := m.GetCorporationOrders(ctx, a.CorporationID); err != nil {
				logger.Errorf("error fetching corp orders: %s", err.Error())
			} else {
				logger.Debugf("fetched %d orders for corporation %d", len(res), a.CorporationID)
			}

			if res, err := m.GetCorporationBlueprints(ctx, a.CorporationID); err != nil {
				logger.Errorf("error fetching corp blueprints: %s", err.Error())
			} else {
				logger.Debugf("fetched %d blueprints for corporation %d", len(res), a.CorporationID)
			}

			if res, err := m.GetCorporationStructures(ctx, a.CorporationID); err != nil {
				logger.Errorf("error fetching corp structures: %s", err.Error())
			} else {
				logger.Debugf("fetched %d structures for corporation %d", len(res), a.CorporationID)
			}
		}
		return nil
	}
}
