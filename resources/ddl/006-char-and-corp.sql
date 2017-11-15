DROP TABLE IF EXISTS app.characters;
CREATE TABLE app.characters
(
  character_id BIGINT PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL,
  race_id INT NOT NULL,
  bloodline_id INT NOT NULL,
  ancestry_id INT NOT NULL,
  corporation_id BIGINT NOT NULL,
  alliance_id BIGINT NOT NULL,
  birth_date TIMESTAMP NOT NULL,
  description VARCHAR(2000) NOT NULL,
  fetched_at TIMESTAMP DEFAULT NOW()
);

DROP TABLE IF EXISTS app.corporations;
CREATE TABLE app.corporations
(
  corporation_id BIGINT PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL,
  alliance_id BIGINT NOT NULL,
  creation_date TIMESTAMP NOT NULL,
  ticker VARCHAR(10) NOT NULL,
  description VARCHAR(2000) NOT NULL,
  fetched_at TIMESTAMP DEFAULT NOW()
);

DROP TABLE IF EXISTS app.corporation_details;
CREATE TABLE app.corporation_details
(
  corporation_id BIGINT PRIMARY KEY NOT NULL,
  ceo_id BIGINT NOT NULL,
  ceo_name VARCHAR(255) NOT NULL,
  hq_station_id BIGINT NOT NULL,
  hq_station_name VARCHAR(255) NOT NULL,
  faction_id INT NOT NULL,
  member_count INT NOT NULL,
  shares INT NOT NULL,
  hangars BYTEA NOT NULL,
  divisions BYTEA NOT NULL,
  fetched_at TIMESTAMP DEFAULT NOW()
);

DROP TABLE IF EXISTS app.corporation_settings;
CREATE TABLE app.corporation_settings
(
  corporation_id BIGINT PRIMARY KEY NOT NULL,
  opted_in BOOLEAN NOT NULL DEFAULT TRUE,
  opted_in_by INT NOT NULL DEFAULT 0,
  opted_in_at TIMESTAMP DEFAULT NULL,
  created_by INT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

DROP TABLE IF EXISTS app.alliances;
CREATE TABLE app.alliances
(
  alliance_id BIGINT PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL,
  ticker VARCHAR(10) NOT NULL,
  founded_date TIMESTAMP NOT NULL
);