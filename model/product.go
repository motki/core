package model

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/motki/motkid/evecentral"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type ProductKind string

const (
	ProductBuy         ProductKind = "buy"
	ProductManufacture             = "build"
)

type Product struct {
	ProductID     int
	parentID      int
	corporationID int

	TypeID         int
	Materials      []*Product
	Quantity       int
	MarketPrice    decimal.Decimal
	MarketRegionID int
	Kind           ProductKind
}

func (p Product) Cost() decimal.Decimal {
	if p.Kind == ProductBuy {
		return p.MarketPrice
	}
	cost := decimal.NewFromFloat(0)
	for _, m := range p.Materials {
		cost = cost.Add(m.Cost().Mul(decimal.NewFromFloat(float64(m.Quantity))))
	}
	return cost
}

func (m *Manager) NewProduct(corpID int, typeID int) (*Product, error) {
	bp, err := m.evedb.GetBlueprint(typeID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create production line for typeID %d", typeID)
	}
	p := &Product{
		corporationID: corpID,

		TypeID:         typeID,
		Materials:      make([]*Product, 0),
		Quantity:       1,
		MarketPrice:    decimal.NewFromFloat(0),
		MarketRegionID: 0,
		Kind:           ProductManufacture,
	}
	for _, mat := range bp.Materials {
		part, err := m.NewProduct(corpID, mat.ID)
		part.Kind = ProductBuy
		part.Quantity = mat.Quantity
		if err != nil {
			return nil, errors.Wrapf(err, "unable to create production line for typeID %d", typeID)
		}
		p.Materials = append(p.Materials, part)
	}
	return p, nil
}

func (m *Manager) UpdateProductMarketPrices(product *Product, regionID int) error {
	stat, err := m.GetMarketStatRegion(regionID, product.TypeID)
	if err != nil {
		return errors.Wrapf(err, "unable to update production line market price for typeID %d", product.TypeID)
	}
	var max = decimal.NewFromFloat(1000000000000)
	var bestSell = max
	for _, s := range stat {
		if s.TypeID != product.TypeID {
			continue
		}
		if s.Kind != evecentral.StatSell {
			continue
		}
		if s.Min.LessThan(bestSell) {
			bestSell = s.Min
		}
	}
	if bestSell.Equals(max) {
		return errors.Errorf("no sell orders found for typeID %d in regionID %d", product.TypeID, regionID)
	}
	product.MarketPrice = bestSell
	product.MarketRegionID = regionID
	for _, p := range product.Materials {
		err = m.UpdateProductMarketPrices(p, regionID)
		if err != nil {
			return errors.Wrapf(err, "unable to update production line market price for typeID %d", product.TypeID)
		}
	}
	return nil
}

// saveProductWithTx attempts to insert or update the given product.
//
// This method does not commit or roll-back the transaction.
//
// If a product is inserted, its ProductID field is updated.
func (m *Manager) saveProductWithTx(tx *sql.Tx, product *Product) error {
	var prodID string
	if prodID = strconv.Itoa(product.ProductID); prodID == "0" {
		prodID = "DEFAULT"
	}
	var parentID sql.NullInt64
	if product.parentID > 0 {
		parentID.Int64 = int64(product.parentID)
		parentID.Valid = true
	}
	r := tx.QueryRow(`INSERT INTO app.production_lines (
		product_id,
		type_id,
		market_price,
		market_region_id,
		quantity,
		kind,
		parent_id,
		corporation_id)
	VALUES(`+prodID+`, $1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT ON CONSTRAINT "production_lines_pkey"
		 DO UPDATE SET market_price = EXCLUDED.market_price,
		     market_region_id = EXCLUDED.market_region_id,
		     kind = EXCLUDED.kind
	RETURNING product_id`,
		product.TypeID,
		product.MarketPrice,
		product.MarketRegionID,
		product.Quantity,
		product.Kind,
		parentID,
		product.corporationID)
	id := 0
	err := r.Scan(&id)
	if err != nil {
		return err
	}
	if id == 0 {
		return errors.New("invalid last insert id")
	}
	product.ProductID = id
	for _, p := range product.Materials {
		p.parentID = id
		if err = m.saveProductWithTx(tx, p); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) SaveProduct(product *Product) error {
	c, err := m.pool.Open()
	if err != nil {
		return err
	}
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

func (m *Manager) GetAllProducts(corpID int) ([]*Product, error) {
	return m.getProducts(corpID)
}

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

func (m *Manager) getProducts(corpID int, productIDs ...int) ([]*Product, error) {
	c, err := m.pool.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	var ids = []string{}
	for _, id := range productIDs {
		ids = append(ids, strconv.Itoa(id))
	}
	var idClause = ""
	if len(ids) > 0 {
		idClause = "AND p1.product_id IN(" + strings.Join(ids, ",") + ")"
	}
	r, err := c.Query(`WITH RECURSIVE chain(f, t) AS (
		SELECT NULL::INT, p1.product_id FROM app.production_lines p1 WHERE p1.corporation_id = $1 AND p1.parent_id IS NULL `+idClause+`
		UNION
		SELECT c.t, p2.product_id
			FROM chain c
			LEFT OUTER JOIN app.production_lines p2 ON p2.parent_id = c.t
			WHERE c.t IS NOT NULL
		)
		SELECT p.product_id,
			 p.type_id,
			 p.market_price,
			 p.market_region_id,
			 p.quantity,
			 p.kind,
			 p.parent_id
		FROM app.production_lines p
		 	JOIN chain c ON c.t = p.product_id
		WHERE c.t IS NOT NULL ORDER BY p.parent_id NULLS FIRST`, corpID)
	if err != nil {
		return nil, err
	}
	prods := make(map[int]*Product)
	roots := make([]int, 0)
	for r.Next() {
		p := &Product{corporationID: corpID}
		var parentID sql.NullInt64
		err := r.Scan(&p.ProductID, &p.TypeID, &p.MarketPrice, &p.MarketRegionID, &p.Quantity, &p.Kind, &parentID)
		if err != nil {
			return nil, err
		}
		p.parentID = int(parentID.Int64)
		prods[p.ProductID] = p
		if p.parentID == 0 {
			roots = append(roots, p.ProductID)
		} else {
			parent, ok := prods[p.parentID]
			if !ok {
				return nil, errors.Errorf("unable to find product with ID %d in product map", p.parentID)
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
