// Package evedb manages interaction with the EVE Static Dump.
//
// This package is intended to abstract access to the various different
// tables and assets provided in the EVE Static Dump.
package evedb // import "github.com/motki/core/evedb"

import (
	"github.com/motki/core/db"
)

// EveDB is the central service for accessing all EVE Static Dump data.
type EveDB struct {
	pool *db.ConnPool
}

// New creates a new EveDB using the given connection pool.
func New(p *db.ConnPool) *EveDB {
	return &EveDB{pool: p}
}
