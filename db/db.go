// Package db manages interaction with an underlying database store.
package db

import (
	"database/sql"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
)

// Config represents configuration details for creating a connection pool.
type Config struct {
	ConnString     string `toml:"connection_string"`
	MaxConnections int    `toml:"max_connections"`
}

// New creates a new ConnPool using the given Config.
//
// Generally only one ConnPool is needed, shared by the entire application.
func New(c Config) (*ConnPool, error) {
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
func (p *ConnPool) Open() (*sql.DB, error) {
	return stdlib.OpenFromConnPool(p.pool)
}
