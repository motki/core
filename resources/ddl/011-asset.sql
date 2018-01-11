DROP TABLE IF EXISTS app.assets;
CREATE TABLE app.assets
(
  corporation_id BIGINT NOT NULL,
  character_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  location_id BIGINT NOT NULL,
  type_id BIGINT NOT NULL,
  quantity BIGINT NOT NULL,
  raw_quantity BIGINT NOT NULL,
  singleton BOOLEAN NOT NULL,
  flag_id BIGINT NOT NULL,
  fetched_at TIMESTAMP NOT NULL DEFAULT NOW(),
  valid BOOLEAN NOT NULL DEFAULT TRUE
);
