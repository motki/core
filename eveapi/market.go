package eveapi

import (
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type MarketOrderKind int

const (
	MarketOrderKindAll MarketOrderKind = iota
	MarketOrderKindSell
	MarketOrderKindBuy
)

type MarketOrder struct {
	Kind         MarketOrderKind
	LocationID   int
	Range        string
	Price        decimal.Decimal
	MinVolume    int // Minimum quantity for buy orders.
	VolumeRemain int
	VolumeTotal  int
	Duration     time.Duration
	DateIssued   time.Time
}

type MarketPrice struct {
	TypeID       int
	AveragePrice decimal.Decimal
	BasePrice    decimal.Decimal
}

func (api *EveAPI) GetMarketPrices() (prices chan MarketPrice, cancelFn func(), err error) {
	res, _, err := api.client.ESI.MarketApi.GetMarketsPrices(nil)
	if err != nil {
		return nil, nil, err
	}
	if len(res) == 0 {
		return nil, nil, errors.New("no results.")
	}
	prices = make(chan MarketPrice, 16)
	abort := make(chan struct{})
	cancelFn = func() {
		close(abort)
	}
	go func() {
		defer func() {
			close(prices)
		}()
		i := 0
		// We know there is at least one result from above
		p := MarketPrice{
			TypeID:       int(res[i].TypeId),
			AveragePrice: decimal.NewFromFloat(float64(res[i].AveragePrice)),
			BasePrice:    decimal.NewFromFloat(float64(res[i].AdjustedPrice)),
		}
		for {
			select {
			case <-abort:
				return
			case prices <- p:
				i += 1
			}
			if i >= len(res) {
				return
			} else {
				p = MarketPrice{
					TypeID:       int(res[i].TypeId),
					AveragePrice: decimal.NewFromFloat(float64(res[i].AveragePrice)),
					BasePrice:    decimal.NewFromFloat(float64(res[i].AdjustedPrice)),
				}
			}

		}
	}()
	return prices, cancelFn, nil
}

func (api *EveAPI) GetMarketOrdersRegionTypeID(regionID, typeID int) ([]*MarketStat, error) {
	// TODO: implement
	panic("not implemented")
	return nil, nil
}

// MarketStat is reported price information for the given type.
type MarketStat struct {
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

type marketStatSlice []*MarketStat

func (s marketStatSlice) Less(i, j int) bool {
	return s[i].Timestamp.Before(s[j].Timestamp)
}

// Len is the number of elements in the collection.
func (s marketStatSlice) Len() int {
	return len(s)
}

// Swap swaps the elements with indexes i and j.
func (s marketStatSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (api *EveAPI) GetMarketHistoryRegionTypeID(regionID, typeID int) ([]*MarketStat, error) {
	his, _, err := api.client.ESI.MarketApi.GetMarketsRegionIdHistory(int32(regionID), int32(typeID), nil)
	if err != nil {
		return nil, err
	}
	var res []*MarketStat
	for _, h := range his {
		t, _ := time.Parse("2006-01-02", h.Date)
		s := &MarketStat{
			TypeID:    typeID,
			Volume:    int(h.Volume),
			Avg:       decimal.NewFromFloat(float64(h.Average)),
			Max:       decimal.NewFromFloat(float64(h.Highest)),
			Min:       decimal.NewFromFloat(float64(h.Lowest)),
			Timestamp: t,
		}
		res = append(res, s)
	}
	sort.Reverse(marketStatSlice(res))
	return res, nil
}
