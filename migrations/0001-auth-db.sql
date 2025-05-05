CREATE DATABASE auth WITH OWNER postgres ENCODING 'UTF8';

\connect auth;

CREATE USER auth WITH ENCRYPTED PASSWORD 'auth_pwd';

CREATE EXTENSION "uuid-ossp";

-- --------------------------------

DROP TABLE IF EXISTS "SignupEmail" CASCADE;
CREATE TABLE "SignupEmail"
(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    passwd_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(127) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at TIMESTAMP NULL
);
--
CREATE INDEX ON "SignupEmail" (email);
--
COMMENT ON TABLE "SignupEmail" is 'Регистрация по емейл';
COMMENT ON COLUMN "SignupEmail".email is 'Емейл пользователя';
COMMENT ON COLUMN "SignupEmail".passwd_hash is 'SHA-256-хеш пароля';
COMMENT ON COLUMN "SignupEmail".salt is 'Соль для хеша';
COMMENT ON COLUMN "SignupEmail".created_at is 'Создание записи по UTC';
COMMENT ON COLUMN "SignupEmail".updated_at is 'Время последнего обновления';

-- --------------------------------

DROP TABLE IF EXISTS "RefreshToken";
CREATE TABLE "RefreshToken"
(
    id              UUID            DEFAULT uuid_generate_v4() PRIMARY KEY,
    account_id      UUID            NOT NULL,
    token           VARCHAR(1020)   NOT NULL UNIQUE,
    expires_at      TIMESTAMP       NOT NULL,
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      TIMESTAMP       NULL
);
--
CREATE INDEX ON "RefreshToken" (account_id);
--
COMMENT ON TABLE "RefreshToken" is 'Таблица для хранения токенов';
COMMENT ON COLUMN "RefreshToken".account_id is 'ID аккаунта, которому принадлежит токен';
COMMENT ON COLUMN "RefreshToken".token is 'Сам токен (Refresh Token)';
COMMENT ON COLUMN "RefreshToken".expires_at is 'Время истечения токена';
COMMENT ON COLUMN "RefreshToken".created_at is 'Время создания записи';
COMMENT ON COLUMN "RefreshToken".updated_at is 'Время последнего обновления';

GRANT ALL PRIVILEGES ON DATABASE auth TO auth;
GRANT ALL ON SCHEMA public TO auth;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO auth;
GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON ALL TABLES IN SCHEMA public TO auth;