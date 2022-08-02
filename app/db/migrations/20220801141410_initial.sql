-- +goose Up
CREATE TABLE countries 
(
  id              BIGSERIAL     PRIMARY KEY,
  code            VARCHAR(5)    NOT NULL,
  name            VARCHAR(50)   NOT NULL,
  currency        VARCHAR(5)    NOT NULL,
  dialling_code   VARCHAR(10)   NOT NULL
);

CREATE UNIQUE INDEX countries_code_uniq_idx ON countries(code);

CREATE TABLE companies 
(
  id                BIGSERIAL         PRIMARY KEY,
  name              VARCHAR(50)       NOT NULL,
  code              VARCHAR(10)       NOT NULL,
  website           VARCHAR(100)      NOT NULL,
  country_code      VARCHAR(5)        NOT NULL,
  number            VARCHAR(15)       NOT NULL,
  updated_at        TIMESTAMP         NOT NULL DEFAULT clock_timestamp(),
  created_at        TIMESTAMP         NOT NULL DEFAULT clock_timestamp()
);

CREATE UNIQUE INDEX companies_name_uniq_idx ON companies(name);
CREATE UNIQUE INDEX companies_code_uniq_idx ON companies(code);
CREATE UNIQUE INDEX companies_website_uniq_idx ON companies(website);
CREATE UNIQUE INDEX companies_number_uniq_idx ON companies(country_code, number);

CREATE TYPE OPERATION_STATUS AS ENUM('active', 'closed', 'pending');

CREATE TABLE company_countries 
(
  id                BIGSERIAL         PRIMARY KEY,
  company_id        BIGINT            NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
  country_id        BIGINT            NOT NULL REFERENCES countries(id),
  operation_status  OPERATION_STATUS  NOT NULL DEFAULT 'pending',
  updated_at        TIMESTAMP         NOT NULL DEFAULT clock_timestamp(),
  created_at        TIMESTAMP         NOT NULL DEFAULT clock_timestamp()
);

CREATE UNIQUE INDEX company_countries_uniq_idx ON company_countries(company_id, country_id);
 
CREATE TYPE USER_STATUS AS ENUM ('active', 'deactivated', 'unverified');

CREATE TABLE users
(
  id                BIGSERIAL       PRIMARY KEY,
  first_name        VARCHAR(20)     NOT NULL,
  last_name         VARCHAR(20)     NOT NULL,
  username          VARCHAR(20)     NOT NULL,
  password_hash     VARCHAR(255)    NOT NULL,
  status            USER_STATUS     NOT NULL,
  created_at        TIMESTAMPTZ     NOT NULL DEFAULT clock_timestamp(),
  updated_at        TIMESTAMPTZ     NOT NULL DEFAULT clock_timestamp()
);

CREATE UNIQUE INDEX users_username_uniq_idx ON users(username);

CREATE TABLE sessions
(
  id                  BIGSERIAL        PRIMARY KEY,
  deactivated_at      TIMESTAMPTZ      NULL,
  ip_address          VARCHAR(255)     NOT NULL,
  last_refreshed_at   TIMESTAMPTZ      NOT NULL,
  user_agent          VARCHAR(255)     NOT NULL,
  user_id             INTEGER          REFERENCES users(id),
  status              USER_STATUS      NOT NULL DEFAULT 'unverified',
  created_at          TIMESTAMPTZ      NOT NULL DEFAULT clock_timestamp(),
  updated_at          TIMESTAMPTZ      NOT NULL DEFAULT clock_timestamp()
);

CREATE INDEX sessions_idx ON sessions(user_id) WHERE user_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS sessions_idx;
DROP TABLE sessions;

DROP INDEX IF EXISTS users_username_uniq_idx;
DROP TABLE users;
DROP TYPE IF EXISTS USER_STATUS;

DROP INDEX IF EXISTS company_countries_uniq_idx;
DROP TABLE IF EXISTS company_countries;

DROP TYPE IF EXISTS OPERATION_STATUS;

DROP INDEX IF EXISTS companies_number_uniq_idx;
DROP INDEX IF EXISTS companies_website_uniq_idx;
DROP INDEX IF EXISTS companies_name_uniq_idx;

DROP TABLE IF EXISTS companies;

DROP INDEX IF EXISTS countries_code_uniq_idx;
DROP TABLE IF EXISTS countries;
