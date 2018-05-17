package model

import (
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/motki/core/cache"
)

var ErrCorpNotRegistered = errors.New("ceo or director is not registered for the given corporation")

type CorporationConfig struct {
	OptIn     bool      `json:"opt_in"`
	OptInBy   int       `json:"opt_in_by"`
	OptInDate time.Time `json:"opt_in_date"`

	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *CorpManager) GetCorporationsOptedIn() ([]int, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	rs, err := c.Query(
		`SELECT
			  c.corporation_id
			FROM app.corporation_settings c
			WHERE c.opted_in = TRUE`)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []int
	for rs.Next() {
		i := 0
		if err = rs.Scan(&i); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func (m *CorpManager) GetCorporationConfig(corpID int) (*CorporationConfig, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	r := c.QueryRow(
		`SELECT
			  c.opted_in
			, c.opted_in_by
			, c.opted_in_at
			, c.created_by
			, c.created_at
			FROM app.corporation_settings c
			WHERE c.corporation_id = $1`, corpID)
	corp := &CorporationConfig{}
	err = r.Scan(
		&corp.OptIn,
		&corp.OptInBy,
		&corp.OptInDate,
		&corp.CreatedBy,
		&corp.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCorpNotRegistered
		}
		return nil, err
	}
	return corp, nil
}

func (m *CorpManager) GetCorporationAuthorization(corpID int) (*Authorization, error) {
	v, err := m.cache.Memoize("corp_auth:"+strconv.Itoa(corpID), func() (cache.Value, error) {
		config, err := m.GetCorporationConfig(corpID)
		if err != nil {
			return nil, err
		}
		if !config.OptIn {
			return nil, ErrCorpNotRegistered
		}
		user := &User{UserID: config.OptInBy}
		return m.user.GetAuthorization(user, RoleDirector)
	})
	if err != nil {
		return nil, err
	} else if a, ok := v.(*Authorization); ok {
		return a, nil
	}
	return nil, errors.Errorf("expected *Authorization from cache, got %T", v)
}

func (m *CorpManager) SaveCorporationConfig(corpID int, detail *CorporationConfig) error {
	db, err := m.pool.Open()
	if err != nil {
		return err
	}
	defer m.pool.Release(db)
	_, err = db.Exec(
		`INSERT INTO app.corporation_settings
			 (
			     corporation_id
			   , opted_in
			   , opted_in_by
			   , opted_in_at
			   , created_by
			   , created_at
			 )
			 VALUES($1, $2, $3, $4, $5, DEFAULT)
			 ON CONFLICT
			   ON CONSTRAINT "corporation_settings_pkey"
			 DO UPDATE
			   SET opted_in = EXCLUDED.opted_in
			     , opted_in_by = EXCLUDED.opted_in_by
			     , opted_in_at = EXCLUDED.opted_in_at`,
		corpID,
		detail.OptIn,
		detail.OptInBy,
		detail.OptInDate,
		detail.CreatedBy,
	)
	return err
}

func (m *CorpManager) corpID(ctx context.Context) (int, bool) {
	if authctx, ok := ctx.(authContext); ok {
		return authctx.CorporationID(), true
	}
	return 0, false
}

func (m *CorpManager) authContext(ctx context.Context, corpID int) (context.Context, error) {
	if authctx, ok := ctx.(authContext); ok {
		if authctx.CorporationID() != corpID {
			return nil, errors.Errorf("corpID mismatch: expected %d, got %d", corpID, authctx.CorporationID())
		}
	}
	a, err := m.GetCorporationAuthorization(corpID)
	if err != nil {
		return nil, err
	}
	return a.Context(), nil
}
