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

DROP TABLE IF EXISTS app.market_orders;
CREATE TABLE app.market_orders
(
  order_id BIGINT PRIMARY KEY NOT NULL,
  corporation_id BIGINT NOT NULL,
  character_id BIGINT NOT NULL,
  station_id BIGINT NOT NULL,
  type_id INT NOT NULL,
  volume_entered BIGINT NOT NULL,
  volume_remaining BIGINT NOT NULL,
  min_volume BIGINT NOT NULL,
  order_state INT NOT NULL,
  range INT NOT NULL,
  account_key INT NOT NULL,
  duration INT NOT NULL,
  escrow NUMERIC NOT NULL,
  price NUMERIC NOT NULL,
  bid SMALLINT NOT NULL DEFAULT 0,
  issued TIMESTAMP NOT NULL,
  fetched_at TIMESTAMP NOT NULL DEFAULT NOW(),
  loner SMALLINT NOT NULL DEFAULT 0
)