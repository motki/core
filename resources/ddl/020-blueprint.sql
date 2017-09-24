DROP TABLE IF EXISTS app.blueprints;
CREATE TABLE app.blueprints
(
  corporation_id BIGINT NOT NULL,
  character_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  location_id BIGINT NOT NULL,
  type_id BIGINT NOT NULL,
  type_name VARCHAR(500) NOT NULL,
  quantity BIGINT NOT NULL,
  kind VARCHAR NOT NULL CONSTRAINT blueprints_valid_kinds CHECK (kind = 'bpc' OR kind = 'bpo'),
  flag_id BIGINT NOT NULL,
  time_efficiency BIGINT NOT NULL,
  material_efficiency BIGINT NOT NULL,
  runs BIGINT NOT NULL,
  fetched_at TIMESTAMP NOT NULL DEFAULT NOW()
);
