package model

import (
	"github.com/tyler-sommer/motki/db"
	"github.com/tyler-sommer/motki/eveapi"
	"github.com/tyler-sommer/motki/evedb"
)

type Manager struct {
	pool   *db.ConnPool
	evedb  *evedb.EveDB
	eveapi *eveapi.EveAPI
}

func NewManager(pool *db.ConnPool, evedb *evedb.EveDB, api *eveapi.EveAPI) *Manager {
	return &Manager{pool: pool, evedb: evedb, eveapi: api}
}
