package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/motki/core/evemarketer"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// MarketStat is region-specific price information for the given type.
type MarketStat struct {
	Kind        evemarketer.StatKind `json:"kind"`
	TypeID      int                  `json:"type_id"`
	Volume      int                  `json:"volume"`
	WAvg        decimal.Decimal      `json:"w_avg"`
	Avg         decimal.Decimal      `json:"avg"`
	Variance    decimal.Decimal      `json:"variance"`
	StdDev      decimal.Decimal      `json:"std_dev"`
	Median      decimal.Decimal      `json:"median"`
	FivePercent decimal.Decimal      `json:"five_percent"`
	Max         decimal.Decimal      `json:"max"`
	Min         decimal.Decimal      `json:"min"`
	Timestamp   time.Time            `json:"timestamp"`
}

func (m MarketStat) View() StatView {
	wavg, _ := m.WAvg.Float64()
	avg, _ := m.Avg.Float64()
	variance, _ := m.Variance.Float64()
	stddv, _ := m.StdDev.Float64()
	median, _ := m.Median.Float64()
	fivePercent, _ := m.FivePercent.Float64()
	max, _ := m.Max.Float64()
	min, _ := m.Min.Float64()
	return StatView{
		Kind:        string(m.Kind),
		TypeID:      m.TypeID,
		Volume:      m.Volume,
		WAvg:        wavg,
		Avg:         avg,
		Variance:    variance,
		StdDev:      stddv,
		Median:      median,
		FivePercent: fivePercent,
		Max:         max,
		Min:         min,
		Timestamp:   m.Timestamp,
	}
}

type StatView struct {
	Kind        string    `json:"kind"`
	TypeID      int       `json:"type_id"`
	Volume      int       `json:"volume"`
	WAvg        float64   `json:"w_avg"`
	Avg         float64   `json:"avg"`
	Variance    float64   `json:"variance"`
	StdDev      float64   `json:"std_dev"`
	Median      float64   `json:"median"`
	FivePercent float64   `json:"five_percent"`
	Max         float64   `json:"max"`
	Min         float64   `json:"min"`
	Timestamp   time.Time `json:"timestamp"`

	unexported struct{}
}

type MarketManager struct {
	bootstrap

	corp *CorpManager
}

func newMarketManager(m bootstrap, corp *CorpManager) *MarketManager {
	return &MarketManager{m, corp}
}

// GetMarketStat gets market information for the given types.
//
// Multiple typeIDs may be specified, but the method signature requires at least
// the first is given.
func (m *MarketManager) GetMarketStat(typeID int, typeIDs ...int) ([]*MarketStat, error) {
	return m.getMarketStatFromDB(0, 0, append(typeIDs, typeID)...)
}

// GetMarketStatRegion gets market information for the given region and types.
//
// Multiple typeIDs may be specified, but the method signature requires at least
// the first is given.
func (m *MarketManager) GetMarketStatRegion(regionID int, typeID int, typeIDs ...int) ([]*MarketStat, error) {
	return m.getMarketStatFromDB(regionID, 0, append(typeIDs, typeID)...)
}

// GetMarketStatSystem gets market information for the given system and types.
//
// Multiple typeIDs may be specified, but the method signature requires at least
// the first is given.
func (m *MarketManager) GetMarketStatSystem(systemID int, typeID int, typeIDs ...int) ([]*MarketStat, error) {
	return m.getMarketStatFromDB(0, systemID, append(typeIDs, typeID)...)
}

