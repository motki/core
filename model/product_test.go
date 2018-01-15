package model_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/motki/core/model"
)

// testProduct is designed to show the characteristics of the Cost algorithm.
var testProduct = &model.Product{
	// Fake TypeIDs to allow differentiation.
	TypeID: 123,
	Materials: []*model.Product{
		{
			TypeID: 124,

			// 1000 means that 45046 (45.045... * 1000, rounded up) and
			// 162163 (162.162... * 1000, rounded up) of each material
			// are required, instead of the 46 and 163 quantity
			// required with a BatchSize of 1.
			// This means the per-item material requirement works out
			// to about 45.046 and 162.163, respectively.
			// This is demonstrated in the test below.
			BatchSize:          1000,
			MaterialEfficiency: decimal.NewFromFloat(0.11),

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
					// Cost without batch size = 46 * 100.00 = 4600 * 1000   = 4600000
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
					// Cost without batch size = 163 * 50.0 = 8150 * 1000   = 8150000
				},
			},
			Quantity:       8, // 8 with ME of 0.06 = 7.547169
			MarketPrice:    decimal.Zero,
			MarketRegionID: 0,
			Kind:           model.ProductBuild,

			// Total cost per 1000 units is 12612750; 12612.75 per unit
			// Total cost without batch size = 12750 per unit
		}, {
			TypeID:             17,
			BatchSize:          1, // Doesn't apply because the product is bought.
			MaterialEfficiency: decimal.Zero,
			Materials:          []*model.Product{},
			Quantity:           2000, // 2000 with ME of 0.06 = 1886.79245283
			MarketPrice:        decimal.NewFromFloat(5.00),
			MarketRegionID:     0,
			Kind:               model.ProductBuy,

			// Total cost per unit is 5.
		},
	},
	BatchSize:          10,
	Quantity:           1,
	MaterialEfficiency: decimal.NewFromFloat(0.06),
	MarketPrice:        decimal.Zero,
	MarketRegionID:     0,
	Kind:               model.ProductBuild,

	// Total cost per unit (including batch size)
	// ((10 * 7.547169 * 12612.75) + (10 * 1886.79 * 5.00)) / 10 =
	//            ((76 * 12612.75) + (18868 * 5.00))        / 10 =
	//                     (958569 + 94340)                 / 10 =
	//                       1052909                        / 10 = 105290.9

	// Total cost without batch size = 102000 + 9435 = 111435
}

func TestProductCost(t *testing.T) {
	prod := testProduct.Clone()
	if len(prod.Materials) != 2 {
		t.Errorf("expected 2 Materials, got %d", len(prod.Materials))
		return
	}
	id124 := prod.Materials[0]
	if cost := id124.Cost(); !cost.Equals(decimal.NewFromFloat(12612.75)) {
		t.Errorf("expected cost for typeID 124 to be 12612.61, got %s", cost)
	}
	id17 := prod.Materials[1]
	if cost := id17.Cost(); !cost.Equals(decimal.NewFromFloat(5)) {
		t.Errorf("expected cost for typeID 17 to be 5.00, got %s", cost)
	}
	if cost := prod.Cost(); !cost.Equals(decimal.NewFromFloat(105290.9)) {
		t.Errorf("expected cost for product to be 105290.9, got %s", cost)
	}
}

func TestProductCostZeroME(t *testing.T) {
	prod := testProduct.Clone()
	if len(prod.Materials) != 2 {
		t.Errorf("expected 2 Materials, got %d", len(prod.Materials))
		return
	}
	id124 := prod.Materials[0]
	if id124.TypeID != 124 {
		t.Errorf("expected first material to be type ID 124, got %d", id124.TypeID)
		return
	}
	id124.MaterialEfficiency = decimal.Zero
	prod.MaterialEfficiency = decimal.Zero
	if cost := prod.Cost(); !cost.Equals(decimal.NewFromFloat(122000)) {
		t.Errorf("expected cost for product to be 122000, got %s", cost)
	}
	// With 0 ME, batch sizing shouldn't matter either. Test this as a double check.
	id124.BatchSize = 1
	prod.BatchSize = 1
	if cost := prod.Cost(); !cost.Equals(decimal.NewFromFloat(122000)) {
		t.Errorf("expected cost for product to be 122000, got %s", cost)
	}
}

func TestProductCostBatchSizeOne(t *testing.T) {
	prod := testProduct.Clone()
	if len(prod.Materials) != 2 {
		t.Errorf("expected 2 Materials, got %d", len(prod.Materials))
		return
	}
	id124 := prod.Materials[0]
	if id124.TypeID != 124 {
		t.Errorf("expected first material to be type ID 124, got %d", id124.TypeID)
		return
	}
	id124.BatchSize = 1
	prod.BatchSize = 1
	if cost := prod.Cost(); !cost.Equals(decimal.NewFromFloat(111435)) {
		t.Errorf("expected cost for product to be 111435, got %s", cost)
	}
}
