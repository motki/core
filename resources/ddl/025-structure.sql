DROP TABLE IF EXISTS app.structures;
CREATE TABLE app.structures
(
  structure_id BIGINT NOT NULL,
  corporation_id BIGINT NOT NULL,
  system_id BIGINT NOT NULL,
  type_id BIGINT NOT NULL,
  name VARCHAR(255) NOT NULL,
  profile_id BIGINT NOT NULL,
  fuel_expires TIMESTAMP NOT NULL,
  services TEXT NOT NULL,
  state_timer_start TIMESTAMP NOT NULL,
  state_timer_end TIMESTAMP NOT NULL,
  curr_state VARCHAR(255) NOT NULL,
  unanchors_at TIMESTAMP NOT NULL,
  reinforce_weekday INT NOT NULL,
  reinforce_hour INT NOT NULL,
  fetched_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (structure_id)
);
