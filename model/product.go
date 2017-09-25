package model

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// ProductKind describes the way a product is acquired.
type ProductKind string

const (
	ProductBuy   ProductKind = "buy"
	ProductBuild ProductKind = "build"
)

// Product represents one part of a production chain.
type Product struct {
	ProductID          int
	TypeID             int
	Materials          []*Product
	Quantity           int
	MarketPrice        decimal.Decimal
	MarketRegionID     int
	MaterialEfficiency decimal.Decimal
	BatchSize          int
	Kind               ProductKind

	ParentID      int
	CorporationID int
}

// Cost returns the total cost for one single unit of the completed parent product.
func (p Product) Cost() decimal.Decimal {
	batchSize := decimal.NewFromFloat(float64(p.BatchSize))
	if p.Kind == ProductBuy {
		return p.MarketPrice
	}
	// Calculate the cost, and be sure to include the tiny savings received on
	// ME% bonuses when calculating larger job sizes. We do this by multiplying
	// the material cost for each component by the batch size, then dividing
	// by the batch size at the end to scale the final cost back to single-
	// product scale.
	cost := decimal.NewFromFloat(0)
	for _, m := range p.Materials {
		// cost = cost + (m.Cost * round(m.Quantity / (1 + p.MaterialEfficiency) * p.BatchSize)
		cost = cost.Add(m.Cost().
			Mul(decimal.NewFromFloat(float64(m.Quantity)).
				Div(p.MaterialEfficiency.Add(decimal.NewFromFloat(1))).
				Mul(batchSize).
				Round(0)))
	}
	// Bring the final cost back to a single-product scale by dividing by the
	// total component cost by the batch size at the end.
	return cost.Div(batchSize)
}

// NewProduct creates a new production chain for the given corporation and type.
func (m *Manager) NewProduct(corpID int, typeID int) (*Product, error) {
	bp, err := m.evedb.GetBlueprint(typeID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create production chain for typeID %d", typeID)
	}
	p := &Product{
		CorporationID:      corpID,
		TypeID:             typeID,
		Materials:          make([]*Product, 0),
		Quantity:           bp.ProducesQty,
		MarketPrice:        decimal.Zero,
		MarketRegionID:     0,
		MaterialEfficiency: decimal.Zero,
		BatchSize:          1,
		Kind:               ProductBuild,
	}
	for _, mat := range bp.Materials {
		part, err := m.NewProduct(corpID, mat.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to create production chain for typeID %d", typeID)
		}
		part.Kind = ProductBuy
		part.Quantity = mat.Quantity
		p.Materials = append(p.Materials, part)
	}
	return p, nil
}

// UpdateProductMarketPrices fetches the latest market data for the production chain.
func (m *Manager) UpdateProductMarketPrices(product *Product, regionID int) error {
	// EveCentral is down! :(
	//stat, err := m.GetMarketStatRegion(regionID, product.TypeID)
	//if err != nil {
	//	return errors.Wrapf(err, "unable to update production chain market price for typeID %d", product.TypeID)
	//}
	//max := decimal.NewFromFloat(1000000000000)
	//var bestSell = max
	//for _, s := range stat {
	//	if s.TypeID != product.TypeID {
	//		continue
	//	}
	//	if s.Kind != evecentral.StatSell {
	//		continue
	//	}
	//	if s.Min.LessThan(bestSell) {
	//		bestSell = s.Min
	//	}
	//}
	//if bestSell.Equals(max) {
	//	return errors.Errorf("no sell orders found for typeID %d in regionID %d", product.TypeID, regionID)
	//}
	//product.MarketPrice = bestSell
	price, err := m.GetMarketPrice(product.TypeID)
	if err != nil {
		return errors.Wrapf(err, "unable to update production chain market price for typeID %d", product.TypeID)
	}
	product.MarketPrice = price.Avg
	product.MarketRegionID = regionID
	for _, p := range product.Materials {
		if err = m.UpdateProductMarketPrices(p, regionID); err != nil {
			return errors.Wrapf(err, "unable to update production chain market price for typeID %d", product.TypeID)
		}
	}
	return nil
}

