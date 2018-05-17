package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/core/cache"
	"github.com/motki/core/eveapi"
)

type CorpManager struct {
	bootstrap

	user *UserManager
	char *CharacterManager

	cache *cache.Bucket
}

func newCorpManager(m bootstrap, user *UserManager, char *CharacterManager) *CorpManager {
	return &CorpManager{m, user, char, cache.New(10 * time.Second)}
}

type Corporation struct {
	CorporationID int       `json:"corporation_id"`
	Name          string    `json:"name"`
	AllianceID    int       `json:"alliance_id"`
	CreationDate  time.Time `json:"creation_date"`
	Description   string    `json:"description"`
	Ticker        string    `json:"ticker"`
}

func (m *CorpManager) GetCorporation(corporationID int) (*Corporation, error) {
	c, err := m.getCorporationFromDB(corporationID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return m.getCorporationFromAPI(corporationID)
	}
	return c, nil
}

func (m *CorpManager) getCorporationFromDB(corporationID int) (*Corporation, error) {
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

func (m *CorpManager) getCorporationFromAPI(corporationID int) (*Corporation, error) {
	char, err := m.eveapi.GetCorporation(corporationID)
	if err != nil {
		return nil, err
	}
	return m.apiCorporationToDB(char)
}

func (m *CorpManager) apiCorporationToDB(corp *eveapi.Corporation) (*Corporation, error) {
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
			ON CONFLICT ON CONSTRAINT "corporations_pkey"
			DO UPDATE
			SET name = EXCLUDED.name
			  , alliance_id = EXCLUDED.alliance_id
			  , ticker = EXCLUDED.ticker
			  , description = EXCLUDED.description
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

// Divisions is a map of division key to division name.
type Divisions map[string]string

func (d Divisions) GetName(idx int) (string, bool) {
	id := strconv.Itoa(idx)
	if v, ok := d[id]; ok {
		return v, true
	}
	return "", false
}

func (d Divisions) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Divisions) Scan(src interface{}) error {
	s, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("invalid value for division: %v", src)
	}
	return json.Unmarshal(s, &d)
}

type CorporationDetail struct {
	CorporationID int    `json:"corporation_id"`
	CEOID         int    `json:"ceoid"`
	CEOName       string `json:"ceo_name"`
	StationID     int    `json:"station_id"`
	StationName   string `json:"station_name"`
	FactionID     int    `json:"faction_id"`
	MemberCount   int    `json:"member_count"`
	Shares        int    `json:"shares"`

	Wallets Divisions `json:"wallets"`
	Hangars Divisions `json:"hangars"`
}

func (m *CorpManager) GetCorporationDetail(corpID int) (*CorporationDetail, error) {
	d, err := m.getCorporationDetailFromDB(corpID)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, ErrCorpNotRegistered
	}
	return d, nil
}

func (m *CorpManager) FetchCorporationDetail(ctx context.Context) (*CorporationDetail, error) {
	return m.getCorporationDetailFromAPI(ctx)
}

func (m *CorpManager) getCorporationDetailFromDB(corporationID int) (*CorporationDetail, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	r := c.QueryRow(
		`SELECT
			  c.corporation_id
			, c.ceo_id
			, c.ceo_name
			, c.hq_station_id
			, c.hq_station_name
			, c.faction_id
			, c.member_count
			, c.shares
			, c.hangars
			, c.divisions
			FROM app.corporation_details c
			WHERE c.corporation_id = $1`, corporationID)
	corp := &CorporationDetail{}
	err = r.Scan(
		&corp.CorporationID,
		&corp.CEOID,
		&corp.CEOName,
		&corp.StationID,
		&corp.StationName,
		&corp.FactionID,
		&corp.MemberCount,
		&corp.Shares,
		&corp.Hangars,
		&corp.Wallets,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return corp, nil
}

func (m *CorpManager) getCorporationDetailFromAPI(ctx context.Context) (*CorporationDetail, error) {
	corpID, ok := m.corpID(ctx)
	if !ok {
		return nil, errors.New("unable to get corpID from context")
	}
	actx, err := m.authContext(ctx, corpID)
	if err != nil {
		return nil, err
	}
	sheet, err := m.eveapi.GetCorporationSheet(actx, corpID)
	if err != nil {
		return nil, err
	}
	ceo, err := m.char.GetCharacter(sheet.CEOID)
	if err != nil {
		return nil, err
	}
	sta, err := m.evedb.GetStation(sheet.StationID)
	if err != nil {
		return nil, err
	}
	return m.apiCorporationDetailToDB(&CorporationDetail{
		CorporationID: sheet.CorporationID,
		CEOID:         sheet.CEOID,
		CEOName:       ceo.Name,
		StationID:     sheet.StationID,
		StationName:   sta.Name,
		FactionID:     sheet.FactionID,
		MemberCount:   sheet.MemberCount,
		Shares:        sheet.Shares,
		Hangars:       Divisions(sheet.Hangars),
		Wallets:       Divisions(sheet.Wallets),
	})
}

func (m *CorpManager) apiCorporationDetailToDB(detail *CorporationDetail) (*CorporationDetail, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	_, err = db.Exec(
		`INSERT INTO app.corporation_details
			 (
			     corporation_id
			   , ceo_id
			   , ceo_name
			   , hq_station_id
			   , hq_station_name
			   , faction_id
			   , member_count
			   , shares
			   , hangars
			   , divisions
			 )
			 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			 ON CONFLICT
			   ON CONSTRAINT "corporation_details_pkey"
			 DO UPDATE
			   SET ceo_id = EXCLUDED.ceo_id
			     , ceo_name = EXCLUDED.ceo_name
			     , hq_station_id = EXCLUDED.hq_station_id
			     , hq_station_name = EXCLUDED.hq_station_name
			     , faction_id = EXCLUDED.faction_id
			     , member_count = EXCLUDED.member_count
			     , shares = EXCLUDED.shares
			     , hangars = EXCLUDED.hangars
			     , divisions = EXCLUDED.divisions
			     , fetched_at = DEFAULT`,
		detail.CorporationID,
		detail.CEOID,
		detail.CEOName,
		detail.StationID,
		detail.StationName,
		detail.FactionID,
		detail.MemberCount,
		detail.Shares,
		detail.Hangars,
		detail.Wallets,
	)
	if err != nil {
		return nil, err
	}
	return detail, nil
}
