package model

import (
	"time"

	"github.com/jackc/pgx"
	"github.com/motki/motki/eveapi"
)

type Corporation struct {
	CorporationID int
	Name          string
	AllianceID    int
	CreationDate  time.Time
	Description   string
	Ticker        string
}

func (m *Manager) GetCorporation(corporationID int) (*Corporation, error) {
	c, err := m.getCorporationFromDB(corporationID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return m.getCorporationFromAPI(corporationID)
	}
	return c, nil
}

func (m *Manager) getCorporationFromDB(corporationID int) (*Corporation, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	r := c.QueryRow(
		`SELECT
			  c.corporation_id
			, c.name
			, c.alliance_id
			, c.creation_date
			, c.ticker
			, c.description
			FROM app.corporations c
			WHERE c.corporation_id = $1
				AND c.fetched_at > NOW() - INTERVAL '7 days'`, corporationID)
	char := &Corporation{}
	err = r.Scan(
		&char.CorporationID,
		&char.Name,
		&char.AllianceID,
		&char.CreationDate,
		&char.Ticker,
		&char.Description,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return char, nil
}

func (m *Manager) getCorporationFromAPI(corporationID int) (*Corporation, error) {
	char, err := m.eveapi.GetCorporation(corporationID)
	if err != nil {
		return nil, err
	}
	return m.apiCorporationToDB(char)
}

func (m *Manager) apiCorporationToDB(corp *eveapi.Corporation) (*Corporation, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	c := &Corporation{
		CorporationID: corp.CorporationID,
		Name:          corp.Name,
		AllianceID:    corp.AllianceID,
		CreationDate:  corp.CreationDate,
		Ticker:        corp.Ticker,
		Description:   corp.Description,
	}
	_, err = db.Exec(
		`INSERT INTO app.corporations
				(corporation_id, name, alliance_id, creation_date, ticker, description)
				VALUES($1, $2, $3, $4, $5, $6)
				ON CONFLICT ON CONSTRAINT "corporations_pkey" DO
				UPDATE SET name = EXCLUDED.name
					, alliance_id = EXCLUDED.alliance_id
					, ticker = EXCLUDED.ticker
					, description = EXCLUDED.ticker
					, fetched_at = DEFAULT`,
		c.CorporationID,
		c.Name,
		c.AllianceID,
		c.CreationDate,
		c.Ticker,
		c.Description,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}
