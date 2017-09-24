DROP TABLE IF EXISTS app.market_data;
CREATE TABLE app.market_data
(
  id SERIAL PRIMARY KEY NOT NULL,
  kind VARCHAR(4) NOT NULL,
  type_id INT NOT NULL,
  region_id INT NOT NULL DEFAULT 0,
  system_id INT NOT NULL DEFAULT 0,
  volume BIGINT NOT NULL,
  wavg NUMERIC NOT NULL,
  avg NUMERIC NOT NULL,
  variance NUMERIC NOT NULL,
  stddev NUMERIC NOT NULL,
  median NUMERIC NOT NULL,
  five_percent NUMERIC NOT NULL,
  max NUMERIC NOT NULL,
  min NUMERIC NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  fetched_at TIMESTAMP NOT NULL DEFAULT NOW()
);

DROP TABLE IF EXISTS app.market_prices;
CREATE TABLE app.market_prices
(
  id SERIAL PRIMARY KEY NOT NULL,
  type_id INT NOT NULL,
  avg NUMERIC NOT NULL,
  base NUMERIC NOT NULL,
  fetched_at TIMESTAMP NOT NULL DEFAULT NOW()
);

DROP INDEX IF EXISTS idx_market_prices_type_id;
CREATE INDEX idx_market_prices_type_id
  ON app.market_prices (type_id);