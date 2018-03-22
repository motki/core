DROP TABLE IF EXISTS app.assets;
CREATE TABLE app.assets
(
  corporation_id BIGINT NOT NULL,
  character_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  location_id BIGINT NOT NULL,
  location_type VARCHAR(500) NOT NULL,
  location_flag VARCHAR(500) NOT NULL,
  type_id BIGINT NOT NULL,
  quantity BIGINT NOT NULL,
  singleton BOOLEAN NOT NULL,
  fetched_at TIMESTAMP NOT NULL DEFAULT NOW(),
  valid BOOLEAN NOT NULL DEFAULT TRUE
);
