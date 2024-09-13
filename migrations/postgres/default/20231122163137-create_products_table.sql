-- +migrate Up

CREATE TABLE IF NOT EXISTS products
(
    id         bigserial    NOT NULL,
    name       varchar(128) NOT NULL DEFAULT '',
    "desc"     varchar(255) NOT NULL DEFAULT '',
    price      int          NOT NULL DEFAULT 0,
    created_at bigint       NOT NULL DEFAULT 0,
    updated_at bigint       NOT NULL DEFAULT 0,
    deleted_at bigint       NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
);

CREATE INDEX ON products (name);
CREATE INDEX ON products (deleted_at);

COMMENT ON COLUMN products.name IS '名称';
COMMENT ON COLUMN products.desc IS '描述';
COMMENT ON COLUMN products.price IS '价格';

COMMENT ON TABLE products IS '产品表';

-- +migrate Down

DROP TABLE IF EXISTS products;