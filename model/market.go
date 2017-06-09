package model

import (
	"time"

	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/tyler-sommer/motki/evecentral"
)

// MarketStat is reported price information for the given type.
type MarketStat struct {
	Kind        evecentral.StatKind
	TypeID      int
	Volume      int
	WAvg        decimal.Decimal
	Avg         decimal.Decimal
	Variance    decimal.Decimal
	StdDev      decimal.Decimal
	Median      decimal.Decimal
	FivePercent decimal.Decimal
	Max         decimal.Decimal
	Min         decimal.Decimal
	Timestamp   time.Time
}

func (m MarketStat) View() marketStatView {
	wavg, _ := m.WAvg.Float64()
	avg, _ := m.Avg.Float64()
	variance, _ := m.Variance.Float64()
	stddv, _ := m.StdDev.Float64()
	median, _ := m.Median.Float64()
	fivePercent, _ := m.FivePercent.Float64()
	max, _ := m.Max.Float64()
	min, _ := m.Min.Float64()
	return marketStatView{
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

type marketStatView struct {
	Kind        string
	TypeID      int
	Volume      int
	WAvg        float64
	Avg         float64
	Variance    float64
	StdDev      float64
	Median      float64
	FivePercent float64
	Max         float64
	Min         float64
	Timestamp   time.Time
}

// GetMarketStat gets market information for the given types.
func (m *Manager) GetMarketStat(typeIDs ...int) ([]*MarketStat, error) {
	res, err := m.getMarketStatFromDB(0, 0, typeIDs...)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return m.getMarketStatFromAPI(0, 0, typeIDs...)
	}
	return res, nil
}

// GetMarketStatRegion gets market information for the given region and types.
func (m *Manager) GetMarketStatRegion(regionID int, typeIDs ...int) ([]*MarketStat, error) {
	res, err := m.getMarketStatFromDB(regionID, 0, typeIDs...)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return m.getMarketStatFromAPI(regionID, 0, typeIDs...)
	}
	return res, nil
}

// GetMarketStatSystem gets market information for the given system and types.
func (m *Manager) GetMarketStatSystem(systemID int, typeIDs ...int) ([]*MarketStat, error) {
	res, err := m.getMarketStatFromDB(0, systemID, typeIDs...)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return m.getMarketStatFromAPI(0, systemID, typeIDs...)
	}
	return res, nil
}

func (m *Manager) getMarketStatFromDB(regionID, systemID int, typeIDs ...int) ([]*MarketStat, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	ids := []string{}
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
	res := []*MarketStat{}
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
		r.Kind = evecentral.StatKind(kind)
		res = append(res, r)
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (m *Manager) getMarketStatFromAPI(regionID, systemID int, typeIDs ...int) ([]*MarketStat, error) {
	var stats []*evecentral.MarketStat
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

func (m *Manager) apiMarketStatToDB(regionID, systemID int, stats []*evecentral.MarketStat) ([]*MarketStat, error) {
	db, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	res := []*MarketStat{}
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
