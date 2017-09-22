// Package db manages interaction with an underlying database store.
//
// This package is designed to abstract the actual library used for interacting
// with the database. The exported API returns standard database/sql
// structures.
package db

import (
	"database/sql"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/motki/motki/log"
)

// Config represents configuration details for creating a connection pool.
type Config struct {
	ConnString     string `toml:"connection_string"`
	MaxConnections int    `toml:"max_connections"`
}

// New creates a new ConnPool using the given Config.
//
// Generally only one ConnPool should be used, shared by the entire
// application.
func New(c Config, l log.Logger) (*ConnPool, error) {
	l.Debugf("db: init database connection pool for: %s", c.ConnString)
	l.Debugf("db: max connections: %d", c.MaxConnections)
	pcon, err := pgx.ParseConnectionString(c.ConnString)
	if err != nil {
		return nil, err
	}
	p, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: pcon, MaxConnections: c.MaxConnections})
	if err != nil {
		return nil, err
	}
	return &ConnPool{p}, nil
}

// Type ConnPool represents a connection pool.
type ConnPool struct {
	pool *pgx.ConnPool
}

// Open acquires a connection for the caller.
//
// This method is designed to allow compatibility with the database/sql
// package.
//
// Remember to Close the DB when you're done with it! Otherwise the connection
// will not be released back to the pool, and eventually you will run out of
// available connections.
func (p *ConnPool) Open() (*sql.DB, error) {
	return stdlib.OpenFromConnPool(p.pool)
}
