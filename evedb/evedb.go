// Package evedb manages interaction with the EVE Static Dump.
package evedb

import (
	"github.com/tyler-sommer/motki/db"
)

type EveDB struct {
	pool *db.ConnPool
}

func New(p *db.ConnPool) *EveDB {
	return &EveDB{pool: p}
}

func (e *EveDB) GetRaces() ([]string, error) {
	c, err := e.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	// TODO: Super unhappy about formatting
	rs, err := c.Query(
		`SELECT
			  "raceID"
			, "raceName"
--			, description
--			, iconID
--			, shortDescription
			FROM evesde."chrRaces"`)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	res := []string{}
	var id int
	var name string
	for rs.Next() {
		err := rs.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		res = append(res, name)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
