BEGIN;

CREATE TABLE IF NOT EXISTS events
(
    id          UUID PRIMARY KEY,
    name        TEXT      NOT NULL,
    description TEXT,
    place       TEXT      NOT NULL,
    age_rating  INTEGER   NOT NULL,
    start_at    TIMESTAMP NOT NULL,
    end_at      TIMESTAMP NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS promocodes
(
    id               UUID PRIMARY KEY,
    limit_use        INTEGER   NOT NULL,
    discount_value   INTEGER,
    discount_percent INTEGER,
    created_at       TIMESTAMP NOT NULL,
    updated_at       TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS ticket_categories
(
    id          UUID PRIMARY KEY,
    event_id    UUID REFERENCES events (id) ON UPDATE CASCADE,
    price       INTEGER   NOT NULL,
    name        TEXT      NOT NULL,
    description TEXT,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS promocodes_ticket_categories
(
    id                 UUID PRIMARY KEY,
    ticket_category_id UUID REFERENCES ticket_categories (id) ON UPDATE CASCADE,
    promocode_id       UUID REFERENCES promocodes (id) ON UPDATE CASCADE,
    created_at         TIMESTAMP NOT NULL,
    updated_at         TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS orders
(
    id          UUID PRIMARY KEY,
    name        TEXT      NOT NULL,
    email       TEXT      NOT NULL,
    phone       TEXT      NOT NULL,
    status      TEXT      NOT NULL,
    payment_id  TEXT,
    payment_url TEXT,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS tickets
(
    id                 UUID PRIMARY KEY,
    order_id           UUID REFERENCES orders (id) ON UPDATE CASCADE,
    ticket_category_id UUID REFERENCES ticket_categories (id) ON UPDATE CASCADE,
    promocode_id       UUID REFERENCES promocodes (id) ON UPDATE CASCADE,
    full_price         INTEGER   NOT NULL,
    buy_price          INTEGER   NOT NULL,
    created_at         TIMESTAMP NOT NULL,
    updated_at         TIMESTAMP NOT NULL
);

COMMIT;
