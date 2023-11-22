-- +migrate Up

CREATE TABLE IF NOT EXISTS `products`
(
    `id`         integer PRIMARY KEY AUTOINCREMENT,
    `name`       varchar(128)    NOT NULL DEFAULT '', -- 名称
    `desc`       varchar(255)    NOT NULL DEFAULT '', -- 描述
    `price`      int             NOT NULL DEFAULT 0,  -- 价格
    `created_at` bigint          NOT NULL DEFAULT 0,
    `updated_at` bigint          NOT NULL DEFAULT 0,
    `deleted_at` bigint unsigned NOT NULL DEFAULT 0
);

CREATE INDEX products_name ON products (name);
CREATE INDEX products_deleted_at ON products (deleted_at);

-- +migrate Down

DROP TABLE IF EXISTS `products`;