// Package evecentral contains a client integration with the eve-central.com API.
package evecentral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const baseURL = "https://api.eve-central.com/api/marketstat/json"

// StatKind describes the type of market prices in a MarketStat.
type StatKind string

const (
	StatBuy  StatKind = "buy"
	StatSell          = "sell"
	StatAll           = "all"
)

// EveCentral is a client for retrieving market data from the eve-central.com API.
type EveCentral struct {
	client *http.Client
}

// MarketStat is reported price information for the given type.
type MarketStat struct {
	Kind        StatKind
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

// New creates a new EveCentral API client.
func New() *EveCentral {
	return &EveCentral{client: &http.Client{}}
}

// GetMarketStat gets market information for the given types.
func (api *EveCentral) GetMarketStat(typeIDs ...int) ([]*MarketStat, error) {
	params := make([]string, 0)
	for _, id := range typeIDs {
		params = append(params, fmt.Sprintf("typeid=%d", id))
	}
	res, err := api.client.Get(fmt.Sprintf("%s?%s", baseURL, strings.Join(params, "&")))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return parseBody(body)
}

// GetMarketStatRegion gets market information for the given region and types.
func (api *EveCentral) GetMarketStatRegion(regionID int, typeIDs ...int) ([]*MarketStat, error) {
	params := make([]string, 0)
	for _, id := range typeIDs {
		params = append(params, fmt.Sprintf("typeid=%d", id))
	}
	res, err := api.client.Get(fmt.Sprintf("%s?%s&regionlimit=%d", baseURL, strings.Join(params, "&"), regionID))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return parseBody(body)
}

// GetMarketStatSystem gets market information for the given system and types.
func (api *EveCentral) GetMarketStatSystem(systemID int, typeIDs ...int) ([]*MarketStat, error) {
	params := make([]string, 0)
	for _, id := range typeIDs {
		params = append(params, fmt.Sprintf("typeid=%d", id))
	}
	res, err := api.client.Get(fmt.Sprintf("%s?%s&usesystem=%d", baseURL, strings.Join(params, "&"), systemID))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return parseBody(body)
}

func parseBody(body []byte) ([]*MarketStat, error) {
	d := make([]map[string]map[string]interface{}, 0)
	err := json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	ret := []*MarketStat{}
	for _, report := range d {
		for kind, info := range report {
			stat := &MarketStat{
				Kind: StatKind(kind),
			}
			for k, v := range info {
				switch k {
				case "forQuery":
					if t, ok := v.(map[string]interface{}); ok {
						if ids, ok := t["types"].([]interface{}); ok {
							stat.TypeID = int(ids[0].(float64))
						}
					}
				case "volume":
					if f, ok := v.(float64); ok {
						stat.Volume = int(f)
					}
				case "wavg":
					if f, ok := v.(float64); ok {
						stat.WAvg = f
					}
				case "avg":
					if f, ok := v.(float64); ok {
						stat.Avg = f
					}
				case "variance":
					if f, ok := v.(float64); ok {
						stat.Variance = f
					}
				case "stdDev":
					if f, ok := v.(float64); ok {
						stat.StdDev = f
					}
				case "median":
					if f, ok := v.(float64); ok {
						stat.Median = f
					}
				case "fivePercent":
					if f, ok := v.(float64); ok {
						stat.FivePercent = f
					}
				case "max":
					if f, ok := v.(float64); ok {
						stat.Max = f
					}
				case "min":
					if f, ok := v.(float64); ok {
						stat.Min = f
					}
				case "generated":
					if f, ok := v.(float64); ok {
						stat.Timestamp = time.Unix(int64(f)/1000, int64(f)%1000)
					}
				}
			}
			ret = append(ret, stat)
		}
	}
	return ret, nil
}
