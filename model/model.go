// Package model encapsulates the persistence layer of the MOTKI application.
package model // import "github.com/motki/core/model"

import (
	"github.com/motki/core/db"
	"github.com/motki/core/eveapi"
	"github.com/motki/core/evedb"
	"github.com/motki/core/evemarketer"
)

// A Manager is used to retrieve and save data.
//
// Most data is stored in the configured database and only
// fetched from external sites when necessary.
type Manager struct {
	pool   *db.ConnPool
	evedb  *evedb.EveDB
	eveapi *eveapi.EveAPI
	ec     *evemarketer.EveMarketer
}

// NewManager creates a new Manager, ready for use.
func NewManager(pool *db.ConnPool, evedb *evedb.EveDB, api *eveapi.EveAPI, ec *evemarketer.EveMarketer) *Manager {
	return &Manager{pool: pool, evedb: evedb, eveapi: api, ec: ec}
}
