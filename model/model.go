// Package model encapsulates the persistence layer.
package model

import (
	"github.com/motki/motki/db"
	"github.com/motki/motki/eveapi"
	"github.com/motki/motki/evecentral"
	"github.com/motki/motki/evedb"
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
