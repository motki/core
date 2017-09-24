DROP TABLE IF EXISTS app.mailing_lists;
CREATE TABLE app.mailing_lists
(
  key VARCHAR(100) NOT NULL,
  name VARCHAR(30) NOT NULL,
  email VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY(key, email)
);