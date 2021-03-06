CREATE TABLE IF NOT EXISTS users
(
    id            INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    login         TEXT  NOT NULL,
    first_name    TEXT  NOT NULL,
    last_name     TEXT  NOT NULL,
    email         TEXT  NOT NULL,
    password_hash BYTEA NOT NULL,
    salt          BYTEA NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS products
(
    id          INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    owner_id     INTEGER REFERENCES users (id) NOT NULL,
    name        TEXT                          NOT NULL,
    per_hour    NUMERIC                       NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product_photos
(
    id         INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    product_id INTEGER REFERENCES products (id) NOT NULL,
    photo      TEXT                             NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
