CREATE TABLE IF NOT EXISTS orders
(
    id          INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id     INTEGER REFERENCES users (id)    NOT NULL,
    product_id  INTEGER REFERENCES products (id) NOT NULL,
    order_start TIMESTAMPTZ                      NOT NULL,
    order_end   TIMESTAMPTZ                      NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT now()
);