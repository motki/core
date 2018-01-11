DROP TABLE IF EXISTS app.inventory_items;
CREATE TABLE app.inventory_items
(
  type_id BIGINT NOT NULL,
  location_id BIGINT NOT NULL,
  curr_level INT NOT NULL,
  min_level NUMERIC NULL,
  corporation_id BIGINT NOT NULL,
  fetched_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (corporation_id, type_id, location_id)
);