// saveProductWithTx attempts to insert or update the given product.
//
// This method does not commit or roll-back the transaction.
//
// If a product is inserted, its ProductID field is updated.
func (m *Manager) saveProductWithTx(tx *pgx.Tx, product *Product) error {
	prodID := "DEFAULT"
	if n := product.ProductID; n > 0 {
		prodID = strconv.Itoa(n)
	}
	var parentID sql.NullInt64
	if product.ParentID > 0 {
		parentID.Int64 = int64(product.ParentID)
		parentID.Valid = true
	}
	r := tx.QueryRow(`INSERT INTO app.production_chains (
		product_id,
		type_id,
		market_price,
		market_region_id,
		quantity,
		material_efficiency,
		batch_size,
		kind,
		parent_id,
		corporation_id)
	VALUES(`+prodID+`, $1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT ON CONSTRAINT "production_chains_pkey"
		 DO UPDATE SET market_price = EXCLUDED.market_price,
		     market_region_id = EXCLUDED.market_region_id,
		     kind = EXCLUDED.kind,
		     material_efficiency = EXCLUDED.material_efficiency,
		     batch_size = EXCLUDED.batch_size
	RETURNING product_id`,
		product.TypeID,
		product.MarketPrice,
		product.MarketRegionID,
		product.Quantity,
		product.MaterialEfficiency,
		product.BatchSize,
		product.Kind,
		parentID,
		product.CorporationID)
	id := 0
	if err := r.Scan(&id); err != nil {
		return err
	}
	if id == 0 {
		return errors.New("invalid last insert id")
	}
	product.ProductID = id
	for _, p := range product.Materials {
		p.ParentID = id
		if err := m.saveProductWithTx(tx, p); err != nil {
			return err
		}
	}
	return nil
}

// SaveProduct saves the given production chain in the database.
//
// This function automatically handles both inserting and updating.
func (m *Manager) SaveProduct(product *Product) error {
	c, err := m.pool.Open()
	if err != nil {
		return err
	}
	defer m.pool.Release(c)
	tx, err := c.Begin()
	if err != nil {
		return err
	}
	err = m.saveProductWithTx(tx, product)
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			err = errors.Wrapf(err, "unable to rollback db transaction: %s", errTx.Error())
		}
		return err
	}
	return errors.Wrap(tx.Commit(), "couldn't commit db transaction")
}

// GetAllProducts returns all production chains associated with the given corporation.
func (m *Manager) GetAllProducts(corpID int) ([]*Product, error) {
	return m.getProducts(corpID)
}

// GetProduct returns a production chain for the given corporation and root product.
func (m *Manager) GetProduct(corpID int, productID int) (*Product, error) {
	prods, err := m.getProducts(corpID, productID)
	if err != nil {
		return nil, err
	}
	if len(prods) == 0 {
		return nil, errors.Errorf("no root product found with corpID %d and productID %d", corpID, productID)
	}
	return prods[0], nil
}

// getProducts handles fetching production chain components.
func (m *Manager) getProducts(corpID int, productIDs ...int) ([]*Product, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer m.pool.Release(c)
	var ids []string
	for _, id := range productIDs {
		ids = append(ids, strconv.Itoa(id))
	}
	var idClause string
	if len(ids) > 0 {
		idClause = "AND p1.product_id IN(" + strings.Join(ids, ",") + ")"
	}
	// This relies on a recursive CTE to generate the full list of products needed for the given root products.
	r, err := c.Query(`WITH RECURSIVE chain(f, t) AS (
		SELECT NULL::INT, p1.product_id FROM app.production_chains p1 WHERE p1.corporation_id = $1 AND p1.parent_id IS NULL `+idClause+`
		UNION
		SELECT c.t, p2.product_id
			FROM chain c
			LEFT OUTER JOIN app.production_chains p2 ON p2.parent_id = c.t
			WHERE c.t IS NOT NULL
		)
		SELECT p.product_id
		     , p.type_id
		     , p.market_price
		     , p.market_region_id
		     , p.quantity
		     , p.material_efficiency
		     , p.batch_size
		     , p.kind
		     , p.parent_id
		FROM app.production_chains p
		 	JOIN chain c ON c.t = p.product_id
		WHERE c.t IS NOT NULL ORDER BY p.parent_id NULLS FIRST, p.product_id`, corpID)
	if err != nil {
		return nil, err
	}
	prods := make(map[int]*Product)
	roots := make([]int, 0)
	for r.Next() {
		p := &Product{CorporationID: corpID}
		var parentID sql.NullInt64
		err := r.Scan(&p.ProductID, &p.TypeID, &p.MarketPrice, &p.MarketRegionID, &p.Quantity, &p.MaterialEfficiency, &p.BatchSize, &p.Kind, &parentID)
		if err != nil {
			return nil, err
		}
		p.ParentID = int(parentID.Int64)
		prods[p.ProductID] = p
		if p.ParentID == 0 {
			roots = append(roots, p.ProductID)
		} else {
			parent, ok := prods[p.ParentID]
			if !ok {
				return nil, errors.Errorf("unable to find product with ID %d in product map", p.ParentID)
			}
			parent.Materials = append(parent.Materials, p)
		}
	}
	prodSlice := make([]*Product, len(roots))
	for i, prodID := range roots {
		if p, ok := prods[prodID]; ok {
			prodSlice[i] = p
		} else {
			return nil, errors.Errorf("unable to find product with ID %d in product map", prodID)
		}
	}
	return prodSlice, nil
}
