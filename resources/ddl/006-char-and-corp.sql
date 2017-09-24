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

DROP TABLE IF EXISTS app.alliances;
CREATE TABLE app.alliances
(
  alliance_id BIGINT PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL,
  ticker VARCHAR(10) NOT NULL,
  founded_date TIMESTAMP NOT NULL
);