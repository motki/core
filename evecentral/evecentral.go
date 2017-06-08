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

type OrderKind string

const (
	OrderBuy  OrderKind = "buy"
	OrderSell           = "sell"
	OrderAny            = "all"
)

type EveCentral struct {
	client *http.Client
}

type Stat struct {
	Kind           OrderKind
	TypeID         int
	Volume         int
	WAvg           float64
	Avg            float64
	Variance       float64
	StdDev         float64
	Median         float64
	FivePercentile float64
	Max            float64
	Min            float64
	Timestamp      time.Time
}

func New() *EveCentral {
	return &EveCentral{client: &http.Client{}}
}

func (api *EveCentral) GetMarketStat(typeIDs ...int) ([]*Stat, error) {
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
	d := make([]map[string]map[string]interface{}, 0)
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	ret := []*Stat{}
	for _, report := range d {
		for kind, info := range report {
			stat := &Stat{
				Kind: OrderKind(kind),
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
						stat.FivePercentile = f
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

func (api *EveCentral) GetMarketStatRegion(typeID int, regionID int) ([]*Stat, error) {
	return nil, nil
}

func (api *EveCentral) GetMarketStatSystem(typeID int, systemID int) ([]*Stat, error) {
	return nil, nil
}
