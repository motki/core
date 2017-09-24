package model

import (
	"time"

	"github.com/jackc/pgx"
	"github.com/motki/motki/eveapi"
)

type Character struct {
	CharacterID   int
	Name          string
	BloodlineID   int
	RaceID        int
	AncestryID    int
	CorporationID int
	AllianceID    int
	BirthDate     time.Time
	Description   string
}

func (m *Manager) GetCharacter(characterID int) (*Character, error) {
	c, err := m.getCharacterFromDB(characterID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return m.getCharacterFromAPI(characterID)
	}
	return c, nil
}

func (m *Manager) getCharacterFromDB(characterID int) (*Character, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	r := c.QueryRow(
		`SELECT
			  c.character_id
			, c.name
			, c.bloodline_id
			, c.race_id
			, c.ancestry_id
			, c.corporation_id
			, c.alliance_id
			, c.birth_date
			, c.description
			FROM app.characters c
			WHERE c.character_id = $1
				AND c.fetched_at > NOW() - INTERVAL '12 hours'`, characterID)
	char := &Character{}
	err = r.Scan(
		&char.CharacterID,
		&char.Name,
		&char.BloodlineID,
		&char.RaceID,
		&char.AncestryID,
		&char.CorporationID,
		&char.AllianceID,
		&char.BirthDate,
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

func (m *Manager) getCharacterFromAPI(characterID int) (*Character, error) {
	char, err := m.eveapi.GetCharacter(characterID)
	if err != nil {
		return nil, err
	}
	return m.apiCharacterToDB(char)
}

func (m *Manager) apiCharacterToDB(char *eveapi.Character) (*Character, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	c := &Character{
		CharacterID:   char.CharacterID,
		Name:          char.Name,
		BloodlineID:   char.BloodlineID,
		RaceID:        char.RaceID,
		AncestryID:    char.AncestryID,
		CorporationID: char.CorporationID,
		AllianceID:    char.AllianceID,
		BirthDate:     char.BirthDate,
		Description:   char.Description,
	}
	_, err = db.Exec(
		`INSERT INTO app.characters
				(character_id, name, bloodline_id, race_id, ancestry_id
				, corporation_id, alliance_id, birth_date, description)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
				ON CONFLICT ON CONSTRAINT "characters_pkey" DO
				UPDATE SET name = EXCLUDED.name
					, corporation_id = EXCLUDED.corporation_id
					, alliance_id = EXCLUDED.alliance_id
					, description = EXCLUDED.description
					, fetched_at = DEFAULT`,
		c.CharacterID,
		c.Name,
		c.BloodlineID,
		c.RaceID,
		c.AncestryID,
		c.CorporationID,
		c.AllianceID,
		c.BirthDate,
		c.Description,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}
