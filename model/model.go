package model

import (
	"github.com/motki/motkid/db"
	"github.com/motki/motkid/eveapi"
	"github.com/motki/motkid/evecentral"
	"github.com/motki/motkid/evedb"
)

type Manager struct {
	pool   *db.ConnPool
	evedb  *evedb.EveDB
	eveapi *eveapi.EveAPI
	ec     *evecentral.EveCentral
}

func NewManager(pool *db.ConnPool, evedb *evedb.EveDB, api *eveapi.EveAPI, ec *evecentral.EveCentral) *Manager {
	return &Manager{pool: pool, evedb: evedb, eveapi: api, ec: ec}
}
