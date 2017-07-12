// Package model encapsulates the persistence layer.
package model

import (
	"github.com/motki/motkid/db"
	"github.com/motki/motkid/eveapi"
	"github.com/motki/motkid/evecentral"
	"github.com/motki/motkid/evedb"
)

// A Manager is used to retrieve and save data.
//
// Most data is stored in the configured database and only
// fetched from external sites when necessary.
type Manager struct {
	pool   *db.ConnPool
	evedb  *evedb.EveDB
	eveapi *eveapi.EveAPI
	ec     *evecentral.EveCentral
}

// NewManager creates a new Manager, ready for use.
func NewManager(pool *db.ConnPool, evedb *evedb.EveDB, api *eveapi.EveAPI, ec *evecentral.EveCentral) *Manager {
	return &Manager{pool: pool, evedb: evedb, eveapi: api, ec: ec}
}