func (m *MarketManager) getMarketStatFromDB(regionID, systemID int, typeIDs ...int) ([]*MarketStat, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	var ids []string
	for _, id := range typeIDs {
		ids = append(ids, fmt.Sprintf("%d", id))
	}
	rs, err := c.Query(
		`SELECT
			  c.kind
			, c.type_id
			, c.volume
			, c.wavg
			, c.avg
			, c.variance
			, c.stddev
			, c.five_percent
			, c.max
			, c.min
			, c.timestamp
			FROM app.market_data c
			WHERE c.type_id = ANY($1::INTEGER[])
			  AND c.region_id = $2
			  AND c.system_id = $3
			  AND c.fetched_at >= (NOW() - INTERVAL '1 day')`, "{"+strings.Join(ids, ",")+"}", regionID, systemID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*MarketStat
	for rs.Next() {
		r := &MarketStat{}
		kind := ""
		err := rs.Scan(
			&kind,
			&r.TypeID,
			&r.Volume,
			&r.WAvg,
			&r.Avg,
			&r.Variance,
			&r.StdDev,
			&r.FivePercent,
			&r.Max,
			&r.Min,
			&r.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		r.Kind = evemarketer.StatKind(kind)
		res = append(res, r)
	}
	if len(res) == 0 {
		// No results, get them from the API
		return m.getMarketStatFromAPI(regionID, systemID, typeIDs...)
	}
	got := map[int]struct{}{}
	for _, s := range res {
		got[s.TypeID] = struct{}{}
	}
	// If we're missing any stats for some of the types, fetch them now.
	if len(got) != len(typeIDs) {
		ids := make([]int, 0)
		for _, id := range typeIDs {
			if _, ok := got[id]; !ok {
				ids = append(ids, id)
			}
		}
		ares, err := m.getMarketStatFromAPI(regionID, systemID, ids...)
		if err != nil {
			return nil, err
		}
		res = append(res, ares...)
	}
	return res, nil
}

func (m *MarketManager) getMarketStatFromAPI(regionID, systemID int, typeIDs ...int) ([]*MarketStat, error) {
	var stats []*evemarketer.MarketStat
	var err error
	switch {
	case regionID != 0:
		stats, err = m.ec.GetMarketStatRegion(regionID, typeIDs...)
	case systemID != 0:
		stats, err = m.ec.GetMarketStatSystem(systemID, typeIDs...)
	default:
		stats, err = m.ec.GetMarketStat(typeIDs...)
	}
	if err != nil {
		return nil, err
	}
	return m.apiMarketStatToDB(regionID, systemID, stats)
}

func (m *MarketManager) apiMarketStatToDB(regionID, systemID int, stats []*evemarketer.MarketStat) ([]*MarketStat, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(db)
	var res []*MarketStat
	for _, stat := range stats {
		s := &MarketStat{
			TypeID:      stat.TypeID,
			Kind:        stat.Kind,
			Volume:      stat.Volume,
			WAvg:        stat.WAvg,
			Avg:         stat.Avg,
			Variance:    stat.Variance,
			StdDev:      stat.StdDev,
			Median:      stat.Median,
			FivePercent: stat.FivePercent,
			Max:         stat.Max,
			Min:         stat.Min,
			Timestamp:   stat.Timestamp,
		}
		_, err = db.Exec(
			"INSERT INTO app.market_data (id, type_id, kind, region_id, system_id, volume, wavg, avg, variance, stddev, median, five_percent, max, min, timestamp, fetched_at) VALUES(DEFAULT, $1, $2, $13, $14, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, DEFAULT)",
			s.TypeID,
			string(s.Kind),
			s.Volume,
			s.WAvg,
			s.Avg,
			s.Variance,
			s.StdDev,
			s.Median,
			s.FivePercent,
			s.Max,
			s.Min,
			s.Timestamp,
			regionID,
			systemID,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	return res, nil
}

// A MarketPrice is a universal average price for a given item.
type MarketPrice struct {
	TypeID int             `json:"type_id"`
	Avg    decimal.Decimal `json:"avg"`
	Base   decimal.Decimal `json:"base"`
}

func (m *MarketManager) GetMarketPrice(typeID int) (*MarketPrice, error) {
	res, err := m.getMarketPricesFromDB(typeID)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.Errorf("expected 1 result, got %d", len(res))
	}
	return res[0], nil
}

func (m *MarketManager) GetMarketPrices(typeID int, typeIDs ...int) ([]*MarketPrice, error) {
	return m.getMarketPricesFromDB(append(typeIDs, typeID)...)
}

func (m *MarketManager) getMarketPricesFromDB(typeIDs ...int) ([]*MarketPrice, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	var ids []string
	for _, id := range typeIDs {
		ids = append(ids, fmt.Sprintf("%d", id))
	}
	rs, err := c.Query(
		`SELECT
			  c.type_id
			, c.avg
			, c.base
			FROM app.market_prices c
			WHERE c.type_id = ANY($1::INTEGER[])
			  AND c.fetched_at >= (NOW() - INTERVAL '1 day')`, "{"+strings.Join(ids, ",")+"}")
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var res []*MarketPrice
	for rs.Next() {
		r := &MarketPrice{}
		err := rs.Scan(
			&r.TypeID,
			&r.Avg,
			&r.Base,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	if len(res) == 0 {
		// No results, get them from the API
		return m.getMarketPricesFromAPI(typeIDs...)
	}
	got := map[int]struct{}{}
	for _, s := range res {
		got[s.TypeID] = struct{}{}
	}
	// If we're missing any stats for some of the types, fetch them now.
	if len(got) != len(typeIDs) {
		ids := make([]int, 0)
		for _, id := range typeIDs {
			if _, ok := got[id]; !ok {
				ids = append(ids, id)
			}
		}
		ares, err := m.getMarketPricesFromAPI(ids...)
		if err != nil {
			return nil, err
		}
		res = append(res, ares...)
	}
	return res, nil
}

func (m *MarketManager) getMarketPricesFromAPI(typeIDs ...int) ([]*MarketPrice, error) {
	p, cancel, err := m.eveapi.GetMarketPrices()
	if err != nil {
		return nil, err
	}
	defer cancel()
	var prices []*MarketPrice
	var results []*MarketPrice
	mp := make(map[int]struct{})
	for _, t := range typeIDs {
		mp[t] = struct{}{}
	}
	for pr := range p {
		mktp := &MarketPrice{
			TypeID: pr.TypeID,
			Avg:    pr.AveragePrice,
			Base:   pr.BasePrice,
		}
		if _, ok := mp[mktp.TypeID]; ok {
			results = append(results, mktp)
		}
		prices = append(prices, mktp)
	}
	// Save all the prices, but only return what we were asked for.
	return results, m.apiMarketPricesToDB(prices)
}

func (m *MarketManager) apiMarketPricesToDB(prices []*MarketPrice) error {
	db, err := m.pool.Open()
	if err != nil {
		return err
	}
	defer m.pool.Release(db)
	var res []*MarketPrice
	for _, stat := range prices {
		s := &MarketPrice{
			TypeID: stat.TypeID,
			Avg:    stat.Avg,
			Base:   stat.Base,
		}
		_, err = db.Exec(
			`INSERT INTO app.market_prices (id, type_id, avg, base, fetched_at)
					VALUES(DEFAULT, $1, $2, $3, DEFAULT)`,
			s.TypeID,
			s.Avg,
			s.Base,
		)
		if err != nil {
			return err
		}
		res = append(res, s)
	}
	return nil
}
