DROP SEQUENCE IF EXISTS app.production_chains_id_seq CASCADE;
CREATE SEQUENCE app.production_chains_id_seq;

DROP TABLE IF EXISTS app.production_chains;
CREATE TABLE app.production_chains
(
  product_id INT PRIMARY KEY NOT NULL DEFAULT NEXTVAL('app.production_chains_id_seq'),
  parent_id INT NULL,
  type_id BIGINT NOT NULL,
  kind VARCHAR(5) NOT NULL CONSTRAINT production_chains_valid_kinds CHECK (kind = 'buy' OR kind = 'build'),
  quantity INT NOT NULL,
  market_price NUMERIC NULL,
  market_region_id BIGINT NULL,
  material_efficiency NUMERIC NOT NULL,
  batch_size INT NOT NULL,
  corporation_id BIGINT NOT NULL
);
