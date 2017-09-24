DROP TABLE IF EXISTS app.structures;
CREATE TABLE app.structures
(
  corporation_id BIGINT NOT NULL,
  structure_id BIGINT NOT NULL,
  system_id BIGINT NOT NULL,
  type_id BIGINT NOT NULL,
  profile_id BIGINT NOT NULL,
  curr_vuln TEXT NOT NULL,
  next_vuln TEXT NOT NULL,
  fetched_at TIMESTAMP NOT NULL DEFAULT NOW()
);
