package model_test

import (
	"testing"

	"github.com/motki/motki/model"
	"github.com/shopspring/decimal"
)

// testProduct is designed to show the characteristics of the Cost algorithm.
var testProduct = &model.Product{
	ProductID: 0,
	TypeID:    123,
	Materials: []*model.Product{
		{
			TypeID: 124,

			// 1000 means that 45046 (45.045... * 1000, rounded up) and
			// 162163 (162.162... * 1000, rounded up) are required.
			// This means the individual cost in material per item is
			// 45.046 and 162.163, instead of the 46 and 163 quantity
			// required with a BatchSize of 1. This is demonstrated in the
			// test below.
			BatchSize: 1000,

			Materials: []*model.Product{
				{
					TypeID:             15,
					Materials:          []*model.Product{},
					Quantity:           50, // 50 with ME of 0.11 = 45.045045...
					MarketPrice:        decimal.NewFromFloat(100.00),
					MarketRegionID:     0,
					MaterialEfficiency: decimal.Zero,
					BatchSize:          1,
					Kind:               model.ProductBuy,

					// Cost per 1000 units = ceil(1000 * 50 / 1.11) * 100.00 = 4504600
				}, {
					TypeID:             16,
					Materials:          []*model.Product{},
					Quantity:           180, // 180 with ME of 0.11 = 162.162162...
					MarketPrice:        decimal.NewFromFloat(50.00),
					MarketRegionID:     0,
					MaterialEfficiency: decimal.Zero,
					BatchSize:          1,
					Kind:               model.ProductBuy,

					// Cost per 1000 units = ceil(1000 * 180 / 1.11) * 50.0 = 8108150
				},
			},
			Quantity:           8, // 8 with ME of 0.06 = 7.547169
			MarketPrice:        decimal.Zero,
			MarketRegionID:     0,
			MaterialEfficiency: decimal.NewFromFloat(0.11),
			Kind:               model.ProductBuild,

			// Total cost per 1000 units is 12612750; per unit is 12612.75
		}, {
			TypeID:             17,
			Materials:          []*model.Product{},
			Quantity:           2000, // 2000 with ME of 0.06 = 1886.79245283
			MarketPrice:        decimal.NewFromFloat(5.00),
			MarketRegionID:     0,
			MaterialEfficiency: decimal.Zero,
			BatchSize:          1,
			Kind:               model.ProductBuy,

			// Total cost per unit is 5.
		},
	},
	// And another layer of batch sizing
	BatchSize:          10,
	Quantity:           1,
	MarketPrice:        decimal.Zero,
	MarketRegionID:     0,
	MaterialEfficiency: decimal.NewFromFloat(0.06),
	Kind:               model.ProductBuild,

	// Total cost per unit is
	// ((10 * 7.547169 * 12612.61) + (10 * 1886.79 * 5.00)) / 10 =
	//            ((76 * 12612.75) + (18868 * 5.00))        / 10 =
	//                     (958569 + 94340)                 / 10 =
	//                       1052909                        / 10 = 105290.9
}

func TestProductCost(t *testing.T) {
	if len(testProduct.Materials) != 2 {
		t.Errorf("expected 2 Materials, got %d", len(testProduct.Materials))
		return
	}
	id124 := testProduct.Materials[0]
	if cost := id124.Cost(); !cost.Equals(decimal.NewFromFloat(12612.75)) {
		t.Errorf("expected cost for typeID 124 to be 12612.61, got %s", cost)
	}
	id17 := testProduct.Materials[1]
	if cost := id17.Cost(); !cost.Equals(decimal.NewFromFloat(5)) {
		t.Errorf("expected cost for typeID 17 to be 5.00, got %s", cost)
	}
	if cost := testProduct.Cost(); !cost.Equals(decimal.NewFromFloat(105290.9)) {
		t.Errorf("expected cost for testProduct to be 105290.9, got %s", cost)
	}
}
